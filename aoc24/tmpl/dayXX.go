package main

// https://adventofcode.com/2024/day/XX

import (
	"flag"
	"fmt"

	"github.com/phicode/challenges/lib"
)

func main() {
	flag.Parse()
	lib.Timed("Part 1", ProcessPart1, "aoc24/dayXX/example.txt")
	lib.Timed("Part 1", ProcessPart1, "aoc24/dayXX/input.txt")

	lib.Timed("Part 2", ProcessPart2, "aoc24/dayXX/example.txt")
	lib.Timed("Part 2", ProcessPart2, "aoc24/dayXX/input.txt")

	//lib.Profile(1, "dayXX.pprof", "Part 2", ProcessPart2, "aoc24/dayXX/input.txt")
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
