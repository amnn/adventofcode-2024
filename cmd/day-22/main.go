package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	hashes := readInput(os.Stdin)
	fmt.Println("Part 1:", part1(hashes))
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

func part1(hs []uint32) (sum int) {
	for _, h := range hs {
		sum += int(hash(h, 2000))
	}
	return
}

func hash(s uint32, times int) uint32 {
	for i := 0; i < times; i++ {
		s = next(s)
	}
	return s
}

func next(i uint32) uint32 {
	i = (i << 6) ^ i
	i %= (1 << 24)
	i = (i >> 5) ^ i
	i = (i << 11) ^ i
	i %= (1 << 24)
	return i
}
