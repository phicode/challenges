package main

// https://adventofcode.com/2024/day/2

import (
	"flag"
	"fmt"

	"github.com/phicode/challenges/lib"
)

const (
	MinDiff = 1
	MaxDiff = 3
)

func main() {
	flag.Parse()
	lib.Timed("Part 1", ProcessPart1, "aoc24/day02/example.txt")
	lib.Timed("Part 1", ProcessPart1, "aoc24/day02/input.txt")

	lib.Timed("Part 2", ProcessPart2, "aoc24/day02/example.txt")
	lib.Timed("Part 2", ProcessPart2, "aoc24/day02/input.txt")
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
	ds := differences(r.Levels)
	return safeP1(ds)
}

func safeP1(diffs []int) bool {
	s := sign(diffs[0])
	for i := 0; i < len(diffs); i++ {
		d := diffs[i]
		ad := lib.AbsInt(d)
		if sign(diffs[i]) != s || ad < MinDiff || ad > MaxDiff {
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

func SolvePart2(input Input) int {
	safe := 0
	for _, report := range input.Reports {
		if report.IsSafe() {
			safe++
			continue
		}
		if report.IsSafeP2Naive() {
			safe++
		}
	}
	return safe
}

// dumb implementation that simply checks every combination
func (r Report) IsSafeP2Naive() bool {
	l := len(r.Levels)
	cpy := make([]int, l-1)
	r2 := Report{Levels: cpy}
	// test each combination
	for excludeIndex := 0; excludeIndex < l; excludeIndex++ {
		cpyIdx := 0
		for i := 0; i < l; i++ {
			if i != excludeIndex {
				cpy[cpyIdx] = r.Levels[i]
				cpyIdx++
			}
		}
		if r2.IsSafe() {
			return true
		}
	}
	return false
}

func (r Report) IsSafeP2() bool {
	ds := differences(r.Levels)

	// check if a sequence of difference numbers can be made "safe"
	// ie: adding two diffs (a, b) results in a new diff (c) which is checked against the rest of the differences
	// sequence: w, x, y, z
	// checks:
	//  - (w+x), y, z
	//  - w, (x+y), z
	//  - w, x, (y+z)
	l := len(ds) - 1
	for i := 0; i < l; i++ {
		c := ds[i] + ds[i+1]
		ac := lib.AbsInt(c)
		if ac < MinDiff || ac > MaxDiff {
			continue
		}
		sc := sign(c)
		if !checkSign(ds, sc, i) {
			continue
		}
		// check the rest of the diffs for range compliance
		for j := i + 2; j < l; j++ {
		}
	}
	return false
}

func checkSign(ds []int, sc int, excludeIndex int) bool {
	for i, d := range ds {
		if i == excludeIndex || i+1 == excludeIndex {
			continue
		}
		if sign(d) != sc {
			return false
		}
	}
	return true
}
