package main

// https://adventofcode.com/2022/day/XX

import (
	"fmt"

	"git.bind.ch/phil/challenges/lib"
)

func main() {
	Process("aoc23/dayXX/example.txt")
	Process("aoc23/dayXX/input.txt")
}

func Process(name string) {
	fmt.Println("input:", name)
	lines := lib.ReadLines(name)
	_ = lines

	fmt.Println()
}

////////////////////////////////////////////////////////////
