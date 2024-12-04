package main

// https://adventofcode.com/2024/day/1

import (
	"fmt"
	"sort"

	"git.bind.ch/phil/challenges/lib"
)

var VERBOSE = 1

func main() {
	lib.Timed("Part 1", ProcessPart1, "aoc24/day01/example.txt")
	lib.Timed("Part 1", ProcessPart1, "aoc24/day01/input.txt")

	lib.Timed("Part 2", ProcessPart2, "aoc24/day01/example.txt")
	lib.Timed("Part 2", ProcessPart2, "aoc24/day01/input.txt")
}

func ProcessPart1(name string) {
	fmt.Println("Part 1 input:", name)
	input := ParseInput(name)
	dist := SolvePart1(input)
	fmt.Println("Distance:", dist)
}

func ProcessPart2(name string) {
	fmt.Println("Part 2 input:", name)
	input := ParseInput(name)
	dist := SolvePart2(input)
	fmt.Println("Distance:", dist)
}

func log(v int, msg string) {
	if v <= VERBOSE {
		fmt.Println(msg)
	}
}

////////////////////////////////////////////////////////////

type Input struct {
	A []int
	B []int
}

func ParseInput(name string) (input Input) {
	lines := lib.ReadLines(name)
	var rv Input
	for _, line := range lines {
		ints := lib.ExtractInts(line)
		rv.A = append(rv.A, ints[0])
		rv.B = append(rv.B, ints[1])
	}
	return rv
}

func SolvePart1(input Input) int {
	sort.Ints(input.A)
	sort.Ints(input.B)
	total := 0
	for i := 0; i < len(input.A); i++ {
		dist := lib.AbsInt(input.A[i] - input.B[i])
		total += dist
	}
	return total
}

func SolvePart2(input Input) int {
	sort.Ints(input.B)
	total := 0
	for _, v := range input.A {
		c := count(v, input.B)
		total += v * c
	}
	return total
}

func count(v int, xs []int) int {
	l := len(xs)
	i := sort.SearchInts(xs, v)
	if i == l || xs[i] != v {
		return 0 // not found
	}
	n := 1
	for i+1 < l && xs[i+1] == v {
		i++
		n++
	}
	return n
}
