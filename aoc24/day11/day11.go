package main

// https://adventofcode.com/2024/day/11

import (
	"flag"
	"fmt"

	"github.com/phicode/challenges/lib"
	"github.com/phicode/challenges/lib/assert"
)

func main() {
	flag.Parse()
	lib.Timed("Part 1", ProcessPart1, "aoc24/day11/example.txt")
	lib.Timed("Part 1", ProcessPart1, "aoc24/day11/input.txt")

	lib.Timed("Part 2", ProcessPart2, "aoc24/day11/example.txt")
	lib.Timed("Part 2", ProcessPart2, "aoc24/day11/input.txt")

	//lib.Profile(1, "day11.pprof", "Part 2", ProcessPart2, "aoc24/dayXX/input.txt")
}

func ProcessPart1(name string) {
	fmt.Println("Part 1 input:", name)
	input := ReadAndParseInput(name)
	result := Solve(input, 25)
	fmt.Println("Result:", result)
}

func ProcessPart2(name string) {
	fmt.Println("Part 2 input:", name)
	input := ReadAndParseInput(name)
	result := Solve(input, 75)
	fmt.Println("Result:", result)
}

////////////////////////////////////////////////////////////

type Input struct {
	Numbers []int
}

func ReadAndParseInput(name string) Input {
	lines := lib.ReadLines(name)
	assert.True(len(lines) == 1)
	return ParseInput(lines[0])
}

func ParseInput(line string) Input {
	numbers := lib.ExtractInts(line)
	return Input{numbers}
}

////////////////////////////////////////////////////////////

func Solve(input Input, blinks int) int {
	cache := make(map[Key]int)
	sum := 0
	for _, n := range input.Numbers {
		sum += solve(cache, n, 1, blinks)
	}
	return sum
}

type Key struct {
	Number int
	Blinks int
}

func solve(cache map[Key]int, num, depth, maxDepth int) int {
	blinks := maxDepth - depth + 1
	key := Key{num, blinks}
	if v, ok := cache[key]; ok {
		return v
	}
	a, b := next(num)
	if depth == maxDepth {
		stones := 1
		if b != -1 {
			stones = 2
		}
		cache[key] = stones
		return stones
	}
	stones := solve(cache, a, depth+1, maxDepth)
	if b != -1 {
		stones += solve(cache, b, depth+1, maxDepth)
	}
	cache[key] = stones
	return stones
}

func next(num int) (int, int) {
	if num == 0 {
		return 1, -1
	}
	if n := digits(num); n%2 == 0 {
		base := 1
		a, b := num, num
		half := n / 2
		for half > 0 {
			half--
			a /= 10
			base *= 10
		}
		return a, b % base
	}
	return num * 2024, -1
}

func digits(num int) int {
	d := 0
	for num > 0 {
		num /= 10
		d++
	}
	return d
}
