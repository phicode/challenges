package main

// https://adventofcode.com/2024/day/19

import (
	"bytes"
	"flag"
	"fmt"
	"strings"

	"github.com/phicode/challenges/lib"
	"github.com/phicode/challenges/lib/assert"
)

func main() {
	flag.Parse()
	lib.Timed("Part 1", ProcessPart1, "aoc24/day19/example.txt")
	lib.Timed("Part 1", ProcessPart1, "aoc24/day19/input.txt")

	lib.Timed("Part 2", ProcessPart2, "aoc24/day19/example.txt")
	lib.Timed("Part 2", ProcessPart2, "aoc24/day19/input.txt")

	//lib.Profile(1, "day19.pprof", "Part 2", ProcessPart2, "aoc24/day19/input.txt")
}

func ProcessPart1(name string) {
	fmt.Println("Part 1 input:", name)
	input := ReadAndParseInput(name)
	result := SolvePart1(input)
	fmt.Println("Result:", result)
}

func ProcessPart2(name string) {
	fmt.Println("Part 2 input:", name)
	input := ReadAndParseInput(name)
	result := SolvePart2(input)
	fmt.Println("Result:", result)
}

////////////////////////////////////////////////////////////

type Input struct {
	Patterns [][]byte
	Designs  [][]byte
}

func ReadAndParseInput(name string) Input {
	lines := lib.ReadLines(name)
	return ParseInput(lines)
}

func ToBytes(s string) []byte { return []byte(s) }

func ParseInput(lines []string) Input {
	parts := strings.Split(lines[0], ", ")
	patterns := lib.Map(parts, ToBytes)
	assert.True(len(lines[1]) == 0)
	designs := lib.Map(lines[2:], ToBytes)
	return Input{
		Patterns: patterns,
		Designs:  designs,
	}
}

////////////////////////////////////////////////////////////

func SolvePart1(input Input) int {
	total := 0
	for _, design := range input.Designs {
		if IsPossible(design, input.Patterns) {
			//fmt.Println("Design", string(design), "is possible")
			total++
		}
	}
	return total
}

func IsPossible(design []byte, patterns [][]byte) bool {
	for _, pattern := range patterns {
		ld := len(design)
		lp := len(pattern)
		if lp > ld {
			continue
		}
		if lp == ld && bytes.Equal(design, pattern) {
			return true
		}
		if lp < ld && bytes.Equal(design[:lp], pattern) {
			if IsPossible(design[lp:], patterns) {
				return true
			}
		}
	}
	return false
}

////////////////////////////////////////////////////////////

func SolvePart2(input Input) int {
	return 0
}
