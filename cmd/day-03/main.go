package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
)

var reMul = regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)`)
var reCmd = regexp.MustCompile(`(do)\(\)|(don't)\(\)|(mul)\((\d{1,3}),(\d{1,3})\)`)

func main() {
	inputBytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	input := string(inputBytes)
	fmt.Println("Part 1:", part1(input))
	fmt.Println("Part 2:", part2(input))
}

func part1(input string) (total int) {
	for _, match := range reMul.FindAllStringSubmatch(input, -1) {
		a, b := match[1], match[2]

		ai, err := strconv.Atoi(a)
		if err != nil {
			continue
		}

		bi, err := strconv.Atoi(b)
		if err != nil {
			continue
		}

		total += ai * bi
	}

	return
}

func part2(input string) (total int) {
	enabled := true
	for _, match := range reCmd.FindAllStringSubmatch(input, -1) {
		switch {
		case match[1] == "do":
			enabled = true
		case match[2] == "don't":
			enabled = false
		case match[3] == "mul":
			if !enabled {
				continue
			}

			a, b := match[4], match[5]
			ai, err := strconv.Atoi(a)
			if err != nil {
				continue
			}

			bi, err := strconv.Atoi(b)
			if err != nil {
				continue
			}

			total += ai * bi
		}
	}

	return
}
