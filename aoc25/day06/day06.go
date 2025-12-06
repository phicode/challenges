package main

// https://adventofcode.com/2025/day/6

import (
	"flag"
	"fmt"
	"strings"

	"github.com/phicode/challenges/lib"
)

func main() {
	flag.Parse()
	lib.Timed("Part 1", ProcessPart1, "aoc25/day06/example.txt")
	lib.Timed("Part 1", ProcessPart1, "aoc25/day06/input.txt")

	lib.Timed("Part 2", ProcessPart2, "aoc25/day06/example.txt")
	lib.Timed("Part 2", ProcessPart2, "aoc25/day06/input.txt")

	//lib.Profile(1, "day06.pprof", "Part 2", ProcessPart2, "aoc25/day06/input.txt")
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

type Problem struct {
	nums []int
	op   Operation
}

var operations = map[string]Operation{
	"*": Operation('*'),
	"+": Operation('+'),
}

type Operation rune

func (op Operation) Start() int {
	switch op {
	case '*':
		return 1
	case '+':
		return 0
	default:
		panic(fmt.Errorf("invalid operation: %c", op))
	}
}
func (op Operation) Acc(a, b int) int {
	switch op {
	case '*':
		return a * b
	case '+':
		return a + b
	default:
		panic(fmt.Errorf("invalid operation: %c", op))
	}
}

type Input struct {
	problems []Problem
}

func ReadAndParseInput(name string) Input {
	lines := lib.ReadLines(name)
	return ParseInput(lines)
}

func ParseInput(lines []string) Input {
	// remove empty end lines
	for l := len(lines); l > 0 && lines[l-1] == ""; {
		lines = lines[:l-1]
	}
	ops := lines[len(lines)-1]
	lines = lines[:len(lines)-1]
	allNums := make([][]int, 0, len(lines))
	for _, line := range lines {
		nums := lib.ExtractInts(line)
		allNums = append(allNums, nums)
	}
	allOps := strings.Fields(ops)
	n := len(allOps)
	for i, nums := range allNums {
		if len(nums) != n {
			panic(fmt.Errorf("number on line %d have len=%d instead of %d", i, len(nums), n))
		}
	}
	var rv Input
	for i := 0; i < n; i++ {
		var p Problem
		p.nums = make([]int, len(allNums))
		op, found := operations[allOps[i]]
		if !found {
			panic(fmt.Errorf("no operation found for input %q", allOps[i]))
		}
		p.op = op
		for j, nums := range allNums {
			p.nums[j] = nums[i]
		}
		rv.problems = append(rv.problems, p)
	}
	return rv
}

func (p Problem) Solve() int {
	acc := p.op.Start()
	for _, num := range p.nums {
		acc = p.op.Acc(acc, num)
	}
	return acc
}

////////////////////////////////////////////////////////////

func SolvePart1(input Input) int {
	total := 0
	for _, problem := range input.problems {
		n := problem.Solve()
		total += n
	}
	return total
}

////////////////////////////////////////////////////////////

func SolvePart2(input Input) int {
	return 0
}
