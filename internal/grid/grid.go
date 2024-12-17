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

// Create a new grid with dimensions `width` and `height`, all filled with zero
// values for the element type.
func New[E comparable](width, height int) *Grid[E] {
	return &Grid[E]{make([]E, width*height), width, height}
}

// Read a grid from a `*bufio.Scanner`. The grid is expected to be a
// rectangular arrangement of bytes, with rows represented by a line, ending in
// a newline character. Lines are assumed to be of the same length as each
// other.
//
// Reading stops on EOF or after encountering an empty line.
func ScanFunc[E comparable](s *bufio.Scanner, f func(byte) E) *Grid[E] {
	var width, height int
	var elems []E = make([]E, 0)
	for s.Scan() {
		line := s.Bytes()
		if len(line) == 0 {
			break
		}

		height += 1
		width = len(line)
		for _, b := range line {
			elems = append(elems, f(b))
		}
	}

	return &Grid[E]{elems, width, height}
}

// Read a grid from an `io.Reader`. The grid is expected to be a rectangular
// arrangement of bytes, with rows represented by a line, ending in a newline
// character. Lines are assumed to be of the same length as each other.
//
// Reading stops on EOF or after encountering an empty line.
func ReadFunc[E comparable](r io.Reader, f func(byte) E) *Grid[E] {
	return ScanFunc(bufio.NewScanner(r), f)
}

func ReadBytes(r io.Reader) *Grid[byte] {
	return ReadFunc(r, func(b byte) byte { return b })
}

// Make a shallow copy of the grid `g`.
func (g *Grid[E]) Copy() *Grid[E] {
	elems := make([]E, len(g.elems))
	copy(elems, g.elems)
	return &Grid[E]{elems, g.Width, g.Height}
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

func (g *Grid[E]) Coords() func(yield func(int, int) bool) {
	return func(yield func(int, int) bool) {
		for i := range g.elems {
			if !yield(i%g.Width, i/g.Width) {
				return
			}
		}
	}
}

// Returns the coordinates of all points in grid `g` that match `e`.
func (g *Grid[E]) FindAll(e E) func(yield func(int, int) bool) {
	return func(yield func(int, int) bool) {
		for i, elem := range g.elems {
			if elem == e {
				if !yield(i%g.Width, i/g.Width) {
					return
				}
			}
		}
	}
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
	return (0b0111&d)<<1 | (d >> 3)
}

func (d Dir) RotateCounterClockwise() Dir {
	return (0b0001&d)<<3 | (d >> 1)
}

func (g *Grid[E]) Format(f fmt.State, _ rune) {
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			fmt.Fprintf(f, "%v ", *g.Get(x, y))
		}
		fmt.Fprintln(f)
	}
}
