package main

import (
	"fmt"
	"internal/grid"
	"os"
)

type cell int

const (
	EMPTY = cell(iota)
	BLOCK
	GUARD
	VISIT
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

	fmt.Println("Part 1", part1(grid))
}

func part1(g *grid.Grid[cell]) int {
	// The guard always starts facing up
	dir := grid.DIR_U
	x, y, found := g.Find(GUARD)
	if !found {
		panic("Guard not found")
	}

	// The guard has visited their starting position
	for {
		curr := g.Get(x, y)
		*curr = VISIT

		dx, dy := dir.Move(x, y, 1)
		cell := g.Get(dx, dy)
		if cell == nil {
			break
		}

		switch *cell {
		case BLOCK:
			dir = dir.RotateClockwise()
		case EMPTY, VISIT:
			x, y = dx, dy
		}
	}

	return g.Count(VISIT)
}
