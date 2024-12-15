package main

import (
	"bufio"
	"fmt"
	"internal/grid"
	"internal/point"
	"io"
	"os"
	"slices"
)

type cell byte

const (
	EMPTY cell = iota
	ROBOT
	BOX
	BOX_L
	BOX_R
	WALL
)

func main() {
	input1, moves := readInput(os.Stdin)
	input2 := expand(input1)

	fmt.Println("Part 1:", part1(input1, moves))
	fmt.Println("Part 2:", part2(input2, moves))
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

func part1(g *grid.Grid[cell], moves []grid.Dir) (coords int) {
	robotX, robotY, found := g.Find(ROBOT)
	if !found {
		panic("robot not found")
	}

	for _, m := range moves {
		robotX, robotY = moveLinear(g, m, robotX, robotY)
	}

	for x, y := range g.FindAll(BOX) {
		coords += 100*y + x
	}

	return
}

func part2(g *grid.Grid[cell], moves []grid.Dir) (coords int) {
	robotX, robotY, found := g.Find(ROBOT)
	if !found {
		panic("robot not found")
	}

	for _, m := range moves {
		switch m {
		case grid.DIR_L, grid.DIR_R:
			robotX, robotY = moveLinear(g, m, robotX, robotY)
		case grid.DIR_U, grid.DIR_D:
			robotX, robotY = moveCascade(g, m, robotX, robotY)
		}
	}

	for x, y := range g.FindAll(BOX_L) {
		coords += 100*y + x
	}

	return
}

func moveLinear(g *grid.Grid[cell], d grid.Dir, x, y int) (int, int) {
	step := 1

steps:
	for ; ; step++ {
		switch *g.Get(d.Move(x, y, step)) {
		case WALL:
			return x, y
		case EMPTY:
			break steps
		case BOX, BOX_L, BOX_R:
			continue
		}
	}

	write := EMPTY
	for i := 0; i <= step; i++ {
		cell := g.Get(d.Move(x, y, i))
		write, *cell = *cell, write
	}

	return d.Move(x, y, 1)
}

func moveCascade(g *grid.Grid[cell], d grid.Dir, x, y int) (int, int) {
	next := 0
	var changes []point.Point
	visited := make(map[point.Point]struct{})

	var push = func(x, y int) {
		p := point.New(x, y)
		if _, ok := visited[p]; !ok {
			visited[p] = struct{}{}
			changes = append(changes, p)
		}
	}

	var pop = func() (int, int) {
		p := changes[next]
		next++
		return p.X, p.Y
	}

	// Gather all the cells that need to be pushed, cascading from the robot's
	// current position, in direction `d`. The cascade stops if it encounters a
	// wall, or if all the next positions are empty, meaning the move can happen.
	push(x, y)
	for next < len(changes) {
		currX, currY := pop()
		pushX, pushY := d.Move(currX, currY, 1)

		switch *g.Get(pushX, pushY) {
		case WALL:
			return x, y
		case EMPTY:
			break
		case BOX_L:
			push(pushX, pushY)
			push(pushX+1, pushY)
		case BOX_R:
			push(pushX, pushY)
			push(pushX-1, pushY)
		}
	}

	// Apply the changes backwards to avoid overwriting cells that need to be
	// referenced later. Leave an empty slot in place of the point.
	for _, p := range slices.Backward(changes) {
		cell := g.Get(p.X, p.Y)
		next := g.Get(d.Move(p.X, p.Y, 1))
		*next, *cell = *cell, EMPTY
	}

	// Move the robot to its next position
	return d.Move(x, y, 1)
}

func expand(input *grid.Grid[cell]) *grid.Grid[cell] {
	output := grid.New[cell](input.Width*2, input.Height)

	for x, y := range input.Coords() {
		switch *input.Get(x, y) {
		case ROBOT:
			*output.Get(x*2, y) = ROBOT
		case WALL:
			*output.Get(x*2, y) = WALL
			*output.Get(x*2+1, y) = WALL
		case BOX:
			*output.Get(x*2, y) = BOX_L
			*output.Get(x*2+1, y) = BOX_R
		}
	}

	return output
}

func (c cell) Format(f fmt.State, _ rune) {
	switch c {
	case EMPTY:
		fmt.Fprint(f, ".")
	case ROBOT:
		fmt.Fprint(f, "@")
	case BOX:
		fmt.Fprint(f, "O")
	case WALL:
		fmt.Fprint(f, "#")
	case BOX_L:
		fmt.Fprint(f, "[")
	case BOX_R:
		fmt.Fprint(f, "]")
	}
}
