package main

// https://adventofcode.com/2023/day/XX

import (
	"fmt"

	"git.bind.ch/phil/challenges/lib"
)

var VERBOSE = 1

func main() {
	ProcessPart1("aoc23/day14/example.txt")
	ProcessPart1("aoc23/day14/input.txt")

	ProcessPart2("aoc23/day14/example.txt")
	ProcessPart2("aoc23/day14/input.txt")
}

func ProcessPart1(name string) {
	fmt.Println("Part 1 input:", name)
	lines := lib.ReadLines(name)
	g := ParseInput(lines)

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

type Grid lib.Grid[byte]

func ParseInput(lines []string) *Grid {
	return (*Grid)(lib.NewByteGridFromStrings(lines))
}

func (g *Grid) MoveUp() {
	for i := 1; i < g.(*lib.Grid[byte]).Rows(); i++ {

	}
}
