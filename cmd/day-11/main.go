package main

import (
	"fmt"
	"io"
	"math/big"
	"os"
)

type stone struct {
	engraving  string
	generation int
}

func main() {
	input := readInput(io.Reader(os.Stdin))

	input1 := make([]big.Int, len(input))
	copy(input1, input)

	input2 := make([]big.Int, len(input))
	copy(input2, input)

	fmt.Println("Part 1:", simulate(input1, 25))
	fmt.Println("Part 2:", simulate(input2, 75))
}

func readInput(r io.Reader) []big.Int {
	input := make([]big.Int, 0)

	var buf int64
	for {
		_, err := fmt.Fscanf(r, "%d", &buf)
		if err != nil {
			break
		}

		input = append(input, *big.NewInt(buf))
	}

	return input
}

func simulate(input []big.Int, reps int) (total int) {
	var (
		bi0    = big.NewInt(int64(0))
		bi1    = big.NewInt(int64(1))
		bi10   = big.NewInt(int64(10))
		bi2024 = big.NewInt(int64(2024))
	)

	cache := make(map[stone]int)
	var count func(*big.Int, int) int
	count = func(engraving *big.Int, reps int) int {
		if reps <= 0 {
			return 1
		}

		key := stone{engraving.Text(10), reps}
		if v, ok := cache[key]; ok {
			return v
		}

		if engraving.Cmp(bi0) == 0 {
			// If the stone is engraved with the number 0, it is replaced by a
			// stone engraved with the number `1`.
			cache[key] = count(bi1, reps-1)
		} else if digits := len(key.engraving); digits%2 == 0 {
			// If the stone is engraved with a number that has an even number of
			// digits, it is replaced by two stones. The left half of the digits
			// are engraved on the new left stone, and the right half of the digits
			// are engraved on the new right stone. (The new numbers don't keep
			// extra leading zeroes: 1000 would become stones 10 and 0.)
			var l, r big.Int
			l.Exp(bi10, big.NewInt(int64(digits/2)), nil)
			l.DivMod(engraving, &l, &r)
			cache[key] = count(&l, reps-1) + count(&r, reps-1)
		} else {
			// If none of the other rules apply, the stone is replaced by a new
			// stone; the old stone's number multiplied by 2024 is engraved on the
			// new stone.
			var n big.Int
			n.Mul(engraving, bi2024)
			cache[key] = count(&n, reps-1)
		}

		return cache[key]
	}

	for _, n := range input {
		total += count(&n, reps)
	}

	return
}
