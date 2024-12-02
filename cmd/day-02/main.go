package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

func main() {
	input := readInput()

	fmt.Println("Part 1:", part1(input))
	fmt.Println("Part 2:", part2(input))
}

func readInput() (input [][]int) {
	for s := bufio.NewScanner(os.Stdin); s.Scan(); {
		fields := strings.Fields(s.Text())
		numbers := make([]int, 0, len(fields))

		for _, field := range fields {
			if num, err := strconv.Atoi(field); err == nil {
				numbers = append(numbers, num)
			}
		}

		input = append(input, numbers)
	}

	return
}

func part1(input [][]int) (safe int) {
	for _, row := range input {
		if isSafe(row) {
			safe++
		}
	}

	return
}

func part2(input [][]int) (almostSafe int) {
	for _, row := range input {
		if isAlmostSafe(row) {
			almostSafe++
		}
	}

	return
}

// A report is considered safe if it is strictly monotonic and the absolute gap
// between consecutive elements does not exceed 3
func isSafe(report []int) bool {
	isIncreasing := slices.IsSortedFunc(report, func(a, b int) int { return a - b })
	isDecreasing := slices.IsSortedFunc(report, func(a, b int) int { return b - a })

	var sign int

	// If these values are the same the slice is either not sorted, or it
	// contains a run of the same number. In either case, it does not match the
	// safety criteria.
	if isIncreasing == isDecreasing {
		return false
	} else if isIncreasing {
		sign = 1
	} else {
		sign = -1
	}

	for i := 1; i < len(report); i++ {
		delta := sign * (report[i] - report[i-1])
		if delta < 1 || 3 < delta {
			return false
		}
	}

	return true
}

// A report is considered almost safe if it is safe, or would be considered
// safe after removing one of its elements.
//
// It's possible to do this in linear time, because the effect of removing an
// element is local to its neighbourhood, but we're going to do the naive thing
// and create copies of the report with elements removed.
func isAlmostSafe(report []int) bool {
	if isSafe(report) {
		return true
	}

	for i := 0; i < len(report); i++ {
		if isSafe(remove(report, i)) {
			return true
		}
	}

	return false
}

func remove(xs []int, i int) []int {
	dst := make([]int, len(xs)-1)
	copy(dst, xs[:i])
	copy(dst[i:], xs[i+1:])
	return dst
}
