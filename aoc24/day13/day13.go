package main

// https://adventofcode.com/2024/day/13

import (
	"flag"
	"fmt"
	"strings"

	"github.com/phicode/challenges/lib"
	"github.com/phicode/challenges/lib/assert"
	"github.com/phicode/challenges/lib/math"
)

func main() {
	flag.Parse()
	lib.Timed("Part 1", ProcessPart1, "aoc24/day13/example.txt")
	lib.Timed("Part 1", ProcessPart1, "aoc24/day13/input.txt")
	lib.Timed("Part 2", ProcessPart2, "aoc24/day13/example.txt")
	lib.Timed("Part 2", ProcessPart2, "aoc24/day13/input.txt")
	//lib.Profile(1, "day13.pprof", "Part 2", ProcessPart2, "aoc24/day13/input.txt")
}

func ProcessPart1(name string) {
	fmt.Println("Part 1 input:", name)
	input := ReadAndParseInput(name)
	result := SolvePart1(input)
	fmt.Println("Result:", result)
}

func ProcessPart2(name string) {
	fmt.Println("Part 2 input:", name)
	input := ReadAndParseInput(name)
	result := SolvePart2(input)
	fmt.Println("Result:", result)
}

////////////////////////////////////////////////////////////

type Input struct {
	machines []Machine
}

type Machine struct {
	A, B  Offset
	Prize Offset
}

type Offset struct {
	X, Y int
}

func ReadAndParseInput(name string) Input {
	lines := lib.ReadLines(name)
	return ParseInput(lines)
}

func ParseInput(lines []string) Input {
	var machines []Machine
	var current Machine
	for _, line := range lines {
		if strings.HasPrefix(line, "Button A") {
			current = Machine{}
			current.A = ParseButton(line)
		} else if strings.HasPrefix(line, "Button B") {
			current.B = ParseButton(line)
		} else if strings.HasPrefix(line, "Prize:") {
			var x, y int
			n, err := fmt.Sscanf(line, "Prize: X=%d, Y=%d", &x, &y)
			assert.True(n == 2 && err == nil)
			current.Prize = Offset{x, y}
			machines = append(machines, current)
			current = Machine{}
		} else {
			assert.True(len(line) == 0)
		}
	}
	return Input{machines}
}

func ParseButton(line string) Offset {
	var x, y int
	// skip 'Button A: ' or 'Button B: '
	line = line[10:]
	n, err := fmt.Sscanf(line, "X+%d, Y+%d", &x, &y)
	assert.True(n == 2 && err == nil)
	return Offset{x, y}
}

////////////////////////////////////////////////////////////

func SolvePart1(input Input) int {
	total := 0
	for _, m := range input.machines {
		cost, ok := m.Solve()
		if ok {
			total += cost
		}
	}
	return total
}

////////////////////////////////////////////////////////////

func SolvePart2(input Input) int {
	total := 0
	for _, m := range input.machines {
		cost, ok := m.AdjustPart2().Solve()
		if ok {
			total += cost
		}
	}
	return total
}

func (m Machine) AdjustPart2() Machine {
	cpy := m
	cpy.Prize.X += 10000000000000
	cpy.Prize.Y += 10000000000000
	return cpy
}

func (m Machine) Solve() (int, bool) {
	// eq1: Xa  Xb  | Xp
	// eq2: Ya  Yb  | Yp
	lcm := math.Lcm(m.A.X, m.A.Y)

	// scale equations so that: Xa==Ya
	s := m.Mul(lcm/m.A.X, lcm/m.A.Y)
	assert.True(s.A.X == s.A.Y)

	// Subtract eq1 from eq2
	s.A.Y -= s.A.X
	s.B.Y -= s.B.X
	s.Prize.Y -= s.Prize.X

	b := s.Prize.Y / s.B.Y
	if b*s.B.Y != s.Prize.Y {
		return 0, false
	}
	x := s.Prize.X - s.B.X*b
	a := x / s.A.X
	if a*s.A.X != x {
		return 0, false
	}
	return a*3 + b, true
}

func (m Machine) Mul(x int, y int) Machine {
	return Machine{
		A:     Offset{m.A.X * x, m.A.Y * y},
		B:     Offset{m.B.X * x, m.B.Y * y},
		Prize: Offset{m.Prize.X * x, m.Prize.Y * y},
	}
}
