package main

import (
	"bufio"
	"fmt"
	"internal/set"
	"io"
	"os"
	"slices"
	"strings"
)

type edge struct {
	a, b string
}

type graph struct {
	edges map[string]set.Set[string]
	nodes set.Set[string]
}

func main() {
	g := readInput(os.Stdin)
	fmt.Println("Part 1:", part1(g))
	fmt.Println("Part 2:", part2(g))
}

func part1(g graph) (total int) {
	for c := range g.nodes {
		for b, as := range g.edges {
			for a := range as {
				if !g.isConnected(a, c) {
					continue
				}

				if !g.isConnected(b, c) {
					continue
				}

				if strings.HasPrefix(a, "t") || strings.HasPrefix(b, "t") || strings.HasPrefix(c, "t") {
					total++
				}
			}
		}
	}

	return total / 6
}

func part2(g graph) (password string) {
	cliques := maximalCliques(g)

	maxClique := slices.MaxFunc(cliques, func(a, b set.Set[string]) int {
		return a.Len() - b.Len()
	})

	var computers []string
	for c := range maxClique {
		computers = append(computers, c)
	}

	slices.Sort(computers)
	return strings.Join(computers, ",")
}

func readInput(r io.Reader) (g graph) {
	g.edges = make(map[string]set.Set[string])
	g.nodes = set.New[string]()

	s := bufio.NewScanner(r)
	for s.Scan() {
		tokens := strings.SplitN(s.Text(), "-", 2)
		g.connect(tokens[0], tokens[1])
	}

	return
}

// Return all maximal cliques in the graph. A maximal clique is one that is not
// contained in some other clique. The algorithm proceeds recursively by
// removing each node from the graph in turn, finding all maximal cliques in
// the reduced graph, and then adding the node back to those cliques, if
// possible.
func maximalCliques(g graph) (cliques []set.Set[string]) {
	clique := set.New[string]()

	var bronKerbosh func(p, x set.Set[string])
	bronKerbosh = func(p, x set.Set[string]) {
		if p.IsEmpty() && x.IsEmpty() {
			cliques = append(cliques, clique.Copy())
			clear(clique)
			return
		}

		for v := range p {
			clique.Add(v)
			bronKerbosh(p.Intersect(g.edges[v]), x.Intersect(g.edges[v]))
			clique.Remove(v)
			p.Remove(v)
			x.Add(v)
		}
	}

	bronKerbosh(g.nodes.Copy(), set.New[string]())
	return
}

func (g graph) connect(a, b string) {
	if _, ok := g.edges[a]; !ok {
		g.edges[a] = set.New[string]()
	}

	if _, ok := g.edges[b]; !ok {
		g.edges[b] = set.New[string]()
	}

	g.edges[a].Add(b)
	g.edges[b].Add(a)
	g.nodes.Add(a)
	g.nodes.Add(b)
}

func (g graph) isConnected(a, b string) bool {
	to, ok := g.edges[a]
	return ok && to.Contains(b)
}
