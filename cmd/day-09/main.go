package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func main() {
	input := readInput(os.Stdin)

	input1 := make([]int, len(input))
	copy(input1, input)

	input2 := make([]int, len(input))
	copy(input2, input)

	fmt.Println("Part 1:", part1(input1))
	fmt.Println("Part 2:", part2(input2))
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

// `disk` represents the layout of a disk with alternating file and free sizes.
//
// The part1 function calculates the checksum of a compacted form of the disk,
// where sections of files are moved into free slots. Files can be cut up to
// fit into free slots. Free slots are filled from low to high addresses, but
// by files in reverse order.
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

// Like `part1`, but now we can't cut up files: A file must be smaller than the
// free slot to fill it.
//
// As before, free slots are still filled from low to high addresses, and the
// file to fill it is picked from high to low addresses, but if the file does
// not fit in the slot completely, it won't be moved.
func part2(disk []int) (sum int) {
	// If the encoding includes an even number of entries, it means it ends on a
	// free slot, which we can ignore for the purposes of calculating the
	// checksum
	if len(disk)%2 == 0 {
		disk = disk[:len(disk)-1]
	}

	for written, lo := 0, 0; lo < len(disk); {
		if lo%2 == 0 {
			// If the lowerbound is even, then it is sitting over a file to add to
			// the checksum. That file may have been moved to some earlier free slot,
			// in which case it will leave a negative value behind, so detect that
			// and skip over it.
			if disk[lo] < 0 {
				written -= disk[lo]
			} else {
				sum += checksum(written, disk[lo]) * (lo / 2)
				written += disk[lo]
			}
			lo++
		} else if hi := fill(disk[lo:]); hi != 0 {
			// If the lowerbound is odd, then it is sitting over a free slot. We need
			// to find some later file that fits in this slot.
			hi += lo
			file := hi / 2
			sum += checksum(written, disk[hi]) * file
			written += disk[hi]

			// Move the found file into the free space we are trying to fill, and
			// update its old slot to include a sentinel value to recognise the move.
			disk[lo] -= disk[hi]
			disk[hi] *= -1

			if disk[lo] == 0 {
				lo++
			}
		} else {
			// The lowerbound is over a free slot but we couldn't find a file to fill
			// it, so we skip over it.
			written += disk[lo]
			lo++
		}
	}

	return
}

func checksum(off, size int) int {
	return size * (2*off + size - 1) / 2
}

// Assumes that `disk` starts with an empty slot, and looks for the latest slot
// containing a file that fits in the first slot.
//
// Returns the index of the file that fits, or 0 if no file fits.
func fill(disk []int) int {
	for i := len(disk) - 1; i >= 0; i -= 2 {
		if 0 < disk[i] && disk[i] <= disk[0] {
			return i
		}
	}

	return 0
}
