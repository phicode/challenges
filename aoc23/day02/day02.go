package main

// https://adventofcode.com/2023/day/2

import (
	"fmt"

	"git.bind.ch/phil/challenges/lib"
)

func main() {
	ProcessStep1("aoc23/day02/example.txt")
	ProcessStep1("aoc23/day02/input.txt")

	ProcessStep2("aoc23/day02/example.txt")
	ProcessStep2("aoc23/day02/input.txt")
}

func ProcessStep1(name string) {
	fmt.Println("input:", name)
	lines := lib.ReadLines(name)
	_ = lines

	fmt.Println()
}

func ProcessStep2(name string) {
	fmt.Println("input:", name)
	lines := lib.ReadLines(name)
	_ = lines

	fmt.Println()
}

////////////////////////////////////////////////////////////
