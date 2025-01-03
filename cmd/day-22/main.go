package main

import (
	"fmt"
	"io"
	"os"
)

type round struct {
	signature uint32
	price     int
}

func main() {
	hashes := readInput(os.Stdin)
	fmt.Println("Part 1:", part1(hashes))
	fmt.Println("Part 2:", part2(hashes))
}

func readInput(r io.Reader) (hs []uint32) {
	for {
		var h uint32
		_, err := fmt.Fscanf(r, "%d\n", &h)
		if err != nil {
			break
		}
		hs = append(hs, h)
	}
	return
}

func part1(seeds []uint32) (sum int) {
	for _, seed := range seeds {
		for i, hash := range hashes(seed) {
			if i == 2000 {
				sum += int(hash)
				break
			}
		}
	}
	return
}

func part2(seeds []uint32) (bestPrice int) {
	bananas := make(map[uint32]int)
	for _, seed := range seeds {
		seen := make(map[uint32]struct{})
		for i, round := range rounds(seed) {
			if i >= 2000 {
				break
			}

			if _, ok := seen[round.signature]; ok {
				continue
			}

			seen[round.signature] = struct{}{}
			bananas[round.signature] += round.price
		}
	}

	for _, price := range bananas {
		if price > bestPrice {
			bestPrice = price
		}
	}

	return
}

// Iterate over the prices and price signatures generated by the given seed.
//
// In each round the iterator produces the index of the round we are at and a
// `round` structure which contains this the signature of the current position
// (an encoding of the last 4 price changes), and the current price.
//
// This iterator only produces values once it can generate the signature and
// hash (so it will not produce anything for the first 4 rounds).
func rounds(s uint32) func(yield func(int, round) bool) {
	return func(yield func(int, round) bool) {
		var (
			r round
			h = s
			i = 0
		)

		update := func() {
			p := int(h % 10)
			d := byte(p - r.price)

			r.signature = (r.signature << 8) | 0xff&uint32(d)
			r.price = p

			h = next(h)
			i++
		}

		// Process enough hashes to generate the first signature
		for ; i < 6; update() {
		}

		// Yield all the remaining results. We subtract one from the round index
		// because `i` represents the round that the hash `h` comes from, but the
		// information in round `r` is computed from the previous value of `h`.
		for ; yield(i-1, r); update() {
		}
	}
}

// Iterate over the hashes generated by the given seed.
//
// In each round the iterator produces the index of the round we are at and the
// hash at that round.
func hashes(s uint32) func(yield func(int, uint32) bool) {
	return func(yield func(int, uint32) bool) {
		for i := 0; ; i, s = i+1, next(s) {
			if !yield(i, s) {
				break
			}
		}
	}
}

func next(i uint32) uint32 {
	i = (i << 6) ^ i
	i %= (1 << 24)
	i = (i >> 5) ^ i
	i = (i << 11) ^ i
	i %= (1 << 24)
	return i
}

func (r round) Format(f fmt.State, _ rune) {
	fmt.Fprintf(f, "%08x -> %d", r.signature, r.price)
}
