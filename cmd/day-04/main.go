package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type grid struct {
	elems  []byte
	width  int
	height int
}

type dir int

const (
	DIR_U = dir(1 << iota)
	DIR_D
	DIR_L
	DIR_R
)

func main() {
	grid := readGrid(os.Stdin)

	fmt.Println("Part 1", part1(&grid))
}

func readGrid(r io.Reader) grid {
	scanner := bufio.NewScanner(r)

	var width, height int
	elems := make([]byte, 0)
	for scanner.Scan() {
		line := scanner.Bytes()
		height += 1
		width = len(line)
		elems = append(elems, line...)
	}

	return grid{
		elems,
		width,
		height,
	}
}

func part1(g *grid) int {
	type candidate struct {
		x, y int
		d    dir
	}

	word := []byte("XMAS")
	candidates := make(map[candidate]bool)

	// Initialise candidates by finding all potential starting points, and
	// checking if the second letter exists at any of the 8 directions eminating
	// from the starting point.
	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			if *g.get(x, y) != word[0] {
				continue
			}

			for _, d := range []dir{
				DIR_U,
				DIR_U | DIR_R,
				DIR_R,
				DIR_D | DIR_R,
				DIR_D,
				DIR_D | DIR_L,
				DIR_L,
				DIR_U | DIR_L,
			} {
				if chr := g.get(g.move(x, y, 1, d)); chr != nil && *chr == word[1] {
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
			chr := g.get(g.move(c.x, c.y, i, c.d))
			if chr == nil || *chr != word[i] {
				delete(candidates, c)
			}
		}
	}

	return len(candidates)
}

func (g *grid) get(x, y int) *byte {
	if y < 0 || g.height <= y {
		return nil
	}

	if x < 0 || g.width <= x {
		return nil
	}

	return &g.elems[y*g.width+x]
}

func (g *grid) move(x, y, step int, d dir) (dx, dy int) {
	dx, dy = x, y

	if DIR_U&d != 0 {
		dy -= step
	}

	if DIR_D&d != 0 {
		dy += step
	}

	if DIR_L&d != 0 {
		dx -= step
	}

	if DIR_R&d != 0 {
		dx += step
	}

	return
}
