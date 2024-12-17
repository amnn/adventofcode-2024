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
	SEAT
)

const (
	STEP_COST = 1
	TURN_COST = 1000
)

func main() {
	maze := readInput(os.Stdin)
	fmt.Println("Part 1:", part1(maze))
	fmt.Println("Part 2:", part2(maze))
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

func part1(g *grid.Grid[cell]) (cost int) {
	dists := dijkstra(g)
	endX, endY, _ := g.Find(END)
	return minCost(dists, endX, endY)
}

func part2(g *grid.Grid[cell]) (seats int) {
	dists := dijkstra(g)
	endX, endY, _ := g.Find(END)
	cost := minCost(dists, endX, endY)

	// Fill the frontier with all the ending configurations that have minimal cost
	frontier := make([]config, 0)
	visited := make(map[config]struct{})
	for _, d := range []grid.Dir{grid.DIR_U, grid.DIR_R, grid.DIR_D, grid.DIR_L} {
		if s, ok := dists[config{endX, endY, d}]; ok && s.dist == cost {
			frontier = append(frontier, config{endX, endY, d})
		}
	}

	// Flood fill the grid from the ending configuration to the starting
	// configuration
	var curr config
	for len(frontier) > 0 {
		last := len(frontier) - 1
		frontier, curr = frontier[:last], frontier[last]
		dist := dists[curr].dist

		if _, ok := visited[curr]; ok {
			continue
		} else {
			visited[curr] = struct{}{}
		}

		*g.Get(curr.x, curr.y) = SEAT

		// Check for optimal paths ending in the current configuration that were
		// preceded by a turn
		prevTurnClockwise := config{curr.x, curr.y, curr.dir.RotateClockwise()}
		if s, ok := dists[prevTurnClockwise]; ok && s.dist == dist-TURN_COST {
			frontier = append(frontier, prevTurnClockwise)
		}

		prevTurnCounterClockwise := config{curr.x, curr.y, curr.dir.RotateCounterClockwise()}
		if s, ok := dists[prevTurnCounterClockwise]; ok && s.dist == dist-TURN_COST {
			frontier = append(frontier, prevTurnCounterClockwise)
		}

		// Check for optimal paths ending in the current configuration that were
		// preceded by a move.
		prevX, prevY := curr.dir.Flip().Move(curr.x, curr.y, 1)
		if c := g.Get(prevX, prevY); c == nil || *c == WALL {
			continue
		}

		prevStep := config{prevX, prevY, curr.dir}
		if s, ok := dists[prevStep]; ok && s.dist == dist-STEP_COST {
			frontier = append(frontier, prevStep)
		}
	}

	fmt.Println(g)
	return g.Count(SEAT)
}

func minCost(dists map[config]*state, x, y int) (cost int) {
	for _, d := range []grid.Dir{grid.DIR_U, grid.DIR_R, grid.DIR_D, grid.DIR_L} {
		if s, ok := dists[config{x, y, d}]; ok {
			if cost == 0 || s.dist < cost {
				cost = s.dist
			}
		}
	}
	return
}

func dijkstra(g *grid.Grid[cell]) map[config]*state {
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
			break
		}

		stepX, stepY := s.dir.Move(s.x, s.y, 1)
		if c := g.Get(stepX, stepY); c != nil && *c != WALL {
			relax(s.dist+STEP_COST, stepX, stepY, s.dir)
		}

		relax(s.dist+TURN_COST, s.x, s.y, s.dir.RotateClockwise())
		relax(s.dist+TURN_COST, s.x, s.y, s.dir.RotateCounterClockwise())
	}

	return dists
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
	case SEAT:
		fmt.Fprint(f, "O")
	}
}
