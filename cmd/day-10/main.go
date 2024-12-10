package main

import (
	"fmt"
	"internal/grid"
	"internal/point"
	"os"
)

type exists struct{}
type set[T comparable] map[T]exists

type path struct {
	from, to point.Point
}

func main() {
	grid := grid.ReadFunc(os.Stdin, func(b byte) byte {
		if '0' <= b && b <= '9' {
			return b - '0'
		} else {
			return 255
		}
	})

	fmt.Println("Part 1:", part1(grid))
	fmt.Println("Part 2:", part2(grid))
}

func part1(g *grid.Grid[byte]) (score int) {
	visited := make(set[path])

	var frontier []path
	for x, y := range g.FindAll(0) {
		origin := point.New(x, y)
		frontier = append(frontier, path{origin, origin})
	}

	for len(frontier) > 0 {
		curr := frontier[0]
		frontier = frontier[1:]

		if _, ok := visited[curr]; ok {
			continue
		} else {
			visited[curr] = exists{}
		}

		pos := g.Get(curr.to.X, curr.to.Y)
		if *pos == 9 {
			score += 1
			continue
		}

		for _, dir := range []grid.Dir{grid.DIR_U, grid.DIR_R, grid.DIR_D, grid.DIR_L} {
			nextX, nextY := dir.Move(curr.to.X, curr.to.Y, 1)
			if next := g.Get(nextX, nextY); next != nil && *next == *pos+1 {
				frontier = append(frontier, path{curr.from, point.New(nextX, nextY)})
			}
		}
	}

	return
}

func part2(g *grid.Grid[byte]) (rating int) {
	var frontier []point.Point
	for x, y := range g.FindAll(0) {
		frontier = append(frontier, point.New(x, y))
	}

	for len(frontier) > 0 {
		curr := frontier[0]
		frontier = frontier[1:]

		pos := g.Get(curr.X, curr.Y)
		if *pos == 9 {
			rating += 1
			continue
		}

		for _, dir := range []grid.Dir{grid.DIR_U, grid.DIR_R, grid.DIR_D, grid.DIR_L} {
			nextX, nextY := dir.Move(curr.X, curr.Y, 1)
			if next := g.Get(nextX, nextY); next != nil && *next == *pos+1 {
				frontier = append(frontier, point.New(nextX, nextY))
			}
		}
	}

	return
}
