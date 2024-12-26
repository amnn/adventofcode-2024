package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type edge struct {
	a, b string
}

type clique []string

type graph struct {
	edges map[edge]struct{}
	nodes map[string]struct{}
}

func main() {
	g := readInput(os.Stdin)
	fmt.Println("Part 1:", part1(g))
}

func part1(g graph) (total int) {
	for c := range g.nodes {
		for e := range g.edges {
			if _, ok := g.edges[connect(e.a, c)]; !ok {
				continue
			}

			if _, ok := g.edges[connect(e.b, c)]; !ok {
				continue
			}

			if strings.HasPrefix(e.a, "t") || strings.HasPrefix(e.b, "t") || strings.HasPrefix(c, "t") {
				total++
			}
		}
	}

	return total / 3
}

func readInput(r io.Reader) (g graph) {
	g.edges = make(map[edge]struct{})
	g.nodes = make(map[string]struct{})

	s := bufio.NewScanner(r)
	for s.Scan() {
		tokens := strings.SplitN(s.Text(), "-", 2)
		g.nodes[tokens[0]] = struct{}{}
		g.nodes[tokens[1]] = struct{}{}
		g.edges[connect(tokens[0], tokens[1])] = struct{}{}
	}

	return
}

func connect(a, b string) edge {
	if a < b {
		return edge{a, b}
	} else {
		return edge{b, a}
	}
}
