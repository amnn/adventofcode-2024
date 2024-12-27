package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type wire uint8
type gate uint8

const (
	Z wire = 0b00
	F wire = 0b01
	T wire = 0b11
)

const (
	ID gate = iota
	OR
	AND
	XOR
)

type node struct {
	output wire
	gate   gate
	inputs []string
	deps   []string
}

type network map[string]*node

func main() {
	n := readInput(os.Stdin)
	n.propagate()

	fmt.Println("Part 1:", part1(n))
}

func part1(n network) (val int) {
	for k, v := range n {
		if !strings.HasPrefix(k, "z") {
			continue
		}

		i, err := strconv.ParseInt(k[1:], 10, 8)
		if err != nil {
			continue
		}

		val |= int(v.output>>1) << uint(i)
	}

	return
}

func readInput(r io.Reader) (n network) {
	n = make(map[string]*node)
	ensure := func(v string) {
		if _, ok := n[v]; !ok {
			n[v] = &node{}
		}
	}

	// Read input wires
	for {
		var v, x string
		if _, err := fmt.Fscanf(r, "%s %s\n", &v, &x); err != nil {
			break
		}

		var output wire
		switch x {
		case "0":
			output = F
		case "1":
			output = T
		default:
			panic("invalid wire")
		}

		n[strings.TrimSuffix(v, ":")] = &node{
			output: output,
			gate:   ID,
		}
	}

	// Read non-trivial gates
	for {
		var vl, op, vr, vo string
		if _, err := fmt.Fscanf(r, "%s %s %s -> %s\n", &vl, &op, &vr, &vo); err != nil {
			break
		}

		ensure(vl)
		ensure(vr)
		ensure(vo)

		switch op {
		case "AND":
			n[vo].gate = AND
		case "OR":
			n[vo].gate = OR
		case "XOR":
			n[vo].gate = XOR
		default:
			panic("invalid gate")
		}

		n[vo].inputs = append(n[vo].inputs, vl, vr)
		n[vl].deps = append(n[vl].deps, vo)
		n[vr].deps = append(n[vr].deps, vo)
	}

	return
}

func (n network) propagate() {
	var work []string
	for _, v := range n {
		if v.output != Z {
			work = append(work, v.deps...)
		}
	}

	var v string
	for len(work) > 0 {
		v, work = work[0], work[1:]

		node := n[v]
		if node.output != Z {
			continue
		}

		if len(node.inputs) != 2 {
			panic("invalid gate")
		}

		l := n[node.inputs[0]].output
		r := n[node.inputs[1]].output
		if l == Z || r == Z {
			continue
		}

		switch node.gate {
		case ID:
			panic("ID gate should have been resolved")
		case OR:
			node.output = l | r
		case AND:
			node.output = l & r
		case XOR:
			node.output = (l ^ r) | F
		}

		work = append(work, node.deps...)
	}
}
