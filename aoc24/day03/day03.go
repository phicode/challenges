package main

// https://adventofcode.com/2024/day/3

import (
	"fmt"
	"strconv"
	"strings"

	"git.bind.ch/phil/challenges/lib"
)

// TODO: timing boilerplate
var VERBOSE = 1

func main() {
	ProcessPart1("aoc24/day03/example.txt")
	ProcessPart1("aoc24/day03/input.txt")

	ProcessPart2("aoc24/day03/example.txt")
	ProcessPart2("aoc24/day03/input.txt")
}

func ProcessPart1(name string) {
	fmt.Println("Part 1 input:", name)
	input := ParseInput(name)
	total := SolvePart1(input)
	fmt.Println("Total:", total)
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

type Input struct {
	Muls []Mul
}

type Mul struct {
	X, Y int
}

func ParseInput(name string) Input {
	lines := lib.ReadLines(name)
	var rv Input
	for _, line := range lines {
		for len(line) > 0 {
			x, y, rem, ok := parseMul(line)
			if ok {
				line = rem
				rv.Muls = append(rv.Muls, Mul{x, y})
			} else {
				line = line[1:]
			}
		}
	}
	return rv
}

func parseMul(s string) (x, y int, rem string, ok bool) {
	if !strings.HasPrefix(s, "mul(") {
		return 0, 0, "", false
	}
	end := strings.IndexRune(s, ')')
	if end == -1 || end-4 < 3 {
		return 0, 0, "", false
	}
	between := s[4:end]
	parts := strings.Split(between, ",")
	if len(parts) != 2 {
		return 0, 0, "", false
	}
	x, err1 := strconv.Atoi(parts[0])
	y, err2 := strconv.Atoi(parts[1])
	if err1 == nil && err2 == nil {
		return x, y, s[end+1:], true
	}
	return 0, 0, "", false
}

func SolvePart1(input Input) int {
	total := 0
	for _, m := range input.Muls {
		//fmt.Println(m.X, "*", m.Y)
		total += m.X * m.Y
	}
	return total
}

////////////////////////////////////////////////////////////

func ParseInputPart2(name string) Input {
	lines := lib.ReadLines(name)
	var rv Input
	for _, line := range lines {
		for len(line) > 0 {

			x, y, rem, ok := parseMul(line)
			if ok {
				line = rem
				rv.Muls = append(rv.Muls, Mul{x, y})
			} else {
				line = line[1:]
			}
		}
	}
	return rv
}
