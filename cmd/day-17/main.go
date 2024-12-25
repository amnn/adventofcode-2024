package main

import (
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
)

type inst byte

const (
	ADV inst = iota // 0
	BXL             // 1
	BST             // 2
	JNZ             // 3
	BXC             // 4
	OUT             // 5
	BDV             // 6
	CDV             // 7
)

type vm struct {
	a, b, c, pc int
	ops, out    []byte
}

func main() {
	m := readInput(os.Stdin)

	fmt.Println(&m)
	fmt.Println("Part 1:", part1(m))
	// fmt.Println("Part 2:", part2(m))
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

func part1(m vm) string {
	for m.step() {
	}

	var tokens []string
	for _, b := range m.out {
		tokens = append(tokens, strconv.Itoa(int(b)))
	}

	return strings.Join(tokens, ",")
}

func part2(m vm) int {
outer:
	for a := 0; ; a++ {
		copy := m
		copy.a = a

		for copy.step() {
			if len(copy.out) > len(copy.ops) {
				continue outer
			}

			if !slices.Equal(copy.out, copy.ops[:len(copy.out)]) {
				continue outer
			}
		}

		if slices.Equal(copy.out, copy.ops) {
			return a
		}
	}
}

func (m *vm) step() bool {
	if m.pc >= len(m.ops) {
		return false
	}

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

	switch inst(m.ops[m.pc]) {
	case ADV:
		m.a >>= combo(m.ops[m.pc+1])
		m.pc += 2
	case BXL:
		m.b ^= int(m.ops[m.pc+1])
		m.pc += 2
	case BST:
		m.b = combo(m.ops[m.pc+1]) % 8
		m.pc += 2
	case JNZ:
		if m.a == 0 {
			m.pc += 2
		} else {
			m.pc = int(m.ops[m.pc+1])
		}
	case BXC:
		m.b ^= m.c
		m.pc += 2
	case OUT:
		m.out = append(m.out, byte(combo(m.ops[m.pc+1])%8))
		m.pc += 2
	case BDV:
		m.b = m.a >> combo(m.ops[m.pc+1])
		m.pc += 2
	case CDV:
		m.c = m.a >> combo(m.ops[m.pc+1])
		m.pc += 2
	}

	return true
}

func (m *vm) Format(f fmt.State, _ rune) {
	fmt.Fprintf(f, "A: %d\nB: %d\nC: %d\n\n", m.a, m.b, m.c)

	jumps := make(map[int]int)
	for label, i := 0, 0; i < len(m.ops); i += 2 {
		if inst(m.ops[i]) == JNZ {
			jumps[int(m.ops[i+1])] = label
			label++
		}
	}

	for i := 0; i < len(m.ops); i += 2 {
		if label, ok := jumps[i]; ok {
			fmt.Fprintf(f, "L%d:\n", label)
		}

		fmt.Fprintf(f, "%04x: ", i)
		switch inst(m.ops[i]) {
		case ADV:
			fmt.Fprint(f, "A >>= ")
			formatAsComboOperand(m.ops[i+1], f)
		case BXL:
			fmt.Fprintf(f, "B ^= %03b", m.ops[i+1])
		case BST:
			fmt.Fprint(f, "B := ")
			formatAsComboOperand(m.ops[i+1], f)
			fmt.Fprint(f, " % 8")
		case JNZ:
			fmt.Fprintf(f, "if A != 0 goto L%d", jumps[int(m.ops[i+1])])
		case BXC:
			fmt.Fprint(f, "B ^= C")
		case OUT:
			fmt.Fprint(f, "output ")
			formatAsComboOperand(m.ops[i+1], f)
			fmt.Fprint(f, " % 8")
		case BDV:
			fmt.Fprint(f, "B := A >> ")
			formatAsComboOperand(m.ops[i+1], f)
		case CDV:
			fmt.Fprint(f, "C := A >> ")
			formatAsComboOperand(m.ops[i+1], f)
		}

		fmt.Fprintln(f, "")
	}
}

func formatAsComboOperand(rand byte, f fmt.State) {
	switch rand {
	case 0, 1, 2, 3:
		fmt.Fprintf(f, "%d", rand)
	case 4:
		fmt.Fprintf(f, "A")
	case 5:
		fmt.Fprintf(f, "B")
	case 6:
		fmt.Fprintf(f, "C")
	default:
		panic("unexpected combo operand")
	}

	return
}
