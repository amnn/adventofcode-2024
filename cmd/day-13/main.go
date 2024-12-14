package main

import (
	"fmt"
	"io"
	"os"
)

type machine struct {
	ax, ay, bx, by, x, y int
}

const (
	A_COST = 3
	B_COST = 1
)

func main() {
	input := readInput(os.Stdin)

	input1 := make([]machine, len(input))
	copy(input1, input)

	input2 := make([]machine, len(input))
	copy(input2, input)

	fmt.Println("Part 1:", part1(input1))
	fmt.Println("Part 2:", part2(input2))
}

func part1(input []machine) (total int) {
	for _, m := range input {
		if cost, ok := solve(m); ok {
			total += cost
		}
	}
	return
}

func part2(input []machine) (total int) {
	for _, m := range input {
		m.x += 10000000000000
		m.y += 10000000000000
		if cost, ok := solve(m); ok {
			total += cost
		}
	}
	return
}

func solve(m machine) (cost int, soluble bool) {
	det := m.ax*m.by - m.ay*m.bx
	if det == 0 {
		return
	}

	a := m.by*m.x - m.bx*m.y
	b := m.ax*m.y - m.ay*m.x

	if a%det != 0 || b%det != 0 {
		return
	}

	soluble = true
	cost = A_COST*(a/det) + B_COST*(b/det)
	return
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
