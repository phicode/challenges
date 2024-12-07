package main

// https://adventofcode.com/2024/day/7

import (
	"fmt"
	"strings"

	"github.com/phicode/challenges/lib"
	"github.com/phicode/challenges/lib/assert"
)

func main() {
	lib.Timed("Part 1", ProcessPart1, "aoc24/day07/example.txt")
	lib.Timed("Part 1", ProcessPart1, "aoc24/day07/input.txt")

	lib.Timed("Part 2", ProcessPart2, "aoc24/day07/example.txt")
	lib.Timed("Part 2", ProcessPart2, "aoc24/day07/input.txt")
}

func ProcessPart1(name string) {
	fmt.Println("Part 1 input:", name)
	input := ParseInput(name)
	result := SolvePart1(input)
	fmt.Println("Result:", result)
}

func ProcessPart2(name string) {
	fmt.Println("Part 2 input:", name)
	input := ParseInput(name)
	result := SolvePart2(input)
	fmt.Println("Result:", result)
}

////////////////////////////////////////////////////////////

type Equation struct {
	Result int
	Values []int
}

type Input struct {
	Equations []Equation
}

type Op func(int, int) int

func Add(a, b int) int { return a + b }
func Mul(a, b int) int { return a * b }
func Or(a, b int) int {
	if b == 0 {
		return a * 10
	}
	rem := b
	for rem > 0 {
		a *= 10
		rem /= 10
	}
	return a + b
}

var OpsPart1 = []Op{Add, Mul}
var OpsPart2 = []Op{Add, Mul, Or}

func ParseInput(name string) Input {
	lines := lib.ReadLines(name)
	var input Input
	for _, line := range lines {
		input.Equations = append(input.Equations, ParseEquation(line))
	}
	return input
}

func ParseEquation(line string) Equation {
	parts := strings.Split(line, ":")
	assert.True(len(parts) == 2)
	result := lib.ToInt(parts[0])
	values := lib.ExtractInts(parts[1])
	return Equation{result, values}
}

////////////////////////////////////////////////////////////

func SolvePart1(input Input) int {
	total := 0
	for _, e := range input.Equations {
		if e.CanSolve(OpsPart1) {
			total += e.Result
		}
	}
	return total
}

////////////////////////////////////////////////////////////

func SolvePart2(input Input) int {
	total := 0
	for _, e := range input.Equations {
		if e.CanSolve(OpsPart2) {
			total += e.Result
		}
	}
	return total
}

////////////////////////////////////////////////////////////

func (e Equation) CanSolve(ops []Op) bool {
	acc := e.Values[0]
	for _, op := range ops {
		if e.CanSolveX(acc, 1, ops, op) {
			return true
		}
	}
	return false
}

func (e Equation) CanSolveX(acc int, idx int, ops []Op, op Op) bool {
	acc = op(acc, e.Values[idx])
	if acc > e.Result { // early termination
		return false
	}
	if idx == len(e.Values)-1 {
		return acc == e.Result
	}
	for _, op := range ops {
		if e.CanSolveX(acc, idx+1, ops, op) {
			return true
		}
	}
	return false
}
