package main

// https://adventofcode.com/2023/day/XX

import (
	"fmt"

	"git.bind.ch/phil/challenges/lib"
)

// TODO: timing boilerplate
var VERBOSE = 1

func main() {
	ProcessPart1("aoc23/dayXX/example.txt")
	ProcessPart1("aoc23/dayXX/input.txt")

	ProcessProcessPart2("aoc23/dayXX/example.txt")
	ProcessProcessPart2("aoc23/dayXX/input.txt")
}

func ProcessPart1(name string) {
	fmt.Println("Part 1 input:", name)
	lines := lib.ReadLines(name)
	_ = lines

	fmt.Println()
}

func ProcessProcessPart2(name string) {
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
