package main

// https://adventofcode.com/2023/day/XX

import (
	"fmt"

	"git.bind.ch/phil/challenges/lib"
)

// TODO: debug boilerplate, timing boilerplate

func main() {
	ProcessStep1("aoc23/dayXX/example.txt")
	ProcessStep1("aoc23/dayXX/input.txt")

	ProcessStep2("aoc23/dayXX/example.txt")
	ProcessStep2("aoc23/dayXX/input.txt")
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

var DEBUG = 1

func debug(v int, msg string) {
	if v <= DEBUG {
		fmt.Println(msg)
	}
}

////////////////////////////////////////////////////////////
