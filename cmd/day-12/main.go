package main

import (
	"fmt"
	"internal/grid"
	"internal/point"
	"os"
)

type field struct {
	area, perimeter int
}

func main() {
	grid := grid.ReadBytes(os.Stdin)
	renumbered := renumber(grid)
	fmt.Println("Part 1:", part1(renumbered))
}

func part1(g *grid.Grid[int]) (cost int) {
	fields := surveyFields(g)
	for _, f := range fields {
		cost += f.area * f.perimeter
	}
	return
}

func surveyFields(g *grid.Grid[int]) map[int]field {
	fields := make(map[int]field)
	for x, y := range g.Coords() {
		color := *g.Get(x, y)

		curr, ok := fields[color]
		if !ok {
			curr = field{}
		}

		curr.area++
		for _, dir := range []grid.Dir{grid.DIR_U, grid.DIR_R, grid.DIR_D, grid.DIR_L} {
			nextX, nextY := dir.Move(x, y, 1)
			if next := g.Get(nextX, nextY); next == nil || *next != color {
				curr.perimeter++
			}
		}

		fields[color] = curr
	}

	return fields
}

func renumber(input *grid.Grid[byte]) *grid.Grid[int] {
	output := grid.New[int](input.Width, input.Height)

	recolor := 0
	for x, y := range input.Coords() {
		// If the current cell has already been filled, don't start filling from here.
		if cell := output.Get(x, y); *cell == 0 {
			recolor++
			*cell = recolor
		} else {
			continue
		}

		color := *input.Get(x, y)
		frontier := []point.Point{{X: x, Y: y}}
		for len(frontier) > 0 {
			var curr point.Point
			curr, frontier = frontier[0], frontier[1:]

			for _, dir := range []grid.Dir{grid.DIR_U, grid.DIR_R, grid.DIR_D, grid.DIR_L} {
				nextX, nextY := dir.Move(curr.X, curr.Y, 1)
				nextI, nextR := input.Get(nextX, nextY), output.Get(nextX, nextY)
				if nextI != nil && *nextI == color && *nextR == 0 {
					*nextR = recolor
					frontier = append(frontier, point.New(nextX, nextY))
				}
			}
		}
	}

	return output
}