package main

import (
	"fmt"
	"sort"
)

func main() {
	var ls, rs []int

	for {
		var l, r int
		n, err := fmt.Scanf("%d %d\n", &l, &r)
		if err != nil || n != 2 {
			break
		}

		ls = append(ls, l)
		rs = append(rs, r)
	}

	fmt.Println("Part 1:", part1(ls, rs))
	fmt.Println("Part 2:", part2(ls, rs))
}

func part1(ls, rs []int) int {
	sort.Ints(ls)
	sort.Ints(rs)

	total := 0
	for i := 0; i < len(ls); i++ {
		if ls[i] > rs[i] {
			total += ls[i] - rs[i]
		} else {
			total += rs[i] - ls[i]
		}
	}

	return total
}

func part2(ls, rs []int) int {
	seen := make(map[int]int)

	for _, v := range ls {
		seen[v] = 0
	}

	for _, v := range rs {
		if count, ok := seen[v]; ok {
			seen[v] = count + 1
		}
	}

	total := 0
	for k, v := range seen {
		total += k * v
	}

	return total
}
