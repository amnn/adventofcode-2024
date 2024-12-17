package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type inst byte

const (
	ADV inst = iota
	BXL
	BST
	JNZ
	BXC
	OUT
	BDV
	CDV
)

type vm struct {
	a, b, c int
	ops     []byte
}

func main() {
	m := readInput(os.Stdin)
	m.run(os.Stdout)
}

func readInput(r io.Reader) (m vm) {
	fmt.Fscanf(r, "Register A: %d\n", &m.a)
	fmt.Fscanf(r, "Register B: %d\n", &m.b)
	fmt.Fscanf(r, "Register C: %d\n", &m.c)
	fmt.Fscanf(r, "\n")

	var ops string
	fmt.Fscanf(r, "Program: %s\n", &ops)

	for _, token := range strings.Split(ops, ",") {
		op, err := strconv.ParseInt(token, 10, 8)
		if err != nil {
			panic(err)
		}

		m.ops = append(m.ops, byte(op))
	}

	return
}

func (m *vm) run(w io.Writer) {
	combo := func(rand byte) int {
		switch rand {
		case 0, 1, 2, 3:
			return int(rand)
		case 4:
			return m.a
		case 5:
			return m.b
		case 6:
			return m.c
		default:
			panic("unexpected combo operand")
		}
	}

	for pc := 0; pc < len(m.ops); {
		switch inst(m.ops[pc]) {
		case ADV:
			m.a >>= combo(m.ops[pc+1])
			pc += 2
		case BXL:
			m.b ^= int(m.ops[pc+1])
			pc += 2
		case BST:
			m.b = combo(m.ops[pc+1]) % 8
			pc += 2
		case JNZ:
			if m.a == 0 {
				pc += 2
			} else {
				pc = int(m.ops[pc+1])
			}
		case BXC:
			m.b ^= m.c
			pc += 2
		case OUT:
			fmt.Fprintf(w, "%d,", combo(m.ops[pc+1])%8)
			pc += 2
		case BDV:
			m.b = m.a >> combo(m.ops[pc+1])
			pc += 2
		case CDV:
			m.c = m.a >> combo(m.ops[pc+1])
			pc += 2
		}
	}

	fmt.Fprintf(w, " (A: %d, B %d, C %d)\n", m.a, m.b, m.c)
}
