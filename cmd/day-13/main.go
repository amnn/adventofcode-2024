package main

import (
	"fmt"
	"io"
	"os"
)

type machine struct {
	ax, ay, bx, by, x, y int
}

func main() {
	input := readInput(os.Stdin)
	fmt.Println("Part 1:", part1(input))
}

func part1(input []machine) (total int) {
	for _, m := range input {
		if cost, ok := solve(m); ok {
			total += cost
		}
	}
	return
}

// Brute force solver which tries all possible combinations of button presses
func solve(m machine) (min int, soluble bool) {
	const (
		kACost      = 3
		kBCost      = 1
		kMaxPresses = 100
		kMaxCost    = kACost*kMaxPresses + kBCost*kMaxPresses
	)

	min = kMaxCost
	for as := 0; as <= kMaxPresses; as++ {
		for bs := 0; bs <= kMaxPresses; bs++ {
			if m.ax*as+m.bx*bs == m.x && m.ay*as+m.by*bs == m.y {
				cost := kACost*as + kBCost*bs
				if cost < min {
					min = cost
				}
			}
		}
	}

	if kMaxCost == min {
		return 0, false
	} else {
		return min, true
	}
}

func readInput(r io.Reader) []machine {
	input := make([]machine, 0)

	var buf machine
	for {
		if _, err := fmt.Fscanf(r, "Button A: X+%d, Y+%d\n", &buf.ax, &buf.ay); err != nil {
			break
		}

		if _, err := fmt.Fscanf(r, "Button B: X+%d, Y+%d\n", &buf.bx, &buf.by); err != nil {
			break
		}

		if _, err := fmt.Fscanf(r, "Prize: X=%d, Y=%d\n", &buf.x, &buf.y); err != nil {
			break
		}

		input = append(input, buf)

		if _, err := fmt.Fscanf(r, "\n"); err != nil {
			break
		}
	}

	return input
}
