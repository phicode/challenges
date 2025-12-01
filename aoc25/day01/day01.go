package main

// https://adventofcode.com/2025/day/1

import (
	"flag"
	"fmt"

	"github.com/phicode/challenges/lib"
)

func main() {
	flag.Parse()
	lib.Timed("Part 1", ProcessPart1, "aoc25/day01/example.txt")
	lib.Timed("Part 1", ProcessPart1, "aoc25/day01/input.txt")

	lib.Timed("Part 2", ProcessPart2, "aoc25/day01/example.txt")
	lib.Timed("Part 2", ProcessPart2, "aoc25/day01/input.txt")

	//lib.Profile(1, "day01.pprof", "Part 2", ProcessPart2, "aoc25/day01/input.txt")
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
	moves []int
}

func ReadAndParseInput(name string) Input {
	lines := lib.ReadLines(name)
	return ParseInput(lines)
}

func ParseInput(lines []string) Input {
	var rv Input
	{
	}
	for _, l := range lines {
		var dir rune
		var move int
		if n, err := fmt.Sscanf(l, "%c%d", &dir, &move); n != 2 || err != nil {
			panic("invalid input: " + l + " err: " + err.Error())
		}
		if dir != 'L' && dir != 'R' {
			panic("invalid input: " + l)
		}
		if dir == 'L' {
			move = -move
		}
		rv.moves = append(rv.moves, move)
	}
	return rv
}

////////////////////////////////////////////////////////////

func SolvePart1(input Input) int {
	pos := 50
	password := 0
	for _, m := range input.moves {
		pos = resolve(pos, m)
		if pos == 0 {
			password++
		}
	}
	return password
}

func resolve(pos int, m int) int {
	pos += m
	res := pos % 100
	if res < 0 {
		res = 100 + res
	}
	return res
}

////////////////////////////////////////////////////////////

func SolvePart2(input Input) int {
	pos := 50
	password := 0
	for _, m := range input.moves {
		var zeros int
		//before := pos
		pos, zeros = resolve2(pos, m)
		//fmt.Printf("%d - %d+%d=%d, zeros=%d\n", i, before, m, pos, zeros)
		password += zeros
	}
	return password
}

func resolve2(pos int, m int) (int, int) {
	zeros := 0
	for m >= 100 {
		m -= 100
		zeros++
	}
	for m <= -100 {
		m += 100
		zeros++
	}
	if m == 0 {
		return pos, zeros
	}
	end := pos + m
	if end == 100 {
		return 0, zeros + 1
	}
	if end > 99 {
		end -= 100
		zeros++
	}
	if end < 0 {
		end += 100
		if pos == 0 {
			return end, zeros
		}
		return end, zeros + 1
	}
	if end == 0 {
		zeros++
	}
	return end, zeros
}
