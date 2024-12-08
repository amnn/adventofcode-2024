package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type equation struct {
	target int64
	terms  []term
}

type term struct {
	value int64
	pad10 int64
}

func main() {
	equations := readInput(os.Stdin)

	fmt.Println("Part 1:", part1(equations))
	fmt.Println("Part 2:", part2(equations))
}

func readInput(r io.Reader) (equations []equation) {
	s := bufio.NewScanner(r)

	for s.Scan() {
		eq := equation{}
		line := s.Text()
		target, rest, found := strings.Cut(line, ":")
		if !found {
			continue
		}

		eq.target, _ = strconv.ParseInt(target, 10, 64)
		for _, field := range strings.Fields(rest) {
			t, _ := strconv.ParseInt(field, 10, 64)
			pad10 := Pow(10, int64(len(field)))
			eq.terms = append(eq.terms, term{t, pad10})
		}

		equations = append(equations, eq)
	}

	return
}

func part1(equations []equation) int64 {
	total := int64(0)

	for _, e := range equations {
		if isSatisfiable(e.target, e.terms /* withConcat */, false) {
			total += e.target
		}
	}

	return total
}

func part2(equations []equation) (total int64) {
	for _, e := range equations {
		if isSatisfiable(e.target, e.terms /* withConcat */, true) {
			total += e.target
		}
	}

	return
}

// Returns true if ther is some combination of operations that will result in
// `e` satisfying its target.
func isSatisfiable(target int64, terms []term, withConcat bool) bool {
	if len(terms) == 1 {
		return terms[0].value == target
	}

	last := terms[len(terms)-1]
	rest := terms[:len(terms)-1]

	if target%last.value == 0 && isSatisfiable(target/last.value, rest, withConcat) {
		return true
	}

	if target > last.value && isSatisfiable(target-last.value, rest, withConcat) {
		return true
	}

	if withConcat && (target-last.value)%last.pad10 == 0 && isSatisfiable(target/last.pad10, rest, withConcat) {
		return true
	}

	return false
}

// Logarithmic time integer exponentiation
func Pow(x, y int64) int64 {
	result, base, exp := int64(1), x, y
	for exp > 0 {
		if exp&1 == 0 {
			base *= base
			exp >>= 1
		} else {
			result *= base
			exp--
		}
	}

	return result
}
