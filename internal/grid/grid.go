package grid

import (
	"bufio"
	"fmt"
	"io"
)

type Grid[E comparable] struct {
	elems  []E
	Width  int
	Height int
}

type Dir int

const (
	DIR_U = Dir(1 << iota)
	DIR_R
	DIR_D
	DIR_L
)

// Read a grid from an `io.Reader`. The grid is expected to be a rectangular
// arrangement of bytes, with rows represented by a line, ending in a newline
// character. Lines are assumed to be of the same length as each other.
func ReadFunc[E comparable](r io.Reader, f func(byte) E) *Grid[E] {
	s := bufio.NewScanner(r)

	var width, height int
	var elems []E = make([]E, 0)
	for s.Scan() {
		line := s.Bytes()
		height += 1
		width = len(line)
		for _, b := range line {
			elems = append(elems, f(b))
		}
	}

	return &Grid[E]{elems, width, height}
}

func ReadBytes(r io.Reader) *Grid[byte] {
	return ReadFunc(r, func(b byte) byte { return b })
}

// Access an element from the grid by its position. Returns `nil` if the access
// is out of bounds.
func (g *Grid[E]) Get(x, y int) *E {
	if y < 0 || g.Height <= y {
		return nil
	}

	if x < 0 || g.Width <= x {
		return nil
	}

	return &g.elems[y*g.Width+x]
}

// Returns the first (going from top to bottom, left to right) matching
// position for element `e` in grid `g`, if there is one. The return values are
// the x and y coordinates of the element, and a flag indicating whether it was
// actually found.
func (g *Grid[E]) Find(e E) (x, y int, found bool) {
	for i, elem := range g.elems {
		if elem == e {
			return i % g.Width, i / g.Width, true
		}
	}

	return 0, 0, false
}

// Counts the number of occurrences of `e` in `g`.
func (g *Grid[E]) Count(e E) (count int) {
	for _, elem := range g.elems {
		if elem == e {
			count += 1
		}
	}

	return
}

// Move `step` units in direction `d` from position `(x, y)`.
//
// If `d` is a combination of directions (e.g. `DIR_U|DIR_L`), moves are made
// in each direction. This also means that conflicting directions will cancel
// themselves out.
func (d Dir) Move(x, y int, step int) (dx, dy int) {
	dx, dy = x, y

	if DIR_U&d != 0 {
		dy -= step
	}

	if DIR_D&d != 0 {
		dy += step
	}

	if DIR_L&d != 0 {
		dx -= step
	}

	if DIR_R&d != 0 {
		dx += step
	}

	return
}

func (d Dir) RotateClockwise() Dir {
	return (0b111&d)<<1 | (d >> 3)
}

func (g *Grid[E]) Format(f fmt.State, _ rune) {
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			fmt.Fprintf(f, "%v ", *g.Get(x, y))
		}
		fmt.Fprintln(f)
	}
}
