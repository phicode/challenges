package main

// https://adventofcode.com/2024/day/02

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
