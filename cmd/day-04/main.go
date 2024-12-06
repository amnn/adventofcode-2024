package main

import (
	"fmt"
	"internal/grid"
	"os"
)

func main() {
	grid := grid.ReadBytes(os.Stdin)

	fmt.Println("Part 1", part1(grid))
	fmt.Println("Part 2", part2(grid))
}

func part1(g *grid.Grid[byte]) int {
	type candidate struct {
		x, y int
		d    grid.Dir
	}

	word := []byte("XMAS")
	candidates := make(map[candidate]bool)

	// Initialise candidates by finding all potential starting points, and
	// checking if the second letter exists at any of the 8 directions eminating
	// from the starting point.
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			if *g.Get(x, y) != word[0] {
				continue
			}

			for _, d := range []grid.Dir{
				grid.DIR_U,
				grid.DIR_U | grid.DIR_R,
				grid.DIR_R,
				grid.DIR_D | grid.DIR_R,
				grid.DIR_D,
				grid.DIR_D | grid.DIR_L,
				grid.DIR_L,
				grid.DIR_U | grid.DIR_L,
			} {
				if chr := g.Get(d.Move(x, y, 1)); chr != nil && *chr == word[1] {
					candidates[candidate{x, y, d}] = true
				}
			}
		}
	}

	// Then go through each successive character in the word and check whether
	// we can find that character at the appropriate distance and direction from
	// the candidate starting point.
	for i := 2; i < len(word); i++ {
		for c := range candidates {
			chr := g.Get(c.d.Move(c.x, c.y, i))
			if chr == nil || *chr != word[i] {
				delete(candidates, c)
			}
		}
	}

	return len(candidates)
}

func part2(g *grid.Grid[byte]) (total int) {
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			if *g.Get(x, y) != 'A' {
				continue
			}

			ul := g.Get((grid.DIR_U | grid.DIR_L).Move(x, y, 1))
			ur := g.Get((grid.DIR_U | grid.DIR_R).Move(x, y, 1))
			dl := g.Get((grid.DIR_D | grid.DIR_L).Move(x, y, 1))
			dr := g.Get((grid.DIR_D | grid.DIR_R).Move(x, y, 1))

			if xMas(ul, dr) && xMas(ur, dl) {
				total += 1
			}
		}
	}

	return
}

func xMas(a, b *byte) bool {
	if a == nil || b == nil {
		return false
	} else {
		return *a == 'M' && *b == 'S' || *a == 'S' && *b == 'M'
	}
}
