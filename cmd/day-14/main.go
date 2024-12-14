package main

import (
	"fmt"
	"internal/grid"
	"internal/point"
	"io"
	"os"
	"time"
)

const (
	WIDTH    = 101
	HEIGHT   = 103
	DURATION = 100
)

type cell byte

type robot struct {
	pos point.Point
	vel point.Vec
}

func main() {
	input := readInput(os.Stdin)
	fmt.Println("Part 1:", part1(input))
	part2(input)
}

func part1(robots []robot) (safety int) {
	var tl, tr, bl, br int

	for _, r := range robots {
		end := r.pos.Move(r.vel.Scale(DURATION))

		end.X %= WIDTH
		end.X += WIDTH
		end.X %= WIDTH

		end.Y %= HEIGHT
		end.Y += HEIGHT
		end.Y %= HEIGHT

		switch {
		case end.X < WIDTH/2 && end.Y < HEIGHT/2:
			tl++
		case end.X > WIDTH/2 && end.Y < HEIGHT/2:
			tr++
		case end.X < WIDTH/2 && end.Y > HEIGHT/2:
			bl++
		case end.X > WIDTH/2 && end.Y > HEIGHT/2:
			br++
		}
	}

	return tl * tr * bl * br
}

func part2(robots []robot) {
	for i := 0; ; i++ {
		g := grid.New[cell](WIDTH, HEIGHT)
		for i, r := range robots {
			(*g.Get(r.pos.X, r.pos.Y))++
			next := r.pos.Move(r.vel)
			next.X %= WIDTH
			next.X += WIDTH
			next.X %= WIDTH
			next.Y %= HEIGHT
			next.Y += HEIGHT
			next.Y %= HEIGHT
			robots[i].pos = next
		}

		// Assume that if we have fewer than this many contiguous regions, it might
		// be an interesting output.
		rs := regions(g)
		if rs < 200 {
			fmt.Printf("\033[2JGeneration %d, Regions %d\n%v\n", i, rs, g)
			time.Sleep(5 * time.Second)
		}
	}
}

func readInput(r io.Reader) (robots []robot) {
	var input robot

	for {
		if _, err := fmt.Fscanf(
			r, "p=%d,%d v=%d,%d\n",
			&input.pos.X, &input.pos.Y,
			&input.vel.Dx, &input.vel.Dy,
		); err != nil {
			break
		}

		robots = append(robots, input)
	}

	return
}

func (c cell) Format(f fmt.State, _ rune) {
	if c == 0 {
		fmt.Fprint(f, " ")
	} else {
		fmt.Fprint(f, "#")
	}
}

func regions(input *grid.Grid[cell]) int {
	output := grid.New[bool](input.Width, input.Height)

	regions := 0
	for x, y := range input.Coords() {
		// If the current cell has already been filled, don't start filling from here.
		if input, output := input.Get(x, y), output.Get(x, y); *input > 0 && !*output {
			regions++
			*output = true
		} else {
			continue
		}

		frontier := []point.Point{{X: x, Y: y}}
		for len(frontier) > 0 {
			var curr point.Point
			curr, frontier = frontier[0], frontier[1:]

			for _, dir := range []grid.Dir{grid.DIR_U, grid.DIR_R, grid.DIR_D, grid.DIR_L} {
				nextX, nextY := dir.Move(curr.X, curr.Y, 1)
				nextI, nextO := input.Get(nextX, nextY), output.Get(nextX, nextY)
				if nextI != nil && *nextI > 0 && !*nextO {
					*nextO = true
					frontier = append(frontier, point.New(nextX, nextY))
				}
			}
		}
	}

	return regions
}
