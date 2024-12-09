package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func main() {
	input := readInput(os.Stdin)
	fmt.Println("Part 1:", part1(input))
}

func readInput(r io.Reader) []int {
	var buf bytes.Buffer
	buf.ReadFrom(r)

	nums := make([]int, 0, len(buf.Bytes()))
	for _, b := range buf.Bytes() {
		if b >= '0' && b <= '9' {
			nums = append(nums, int(b-'0'))
		}
	}

	return nums
}

func part1(disk []int) (sum int) {
	// If the encoding includes an even number of entries, it means it ends on a
	// free slot, which we can ignore for the purposes of calculating the
	// checksum
	if len(disk)%2 == 0 {
		disk = disk[:len(disk)-1]
	}

	for written, lo, hi := 0, 0, len(disk); lo < hi; {
		if lo%2 == 0 {
			// If the lowerbound is even, then it is sitting over a file to add to
			// the checksum
			file := lo / 2
			sum += checksum(written, disk[lo]) * file
			written += disk[lo]
			lo++
		} else if disk[lo] <= disk[hi-1] {
			// The lowerbound is odd, meaning it is sitting over a free slot, and the
			// file that would go in that slot is too big, so the free slot will be
			// entirely filled up, and the size deducted from the file.
			file := (hi - 1) / 2
			sum += checksum(written, disk[lo]) * file
			written += disk[lo]
			disk[hi-1] -= disk[lo]
			lo++
		} else {
			// The lowerbound is over a free slot that's bigger than the next file to
			// go in it (from the end), write that file, consuming it, and update the
			// free space remaining. When writing and consuming a file at the end, we
			// also need to consume the free space immediately before it.
			file := (hi - 1) / 2
			sum += checksum(written, disk[hi-1]) * file
			written += disk[hi-1]
			disk[lo] -= disk[hi-1]
			hi -= 2
		}
	}

	return
}

func checksum(off, size int) int {
	return size * (2*off + size - 1) / 2
}
