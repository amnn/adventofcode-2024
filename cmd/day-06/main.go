package main

import (
	"fmt"
	"internal/grid"
	"os"
	"sync"
	"sync/atomic"
)

type cell int

// Assign a power of two to each cell state so that we can super-impose the
// states corresponding to visiting the cell while moving in a particular
// direction.
const (
	EMPTY = cell(iota)
	BLOCK = cell(1 << (iota - 1))
	GUARD
	VISIT_U
	VISIT_R
	VISIT_D
	VISIT_L
)

func main() {
	grid := grid.ReadFunc(os.Stdin, func(b byte) cell {
		switch b {
		case '.':
			return EMPTY
		case '#':
			return BLOCK
		case '^':
			return GUARD
		default:
			panic("Invalid character")
		}
	})

	fmt.Println("Part 1", part1(grid.Copy()))
	fmt.Println("Part 2", part2(grid.Copy()))
}

func part1(g *grid.Grid[cell]) int {
	traverse(g)
	total := g.Width * g.Height
	total -= g.Count(EMPTY)
	total -= g.Count(BLOCK)
	return total
}

func part2(g *grid.Grid[cell]) int {
	var found atomic.Uint64

	var wg sync.WaitGroup
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			if *g.Get(x, y) != EMPTY {
				continue
			}

			wg.Add(1)
			clone := g.Copy()
			*clone.Get(x, y) = BLOCK

			go func() {
				defer wg.Done()
				if traverse(clone) {
					found.Add(1)
				}
			}()
		}
	}

	wg.Wait()
	return int(found.Load())
}

// Simulate the guard walking across the grid, filling in cells as they go.
// Returns `true` if the grid causes the guard to enter a cycle, and false
// otherwise.
func traverse(g *grid.Grid[cell]) bool {
	// The guard always starts facing up
	dir := grid.DIR_U
	x, y, found := g.Find(GUARD)
	if !found {
		panic("Guard not found")
	}

	for {
		// If the guard has visited this cell in this direction before, then we
		// have entered a cycle, otherwise, record that configuration and try and
		// make the next move.
		curr := g.Get(x, y)
		visit := cell(dir << 2)
		if *curr&visit != 0 {
			return true
		} else {
			*curr |= visit
		}

		dx, dy := dir.Move(x, y, 1)
		cell := g.Get(dx, dy)
		if cell == nil {
			break
		}

		if *cell == BLOCK {
			dir = dir.RotateClockwise()
		} else {
			x, y = dx, dy
		}
	}

	return false
}

// Print the cell as a 6-bit binary number so we can easily visualize the
// configurations the guard has visited this cell in, as a bitset.
func (c cell) Format(f fmt.State, _ rune) {
	fmt.Fprintf(f, "%06b", int(c))
}
