package main

import (
	"bufio"
	"fmt"
	"internal/grid"
	"io"
	"os"
)

type cell byte

const (
	EMPTY cell = iota
	ROBOT
	BOX
	WALL
)

func main() {
	grid, moves := readInput(os.Stdin)
	fmt.Println("Part 1:", part1(grid, moves))
}

func part1(g *grid.Grid[cell], moves []grid.Dir) (coords int) {
	simulate(g, moves)
	for x, y := range g.FindAll(BOX) {
		coords += 100*y + x
	}
	return
}

func simulate(g *grid.Grid[cell], moves []grid.Dir) {
	robotX, robotY, found := g.Find(ROBOT)
	if !found {
		panic("robot not found")
	}

moves:
	for _, m := range moves {
		step := 1

	steps:
		for ; ; step++ {
			switch *g.Get(m.Move(robotX, robotY, step)) {
			case WALL:
				continue moves
			case BOX:
				continue steps
			case EMPTY:
				break steps
			}
		}

		write := EMPTY
		for i := 0; i <= step; i++ {
			cell := g.Get(m.Move(robotX, robotY, i))
			write, *cell = *cell, write
		}

		robotX, robotY = m.Move(robotX, robotY, 1)
	}
}

func readInput(r io.Reader) (*grid.Grid[cell], []grid.Dir) {
	s := bufio.NewScanner(r)
	g := grid.ScanFunc(s, func(b byte) cell {
		switch b {
		case '.':
			return EMPTY
		case '@':
			return ROBOT
		case '#':
			return WALL
		case 'O':
			return BOX
		default:
			panic("invalid cell")
		}
	})

	var moves []grid.Dir
	for s.Scan() {
		line := s.Bytes()
		for _, b := range line {
			switch b {
			case '^':
				moves = append(moves, grid.DIR_U)
			case '>':
				moves = append(moves, grid.DIR_R)
			case 'v':
				moves = append(moves, grid.DIR_D)
			case '<':
				moves = append(moves, grid.DIR_L)
			default:
				panic("invalid move")
			}
		}
	}

	return g, moves
}
