package main

import (
	"bufio"
	"fmt"
	"internal/point"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

type cacheKey struct {
	move point.Vec
	bias bias
}

type cache map[cacheKey]string

// Whether we need to bias the first direction we move in to avoid the X slot.
type bias byte

type keyPad map[rune]point.Point

const (
	VERT bias = 0b01
	HORZ bias = 0b10
	NONE bias = 0b11
)

var NUMERIC_KEYS = keyPad{
	'7': {X: 0, Y: 0},
	'8': {X: 1, Y: 0},
	'9': {X: 2, Y: 0},
	'4': {X: 0, Y: 1},
	'5': {X: 1, Y: 1},
	'6': {X: 2, Y: 1},
	'1': {X: 0, Y: 2},
	'2': {X: 1, Y: 2},
	'3': {X: 2, Y: 2},
	'X': {X: 0, Y: 3},
	'0': {X: 1, Y: 3},
	'A': {X: 2, Y: 3},
}

var DIRECTION_KEYS = keyPad{
	'X': {X: 0, Y: 0},
	'^': {X: 1, Y: 0},
	'A': {X: 2, Y: 0},
	'<': {X: 0, Y: 1},
	'v': {X: 1, Y: 1},
	'>': {X: 2, Y: 1},
}

func main() {
	codes := readInput(os.Stdin)
	fmt.Println("Part 1:", part1(codes))
	fmt.Println("Part 2:", part2(codes))
}

func part1(codes []string) (total int) {
	var caches []cache
	for i := 0; i < 3; i++ {
		caches = append(caches, make(cache))
	}

	for _, code := range codes {
		path := cascadingOptimalPath(code, NUMERIC_KEYS, caches)
		n, _ := strconv.Atoi(code[:len(code)-1])
		total += len(path) * n
	}

	return
}

func part2(codes []string) (total int) {
	var caches []cache
	for i := 0; i < 25; i++ {
		caches = append(caches, make(cache))
	}

	for _, code := range codes {
		path := cascadingOptimalPath(code, NUMERIC_KEYS, caches)
		n, _ := strconv.Atoi(code[:len(code)-1])
		total += len(path) * n
	}

	return
}

func readInput(r io.Reader) (codes []string) {
	s := bufio.NewScanner(r)

	for s.Scan() {
		codes = append(codes, s.Text())
	}

	return
}

func cascadingOptimalPath(
	code string,
	pad keyPad,
	caches []cache,
) string {
	if len(caches) == 0 {
		return code
	}

	return optimalPath(code, pad, caches[0], func(step string) string {
		return cascadingOptimalPath(step, DIRECTION_KEYS, caches[1:])
	})
}

func optimalPath(
	code string,
	pad keyPad,
	cache cache,
	plan func(string) string,
) string {
	curr := 'A'
	var b strings.Builder
	for _, next := range code {
		b.WriteString(optimalStep(curr, next, pad, cache, plan))
		curr = next
	}

	return b.String()
}

func optimalStep(
	curr, next rune,
	pad keyPad,
	cache cache,
	plan func(code string) string,
) string {
	corner := pad['X']
	from, to := pad[curr], pad[next]
	v := to.Sub(from)

	var bias bias
	switch {
	case from.X == corner.X && to.Y == corner.Y:
		bias = HORZ
	case from.Y == corner.Y && to.X == corner.X:
		bias = VERT
	default:
		bias = NONE
	}

	key := cacheKey{v, bias}
	if path, ok := cache[key]; ok {
		return path
	}

	var x, y string
	if v.Dy < 0 {
		y = strings.Repeat("^", -v.Dy)
	} else {
		y = strings.Repeat("v", v.Dy)
	}

	if v.Dx < 0 {
		x = strings.Repeat("<", -v.Dx)
	} else {
		x = strings.Repeat(">", v.Dx)
	}

	minLen := math.MaxInt
	var b strings.Builder
	var path string

	if bias&HORZ != 0 {
		b.Reset()
		b.WriteString(x)
		b.WriteString(y)
		b.WriteRune('A')
		xy := plan(b.String())
		if len(xy) < minLen {
			minLen = len(xy)
			path = xy
		}
	}

	if bias&VERT != 0 {
		b.Reset()
		b.WriteString(y)
		b.WriteString(x)
		b.WriteRune('A')
		yx := plan(b.String())
		if len(yx) < minLen {
			minLen = len(yx)
			path = yx
		}
	}

	cache[key] = path
	return path
}
