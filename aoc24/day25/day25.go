package main

// https://adventofcode.com/2024/day/25

import (
	"flag"
	"fmt"

	"github.com/phicode/challenges/lib"
	"github.com/phicode/challenges/lib/assert"
)

func main() {
	flag.Parse()
	lib.Timed("Part 1", ProcessPart1, "aoc24/day25/example.txt")
	lib.Timed("Part 1", ProcessPart1, "aoc24/day25/input.txt")
	//
	//lib.Timed("Part 2", ProcessPart2, "aoc24/day25/example.txt")
	//lib.Timed("Part 2", ProcessPart2, "aoc24/day25/input.txt")

	//lib.Profile(1, "day25.pprof", "Part 2", ProcessPart2, "aoc24/day25/input.txt")
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
	locks []Combination
	keys  []Combination
}

type Combination struct {
	Key  bool
	Cols []int
}

func (c Combination) fits(lock Combination) bool {
	for i := 0; i < len(c.Cols); i++ {
		sum := c.Cols[i] + lock.Cols[i]
		if sum > 5 {
			return false
		}
	}
	return true
}

func ReadAndParseInput(name string) Input {
	lines := lib.ReadLines(name)
	return ParseInput(lines)
}

func ParseInput(lines []string) Input {
	var locks []Combination
	var keys []Combination
	for len(lines) > 0 {
		// 7 rows are required for a key/lock
		assert.True(len(lines) >= 7)
		c, key := ParseCombination(lines[:7])
		if key {
			keys = append(keys, c)
		} else {
			locks = append(locks, c)
		}
		if len(lines) == 7 {
			break
		}
		assert.True(lines[7] == "")
		lines = lines[8:]
	}
	return Input{keys, locks}
}

func ParseCombination(lines []string) (Combination, bool) {
	cols := make([]int, len(lines[0]))
	key := lines[0][0] == '.'
	for _, line := range lines {
		for i, c := range line {
			if c == '#' {
				cols[i]++
			}
		}
	}
	for i, _ := range cols {
		cols[i]-- // subtract first/last line
	}
	return Combination{key, cols}, key
}

////////////////////////////////////////////////////////////

func SolvePart1(input Input) int {
	n := 0
	for _, lock := range input.locks {
		for _, key := range input.keys {

			if key.fits(lock) {
				//fmt.Printf("%v fits %v\n", key, lock)
				n++
			}
		}
	}
	return n
}

////////////////////////////////////////////////////////////

func SolvePart2(input Input) int {
	return 0
}
