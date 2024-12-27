package main

import (
	"bufio"
	"fmt"
	"internal/grid"
	"io"
	"os"
)

type pins []byte

func main() {
	locks, keys := readInput(os.Stdin)

	fmt.Println("Part 1:", part1(locks, keys))
}

func part1(locks, keys []pins) (matching int) {
	for _, lock := range locks {
	nextKey:
		for _, key := range keys {
			for i := range lock {
				if lock[i]+key[i] > 5 {
					continue nextKey
				}
			}

			matching++
		}
	}

	return
}

func readInput(r io.Reader) (locks []pins, keys []pins) {
	s := bufio.NewScanner(r)

	for {
		g := grid.ScanFunc(s, func(b byte) byte {
			return b
		})

		if g.Height == 0 {
			break
		}

		if *g.Get(0, 0) == '#' {
			locks = append(locks, readPins(g))
		} else {
			keys = append(keys, readPins(g))
		}
	}

	return
}

func readPins(g *grid.Grid[byte]) pins {
	p := make(pins, g.Width)

	for x := range g.FindAll('#') {
		p[x]++
	}

	for i := range p {
		p[i]--
	}

	return p
}
