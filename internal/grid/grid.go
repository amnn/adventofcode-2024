package grid

import (
	"bufio"
	"io"
)

type Grid[E any] struct {
	elems  []E
	Width  int
	Height int
}

type Dir int

const (
	DIR_U = Dir(1 << iota)
	DIR_D
	DIR_L
	DIR_R
)

// Read a grid from an `io.Reader`. The grid is expected to be a rectangular
// arrangement of bytes, with rows represented by a line, ending in a newline
// character. Lines are assumed to be of the same length as each other.
func ReadFunc[E any](r io.Reader, f func(byte) E) *Grid[E] {
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
