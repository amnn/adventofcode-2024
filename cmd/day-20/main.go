package main

import (
	"fmt"
	"internal/grid"
	"io"
	"math"
	"os"
)

type cell int

const (
	EMPTY cell = math.MinInt
	START      = math.MinInt + 1
	END        = math.MinInt + 2
	WALL       = math.MinInt + 3
)

func main() {
	g := readInput(os.Stdin)
	fmt.Println("Part 1:", part1(g))
}

func part1(g *grid.Grid[cell]) int {
	floodFill(g)
	return countShortcuts(g, 100)
}

func readInput(r io.Reader) *grid.Grid[cell] {
	return grid.ReadFunc(r, func(b byte) cell {
		switch b {
		case '.':
			return EMPTY
		case 'S':
			return START
		case 'E':
			return END
		case '#':
			return WALL
		}

		panic("invalid cell")
	})
}

func floodFill(g *grid.Grid[cell]) {
	x, y, ok := g.Find(START)
	if !ok {
		panic("no start position")
	}

outer:
	for curr := 0; ; curr++ {
		*g.Get(x, y) = cell(curr)

		for _, d := range []grid.Dir{grid.DIR_U, grid.DIR_R, grid.DIR_D, grid.DIR_L} {
			nx, ny := d.Move(x, y, 1)
			if c := g.Get(nx, ny); c != nil && (*c == EMPTY || *c == END) {
				x, y = nx, ny
				continue outer
			}
		}

		break
	}
}

func countShortcuts(g *grid.Grid[cell], minSaving int) (shortcuts int) {
	for x, y := range g.Coords() {
		from := g.Get(x, y)
		if *from == WALL {
			continue
		}

		for _, d := range []grid.Dir{grid.DIR_U, grid.DIR_R, grid.DIR_D, grid.DIR_L} {
			if to := g.Get(d.Move(x, y, 2)); to != nil && *to != WALL {
				if int(*to)-int(*from)-2 >= minSaving {
					shortcuts++
				}
			}
		}

		for _, d := range []grid.Dir{
			grid.DIR_U | grid.DIR_R,
			grid.DIR_R | grid.DIR_D,
			grid.DIR_D | grid.DIR_L,
			grid.DIR_L | grid.DIR_U,
		} {
			if to := g.Get(d.Move(x, y, 1)); to != nil && *to != WALL {
				if int(*to)-int(*from)-2 >= minSaving {
					shortcuts++
				}
			}
		}
	}

	return
}

func (c cell) Format(f fmt.State, r rune) {
	switch c {
	case EMPTY:
		fmt.Fprint(f, "....")
	case START:
		fmt.Fprint(f, "SSSS")
	case END:
		fmt.Fprint(f, "EEEE")
	case WALL:
		fmt.Fprint(f, "####")
	default:
		fmt.Fprintf(f, "%04x", int(c))
	}
}
