package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	towels, patterns := readInput(os.Stdin)
	fmt.Println("Part 1:", part1(towels, patterns))
}

func readInput(r io.Reader) (towels []string, patterns []string) {
	s := bufio.NewScanner(r)

	if !s.Scan() {
		panic("no towels")
	}

	towels = strings.Split(s.Text(), ", ")
	s.Scan()

	for s.Scan() {
		patterns = append(patterns, s.Text())
	}

	return
}

func part1(towels []string, patterns []string) (possible int) {
	cache := make(map[string]bool)
	cache[""] = true

	for _, pattern := range patterns {
		if patternPossible(pattern, towels, cache) {
			possible++
		}
	}
	return
}

func patternPossible(pattern string, towels []string, cache map[string]bool) (possible bool) {
	if possible, ok := cache[pattern]; ok {
		return possible
	}

	for _, towel := range towels {
		if !strings.HasPrefix(pattern, towel) {
			continue
		}

		if patternPossible(pattern[len(towel):], towels, cache) {
			possible = true
			break
		}
	}

	cache[pattern] = possible
	return
}
