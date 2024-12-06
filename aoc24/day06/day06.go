package main

// https://adventofcode.com/2024/day/6

import (
	"fmt"

	"git.bind.ch/phil/challenges/lib"
)

var VERBOSE = 1

func main() {
	lib.Timed("Part 1", ProcessPart1, "aoc24/day06/example.txt")
	lib.Timed("Part 1", ProcessPart1, "aoc24/day06/input.txt")

	lib.Timed("Part 2", ProcessPart2, "aoc24/day06/example.txt")
	lib.Timed("Part 2", ProcessPart2, "aoc24/day06/input.txt")
}

func ProcessPart1(name string) {
	fmt.Println("Part 1 input:", name)
	lines := lib.ReadLines(name)
	_ = lines

	fmt.Println()
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
