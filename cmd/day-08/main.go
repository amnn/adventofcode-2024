package main

import (
	"fmt"
	"internal/grid"
	"internal/point"
	"os"
)

type antennae map[byte][]point.Point
type exists struct{}

func main() {
	grid := grid.ReadBytes(os.Stdin)
	fmt.Println("Part 1:", part1(grid))
	fmt.Println("Part 2:", part2(grid))
}

func part1(g *grid.Grid[byte]) int {
	a := findAntennae(g)
	antiNodes := findSingleAntiNodes(g, a)
	return len(antiNodes)
}

func part2(g *grid.Grid[byte]) int {
	a := findAntennae(g)
	antiNodes := findMultiAntiNodes(g, a)
	return len(antiNodes)
}

func findAntennae(g *grid.Grid[byte]) antennae {
	antennae := make(antennae)

	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			if chr := *g.Get(x, y); chr != '.' {
				antennae[chr] = append(antennae[chr], point.New(x, y))
			}
		}
	}

	return antennae
}

func findSingleAntiNodes(g *grid.Grid[byte], a antennae) map[point.Point]exists {
	antiNodes := make(map[point.Point]exists)

	for _, points := range a {
		for _, p := range points {
			for _, q := range points {
				if p == q {
					continue
				}

				// The points below can be visualized as follows:
				//
				//     n <- v -- p <- v -- q - -v -> m
				//
				// Starting from p and q, we calculate the vector v, and use that to
				// move in either direction to get n and m, and we add them to the list
				// of anti-nodes as long as they are within the bounds of the grid.

				v := p.Sub(q)
				n, m := p.Move(v), q.Move(v.Neg())

				if g.Get(n.X, n.Y) != nil {
					antiNodes[n] = exists{}
				}

				if g.Get(m.X, m.Y) != nil {
					antiNodes[m] = exists{}
				}
			}
		}
	}

	return antiNodes
}

// Like `findSingleAntiNodes`, but keep extending anti-nodes out in both
// directions until we hit a wall.
func findMultiAntiNodes(g *grid.Grid[byte], a antennae) map[point.Point]exists {
	antiNodes := make(map[point.Point]exists)

	for _, points := range a {
		for i, p := range points {
			for j, q := range points {
				if i == j {
					continue
				}

				v := p.Sub(q)
				w := v.Neg()

				for n := p; g.Get(n.X, n.Y) != nil; n = n.Move(v) {
					antiNodes[n] = exists{}
				}

				for m := q; g.Get(m.X, m.Y) != nil; m = m.Move(w) {
					antiNodes[m] = exists{}
				}
			}
		}
	}

	return antiNodes
}
