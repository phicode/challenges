package main

// https://adventofcode.com/2024/day/7

import (
	"fmt"
	"strings"

	"git.bind.ch/phil/challenges/lib"
	"git.bind.ch/phil/challenges/lib/assert"
)

var VERBOSE = 1

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

func log(v int, msg string) {
	if v <= VERBOSE {
		fmt.Println(msg)
	}
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
		if e.CanSolve() {
			total += e.Result
		}
	}
	return total
}

func (e Equation) CanSolve() bool {
	acc := e.Values[0]
	if e.CanSolveX(acc, 1, Add) {
		return true
	}
	if e.CanSolveX(acc, 1, Mul) {
		return true
	}
	return false
}

func (e Equation) CanSolveX(acc int, idx int, op Op) bool {
	acc = op(acc, e.Values[idx])
	if acc > e.Result { // early termination
		return false
	}
	if idx == len(e.Values)-1 {
		return acc == e.Result
	}
	if e.CanSolveX(acc, idx+1, Add) {
		return true
	}
	if e.CanSolveX(acc, idx+1, Mul) {
		return true
	}
	return false
}

////////////////////////////////////////////////////////////

func SolvePart2(input Input) int {
	return 0
}
