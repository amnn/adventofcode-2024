package main

import (
	"fmt"
	"internal/grid"
	"internal/point"
	"io"
	"math"
	"os"
)

type progress struct {
	p point.Point
	d int
}

const (
	DIM  = 71
	DROP = 1024
)

func main() {
	points := readInput(os.Stdin)
	fmt.Println("Part 1:", part1(points))
	fmt.Println("Part 2:", part2(points))
}

func readInput(r io.Reader) []point.Point {
	var points []point.Point

	for {
		var p point.Point
		_, err := fmt.Fscanf(r, "%d,%d\n", &p.X, &p.Y)
		if err != nil {
			break
		}

		points = append(points, p)
	}

	return points
}

func part1(points []point.Point) int {
	return floodFillAfter(points[:DROP])
}

func part2(points []point.Point) point.Point {
	for i, p := range points {
		if floodFillAfter(points[:i+1]) == 0 {
			return p
		}
	}

	panic("no solution")
}

func floodFillAfter(points []point.Point) int {
	g := grid.New[int](DIM, DIM)

	// Simulate falling bytes
	for _, p := range points {
		*g.Get(p.X, p.Y) = math.MaxInt
	}

	// Flood fill
	frontier := []progress{{point.Point{X: 0, Y: 0}, 0}}

	var curr progress
	for len(frontier) > 0 {
		curr, frontier = frontier[0], frontier[1:]
		if c := g.Get(curr.p.X, curr.p.Y); c != nil && *c == 0 {
			*c = curr.d
		} else {
			continue
		}

		for _, d := range []grid.Dir{grid.DIR_U, grid.DIR_R, grid.DIR_D, grid.DIR_L} {
			nextX, nextY := d.Move(curr.p.X, curr.p.Y, 1)
			next := progress{point.Point{X: nextX, Y: nextY}, curr.d + 1}
			frontier = append(frontier, next)
		}
	}

	return *g.Get(DIM-1, DIM-1)
}
