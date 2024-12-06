package main

// https://adventofcode.com/2024/day/XX

import (
	"fmt"

	"git.bind.ch/phil/challenges/lib"
	"git.bind.ch/phil/challenges/lib/rowcol"
)

var VERBOSE = 1

func main() {
	lib.Timed("Part 1", ProcessPart1, "aoc24/dayXX/example.txt")
	lib.Timed("Part 1", ProcessPart1, "aoc24/dayXX/input.txt")

	lib.Timed("Part 2", ProcessPart2, "aoc24/dayXX/example.txt")
	lib.Timed("Part 2", ProcessPart2, "aoc24/dayXX/input.txt")
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

type Grid struct {
	rowcol.Grid[byte]
}

func ParseInput(name string) Input {
	lines := lib.ReadLines(name)
	_ = lines
	return Input{}
}

////////////////////////////////////////////////////////////

func SolvePart1(input Input) int {
	return 0
}

////////////////////////////////////////////////////////////

func SolvePart2(input Input) int {
	return 0
}
