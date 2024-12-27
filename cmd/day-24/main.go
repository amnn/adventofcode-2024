package main

import (
	"fmt"
	"internal/set"
	"io"
	"os"
	"slices"
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
	part2(n)
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

func part2(n network) {
	r := make(map[string]string)
	r["bmn"] = "z23"
	r["jss"] = "rds"
	r["mvb"] = "z08"
	r["rds"] = "jss"
	r["wss"] = "z18"
	r["z08"] = "mvb"
	r["z18"] = "wss"
	r["z23"] = "bmn"

	m, rename := adderRename(n.rewired(r))
	fmt.Println(rename)
	fmt.Println(m)
}

func adderRename(n network) (network, map[string]string) {
	r := make(map[string]string)

	// Identify the sum and carry nodes from the initial half adders:
	//
	//   x(i) ^ y(i) -> s(i)
	//   x(i) & y(i) -> c(i)
	for k, v := range n {
		if len(v.inputs) != 2 {
			continue
		}

		var i int
		slices.Sort(v.inputs)
		if k[:1] == "z" || v.inputs[0][:1] != "x" || v.inputs[1][:1] != "y" {
			continue
		} else if xd, xe := strconv.Atoi(v.inputs[0][1:]); xe != nil {
			continue
		} else if yd, ye := strconv.Atoi(v.inputs[1][1:]); ye != nil {
			continue
		} else if xd != yd {
			continue
		} else {
			i = xd
		}

		switch v.gate {
		case XOR:
			r[k] = fmt.Sprintf("s%02d", i)
		case AND:
			r[k] = fmt.Sprintf("c%02d", i)
		}
	}

	m := n.renamed(r)

	// Identify the ripple carry part of the full adder:
	//
	//   c(i-1) | ... -> C(i)
	for k, v := range m {
		if len(v.inputs) != 2 {
			continue
		}

		ap, ad := v.inputs[0][:1], v.inputs[0][1:]
		bp, bd := v.inputs[1][:1], v.inputs[1][1:]

		if k[:1] == "z" || v.gate != OR {
			continue
		} else if ca, ae := strconv.Atoi(ad); ae == nil && ap == "c" {
			r[k] = fmt.Sprintf("C%02d", ca+1)
		} else if cb, be := strconv.Atoi(bd); be == nil && bp == "c" {
			r[k] = fmt.Sprintf("C%02d", cb+1)
		}
	}

	m = n.renamed(r)
	for _, v := range m {
		slices.Sort(v.inputs)
	}

	return m, r
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

func (n network) renamed(r map[string]string) network {
	rename := func(x string) string {
		if y, ok := r[x]; ok {
			return y
		} else {
			return x
		}
	}

	m := make(network)
	for k, v := range n {
		k_ := rename(k)
		m[k_] = &*v

		var deps, inputs []string

		for _, d := range v.deps {
			deps = append(deps, rename(d))
		}

		for _, i := range v.inputs {
			inputs = append(inputs, rename(i))
		}

		m[k_].deps = deps
		m[k_].inputs = inputs
	}

	return m
}

func (n network) rewired(r map[string]string) network {
	rewire := func(x string) string {
		if y, ok := r[x]; ok {
			return y
		} else {
			return x
		}
	}

	m := make(network)
	for k, v := range n {
		m[rewire(k)] = &*v
	}

	return m
}

func (n network) Format(f fmt.State, _ rune) {
	var xs, ys, zs, gs []string
	for k := range n {
		switch k[:1] {
		case "x":
			xs = append(xs, k)
		case "y":
			ys = append(ys, k)
		case "z":
			zs = append(zs, k)
		default:
			gs = append(gs, k)
		}
	}

	slices.Sort(xs)
	slices.Sort(ys)
	slices.Sort(zs)
	slices.Sort(gs)

	// Inputs
	for _, x := range xs {
		v := n[x]
		fmt.Fprintf(f, "%s: %v\n", x, v.output)
	}

	fmt.Fprintln(f)
	for _, y := range ys {
		v := n[y]
		fmt.Fprintf(f, "%s: %v\n", y, v.output)
	}

	// Intermediate gates
	fmt.Fprintln(f)
	for _, g := range gs {
		v := n[g]
		fmt.Fprintf(f, "%s %v %s -> %s = %v\n", v.inputs[0], v.gate, v.inputs[1], g, v.output)
	}

	// Outputs
	fmt.Fprintln(f)
	for _, z := range zs {
		v := n[z]
		fmt.Fprintf(f, "%s %v %s -> %s = %v\n", v.inputs[0], v.gate, v.inputs[1], z, v.output)
	}
}

func (w wire) Format(f fmt.State, _ rune) {
	switch w {
	case Z:
		fmt.Fprint(f, "Z")
	case F:
		fmt.Fprint(f, "0")
	case T:
		fmt.Fprint(f, "1")
	}
}

func (g gate) Format(f fmt.State, _ rune) {
	switch g {
	case OR:
		fmt.Fprint(f, "|")
	case AND:
		fmt.Fprint(f, "&")
	case XOR:
		fmt.Fprint(f, "^")
	}
}
