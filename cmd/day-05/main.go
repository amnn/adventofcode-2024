package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
)

type order map[edge]bool

type edge struct {
	before, after int
}

func main() {
	rules, updates := readInput(os.Stdin)

	fmt.Println("Part 1:", part1(rules, updates))
	fmt.Println("Part 2:", part2(rules, updates))
}

func readInput(r io.Reader) (rules order, updates [][]int) {
	s := bufio.NewScanner(r)
	rules = make(order)

	// Read ordering rules
	for s.Scan() {
		var before, after int
		if _, err := fmt.Sscanf(s.Text(), "%d|%d", &before, &after); err == nil {
			rules[edge{before, after}] = true
		} else {
			break
		}
	}

	// Read pages
	for s.Scan() {
		tokens := strings.Split(s.Text(), ",")
		pages := make([]int, 0, len(tokens))
		for _, token := range tokens {
			if page, err := strconv.Atoi(token); err == nil {
				pages = append(pages, page)
			} else {
				panic("Failed to read number in updates")
			}
		}
		updates = append(updates, pages)
	}

	return
}

func (o order) before(a, b int) bool {
	_, ok := o[edge{a, b}]
	return ok
}

// Convert the ordering table into a comparator function
func (o order) cmp() func(a, b int) int {
	return func(i, j int) int {
		switch {
		case o.before(i, j):
			return -1
		case o.before(j, i):
			return +1
		default:
			return 0
		}
	}
}

func part1(rules order, updates [][]int) (total int) {
	for _, update := range updates {
		if !slices.IsSortedFunc(update, rules.cmp()) {
			continue
		}

		// If the criteria is met, then add the middle element to the total
		total += update[len(update)/2]
	}

	return
}

func part2(rules order, updates [][]int) (total int) {
	for _, update := range updates {
		if slices.IsSortedFunc(update, rules.cmp()) {
			continue
		}

		// If the criteria is not met, figure out what it would look like if it was
		// met (in a copy), and then get that copy's middle element.
		copied := make([]int, len(update))
		copy(copied, update)

		slices.SortStableFunc(copied, rules.cmp())
		total += copied[len(copied)/2]
	}

	return

}
