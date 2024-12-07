package main

// https://adventofcode.com/2022/day/11

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/phicode/challenges/lib/assets"
)

func main() {
	Process("aoc22/day11/example.txt", Div3, 20)
	Process("aoc22/day11/input.txt", Div3, 20)
	Process("aoc22/day11/example.txt", nil, 10000)
	Process("aoc22/day11/input.txt", nil, 10000)
}

func Process(name string, wm WorryMod, rounds int) {
	fmt.Println("input:", name)
	lines := ReadInput(name)
	monkeys := ParseMonkeys(lines)
	if wm == nil {
		wm = FindAbsoluteLimit(monkeys)
	}
	for i := 1; i <= rounds; i++ {
		Round(monkeys, wm)
		if i == 1 || i == 20 || i%1000 == 0 {
			Print(i, monkeys)
		}
	}
	Print(rounds, monkeys)
	fmt.Println("monkey-business:", MonkeyBusiness(monkeys))
	fmt.Println()
}

// first real stumbling block in AOC22
// hint found in: https://github.com/ericwburden/advent_of_code_2022/blob/main/src/day11/part2.rs
func FindAbsoluteLimit(monkeys []*Monkey) WorryMod {
	product := 1
	for _, m := range monkeys {
		product *= m.TestDivisible
	}
	return func(x int) int {
		return x % product
	}
}

func MonkeyBusiness(monkeys []*Monkey) any {
	var inspects []int
	for _, m := range monkeys {
		inspects = append(inspects, m.Inspects)
	}
	sort.Ints(inspects)
	l := len(inspects)
	return inspects[l-1] * inspects[l-2]
}

func ReadInput(name string) []string {
	f, err := os.Open(assets.MustFind(name))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	s := bufio.NewScanner(f)

	var lines []string
	for s.Scan() {
		line := s.Text()
		lines = append(lines, line)
	}

	if err := s.Err(); err != nil {
		panic(s.Err())
	}
	return lines
}

////////////////////////////////////////////////////////////

type Monkey struct {
	Inspects      int
	Items         []int
	Op            func(int) int
	TestDivisible int
	ThrowTrue     int
	ThrowFalse    int
}

func ParseMonkeys(lines []string) []*Monkey {
	var mks []*Monkey
	idx := 0
	for len(lines) >= 6 {
		mks = append(mks, ParseMonkey(idx, lines))
		idx++
		lines = lines[6:]
		if len(lines) > 0 {
			if lines[0] != "" {
				panic("invalid input")
			}
			lines = lines[1:]
		} else {
		}
	}
	return mks
}

func ParseMonkey(idx int, lines []string) *Monkey {
	// line 0: Monkey 0:
	var midx int
	n, err := fmt.Sscanf(lines[0], "Monkey %d:", &midx)
	if n != 1 || err != nil || idx != midx {
		panic(fmt.Errorf("invalid monkey definition start: %q; expect index=%d", lines[0], idx))
	}

	// line 1:   Starting items: 79, 98
	sitems, found := strings.CutPrefix(lines[1], "  Starting items: ")
	if !found {
		panic(fmt.Errorf("invalid starting items: %q", lines[1]))
	}
	items := split(sitems)

	// line 2:   Operation: new = old * 19
	op := parseOp(lines[2])

	// line 3:   Test: divisible by 23
	var div int
	n, err = fmt.Sscanf(lines[3], "  Test: divisible by %d", &div)
	if n != 1 || err != nil {
		panic("invalid division line")
	}

	// line 4:     If true: throw to monkey 2
	var throwtrue int
	n, err = fmt.Sscanf(lines[4], "    If true: throw to monkey %d", &throwtrue)
	if n != 1 || err != nil {
		panic("invalid throw true line")
	}

	// line 5:     If false: throw to monkey 3
	var throwfalse int
	n, err = fmt.Sscanf(lines[5], "    If false: throw to monkey %d", &throwfalse)
	if n != 1 || err != nil {
		panic("invalid throw false line")
	}

	return &Monkey{
		Items:         items,
		Op:            op,
		TestDivisible: div,
		ThrowTrue:     throwtrue,
		ThrowFalse:    throwfalse,
	}
}

func parseOp(s string) func(int) int {
	if s == "  Operation: new = old * old" {
		return func(x int) int {
			return x * x
		}
	}

	var op rune
	var coeff int
	n, err := fmt.Sscanf(s, "  Operation: new = old %c %d", &op, &coeff)
	if n != 2 || err != nil {
		panic(fmt.Errorf("invalid operation line: %q, err=%w", s, err))
	}
	if op == '+' {
		return func(x int) int {
			return x + coeff
		}
	}
	if op == '*' {
		return func(x int) int {
			return x * coeff
		}
	}
	panic(fmt.Errorf("unknown operation: %c", op))
}

func split(s string) []int {
	parts := strings.Split(s, ", ")
	items := make([]int, len(parts))
	for i, part := range parts {
		var err error
		items[i], err = strconv.Atoi(part)
		if err != nil {
			panic(fmt.Errorf("failed to parse item number: %q", part))
		}
	}
	return items
}

func (m *Monkey) Round(all []*Monkey, wm WorryMod) {
	items := m.Items
	m.Items = m.Items[:0] // throw all items away
	m.Inspects += len(items)
	for _, item := range items {
		item = m.Op(item)
		item = wm(item)
		var receiver *Monkey
		if item%m.TestDivisible == 0 {
			receiver = all[m.ThrowTrue]
		} else {
			receiver = all[m.ThrowFalse]
		}
		receiver.Items = append(receiver.Items, item)
	}
}

func Round(all []*Monkey, wm WorryMod) {
	for _, m := range all {
		m.Round(all, wm)
	}
}

func Print(round int, monkeys []*Monkey) {
	fmt.Printf("=== After round %d ===\n", round)
	for i, m := range monkeys {
		fmt.Printf("Monkey %d: %v\n", i, m.Items)
	}
	for i, m := range monkeys {
		fmt.Printf("Monkey %d inspected items %d times.\n", i, m.Inspects)
	}
	fmt.Println()
}

type WorryMod func(int) int

func NoChange(x int) int {
	return x
}

func Div3(x int) int {
	return x / 3
}
