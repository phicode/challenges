package main

// https://adventofcode.com/2024/day/2

import (
	"fmt"

	"git.bind.ch/phil/challenges/lib"
)

// TODO: timing boilerplate
var VERBOSE = 1

func main() {
	ProcessPart1("aoc24/day02/example.txt")
	ProcessPart1("aoc24/day02/input.txt")

	ProcessPart2("aoc24/day02/example.txt")
	ProcessPart2("aoc24/day02/input.txt")
}

func ProcessPart1(name string) {
	fmt.Println("Part 1 input:", name)
	input := ParseInput(name)
	safe := SolvePart1(input)
	fmt.Println("Safe:", safe)
}

func ProcessPart2(name string) {
	fmt.Println("Part 2 input:", name)
	lines := lib.ReadLines(name)
	_ = lines

	fmt.Println()
}

func log(v int, msg string) {
	if v <= VERBOSE {
		fmt.Println(msg)
	}
}

////////////////////////////////////////////////////////////

type Input struct {
	Reports []Report
}
type Report struct {
	Levels []int
}

func ParseInput(name string) Input {
	lines := lib.ReadLines(name)
	var rv Input
	for _, line := range lines {
		rv.Reports = append(rv.Reports, Report{Levels: lib.ExtractInts(line)})
	}
	return rv
}

////////////////////////////////////////////////////////////

func SolvePart1(input Input) int {
	safe := 0
	for _, report := range input.Reports {
		if report.IsSafe() {
			safe++
		}
	}
	return safe
}

func (r Report) IsSafe() bool {
	d := differences(r.Levels)
	return sameSign(d) && withinLimits(d, 1, 3)
}

func withinLimits(xs []int, _min, _max int) bool {
	for _, x := range xs {
		v := lib.AbsInt(x)
		if v < _min || v > _max {
			return false
		}
	}
	return true
}

func sameSign(xs []int) bool {
	s := sign(xs[0])
	for i := 1; i < len(xs); i++ {
		if sign(xs[i]) != s {
			return false
		}
	}
	return true
}

func sign(x int) int {
	if x == 0 {
		return 0
	}
	if x > 0 {
		return 1
	}
	return -1
}

func differences(xs []int) []int {
	n := len(xs) - 1
	diffs := make([]int, n)
	for i := 0; i < n; i++ {
		diffs[i] = xs[i+1] - xs[i]
	}
	return diffs
}

////////////////////////////////////////////////////////////
