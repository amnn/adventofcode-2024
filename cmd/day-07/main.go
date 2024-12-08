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
		if isSatisfiable(e) {
			total += e.target
		}
	}

	return total
}

func part2(equations []equation) (total int64) {
	for _, e := range equations {
		if isSatifisableWithConcat(e) {
			total += e.target
		}
	}

	return
}

// Returns true if ther is some combination of operations that will result in
// `e` satisfying its target.
func isSatisfiable(e equation) bool {
	for ops := int64(0); ops < int64(1)<<len(e.terms); ops++ {
		if evaluate(e, ops) {
			return true
		}
	}

	return false
}

// Evaluate an equation with the operations specified in `ops`. Operations can
// be `+` or `*`, evaluated from left to right, and they are represented by the
// bits of an integer. The least significant bit represents the first
// operation, and so on. A bit set to `0` means that the operation is `+`, and
// a bit set to `1` means that the operation is `*`.
//
// Returns `true` if the equation evaluates to the target value, `false`
// otherwise.
func evaluate(e equation, ops int64) bool {
	actual := e.terms[0].value

	for i, term := range e.terms[1:] {
		// If the actual value exceeds the target, we can stop early, because it's
		// never going to get smaller
		if actual > e.target {
			return false
		}

		if ops&(1<<i) == 0 {
			actual += term.value
		} else {
			actual *= term.value
		}
	}

	return actual == e.target
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

// Like `isSatisfiable` but with an additional concatenation operation.
func isSatifisableWithConcat(e equation) bool {
	for ops := int64(0); ops < Pow(int64(3), int64(len(e.terms))); ops++ {
		if evaluateWithConcat(e, ops) {
			return true
		}
	}

	return false
}

// Like `evaluate`, but `ops` is a stream of one of three operations: addition,
// multiplication, and concatenation. It is still represented by a number, but
// it is decoded by converting into base 3.
func evaluateWithConcat(e equation, ops int64) bool {
	actual := e.terms[0].value

	for _, term := range e.terms[1:] {
		// If the actual value exceeds the target, we can stop early, because it's
		// never going to get smaller
		if actual > e.target {
			return false
		}

		op := ops % 3
		ops /= 3

		switch op {
		case 0:
			actual += term.value
		case 1:
			actual *= term.value
		case 2:
			actual *= term.pad10
			actual += term.value
		}
	}

	return actual == e.target
}
