package main

import (
	"container/heap"
	"fmt"
	"internal/grid"
	"io"
	"os"
)

type cell byte

type config struct {
	x, y int
	dir  grid.Dir
}

type state struct {
	x, y  int
	dir   grid.Dir
	dist  int
	index int
}

type priorityQueue []*state

const (
	EMPTY cell = iota
	WALL
	START
	END
)

const (
	STEP_COST = 1
	TURN_COST = 1000
)

func main() {
	maze := readInput(os.Stdin)
	fmt.Println("Part 1:", part1(maze))
}

func readInput(r io.Reader) *grid.Grid[cell] {
	return grid.ReadFunc(r, func(b byte) cell {
		switch b {
		case '.':
			return EMPTY
		case '#':
			return WALL
		case 'S':
			return START
		case 'E':
			return END
		default:
			panic("invalid character")
		}
	})
}

func part1(g *grid.Grid[cell]) int {
	startX, startY, foundStart := g.Find(START)
	if !foundStart {
		panic("no start found")
	}

	endX, endY, foundEnd := g.Find(END)
	if !foundEnd {
		panic("no end found")
	}

	pq := make(priorityQueue, 0)
	dists := make(map[config]*state)

	relax := func(d, x, y int, dir grid.Dir) {
		c := config{x, y, dir}
		if s, ok := dists[c]; !ok {
			s := &state{x, y, dir, d, -1}
			dists[c] = s
			heap.Push(&pq, s)
		} else if d < s.dist {
			s.dist = d
			heap.Fix(&pq, s.index)
		}
	}

	relax(0, startX, startY, grid.DIR_R)

	for pq.Len() > 0 {
		s := heap.Pop(&pq).(*state)
		if s.x == endX && s.y == endY {
			return s.dist
		}

		stepX, stepY := s.dir.Move(s.x, s.y, 1)
		if c := g.Get(stepX, stepY); c != nil && *c != WALL {
			relax(s.dist+STEP_COST, stepX, stepY, s.dir)
		}

		relax(s.dist+TURN_COST, s.x, s.y, s.dir.RotateClockwise())
		relax(s.dist+TURN_COST, s.x, s.y, s.dir.RotateCounterClockwise())
	}

	return 0
}

func (pq priorityQueue) Len() int {
	return len(pq)
}

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].dist < pq[j].dist
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue) Push(x interface{}) {
	n := len(*pq)
	s := x.(*state)

	s.index = n
	*pq = append(*pq, s)
}

func (pq *priorityQueue) Pop() interface{} {
	q := *pq
	n := len(q)
	s := q[n-1]
	s.index = -1
	q[n-1] = nil
	*pq = q[:n-1]
	return s
}

func (c cell) Format(f fmt.State, _ rune) {
	switch c {
	case EMPTY:
		fmt.Fprint(f, ".")
	case WALL:
		fmt.Fprint(f, "#")
	case START:
		fmt.Fprint(f, "S")
	case END:
		fmt.Fprint(f, "E")
	}
}
