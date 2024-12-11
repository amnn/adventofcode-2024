package main

import (
	"fmt"
	"io"
	"math/big"
	"os"
)

func main() {
	input := readInput(io.Reader(os.Stdin))
	fmt.Println("Part 1:", part1(input))
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

func part1(input []big.Int) int {
	curr, next := input, make([]big.Int, 0)

	for i := 0; i < 25; i++ {
		blink(&curr, &next)
		curr, next = next, curr
	}

	return len(curr)
}

func blink(input, output *[]big.Int) {
	*output = (*output)[:0]

	var (
		bi0    = big.NewInt(int64(0))
		bi1    = big.NewInt(int64(1))
		bi10   = big.NewInt(int64(10))
		bi2024 = big.NewInt(int64(2024))
	)

	for _, n := range *input {
		// If the stone is engraved with the number 0, it is replaced by a stone
		// engraved with the number `1`.
		if n.Cmp(bi0) == 0 {
			*output = append(*output, *bi1)
			continue
		}

		// If the stone is engraved with a number that has an even number of
		// digits, it is replaced by two stones. The left half of the digits are
		// engraved on the new left stone, and the right half of the digits are
		// engraved on the new right stone. (The new numbers don't keep extra
		// leading zeroes: 1000 would become stones 10 and 0.)
		digits := len(n.Text(10))
		if digits%2 == 0 {
			*output = append(*output, big.Int{}, big.Int{})
			s0, s1 := &(*output)[len(*output)-2], &(*output)[len(*output)-1]
			s0.Exp(bi10, big.NewInt(int64(digits/2)), nil)
			s0.DivMod(&n, s0, s1)
			continue
		}

		*output = append(*output, big.Int{})
		(*output)[len(*output)-1].Mul(&n, bi2024)
	}
}
