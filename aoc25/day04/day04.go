package main

// https://adventofcode.com/2025/day/4

import (
	"flag"
	"fmt"

	"github.com/phicode/challenges/lib"
	"github.com/phicode/challenges/lib/rowcol"
)

func main() {
	flag.Parse()
	lib.Timed("Part 1", ProcessPart1, "aoc25/day04/example.txt")
	lib.Timed("Part 1", ProcessPart1, "aoc25/day04/input.txt")

	lib.Timed("Part 2", ProcessPart2, "aoc25/day04/example.txt")
	lib.Timed("Part 2", ProcessPart2, "aoc25/day04/input.txt")

	//lib.Profile(1, "day04.pprof", "Part 2", ProcessPart2, "aoc25/day04/input.txt")
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

const TOILET_PAPER = '@'

type Input struct {
	grid rowcol.Grid[byte]
}

func ReadAndParseInput(name string) Input {
	lines := lib.ReadLines(name)
	return ParseInput(lines)
}

func ParseInput(lines []string) Input {
	g := rowcol.NewByteGridFromStrings(lines)
	return Input{grid: g}
}

////////////////////////////////////////////////////////////

func SolvePart1(input Input) int {
	accessible := 0
	for pos, v := range input.grid.Iterator() {
		if v == TOILET_PAPER {
			adjacent := numAdjacentRolls(input.grid, pos)
			if adjacent < 4 {
				accessible++
			}
		}
	}
	return accessible
}

var directions = []rowcol.Pos{
	rowcol.Pos{-1, -1},
	rowcol.Pos{-1, 0},
	rowcol.Pos{-1, 1},
	rowcol.Pos{0, -1},
	rowcol.Pos{0, 1},
	rowcol.Pos{1, -1},
	rowcol.Pos{1, 0},
	rowcol.Pos{1, 1},
}

func numAdjacentRolls(g rowcol.Grid[byte], pos rowcol.Pos) int {
	n := 0
	for _, dir := range directions {
		p := pos.Add(dir)
		// fmt.Println("checking:", p)
		if g.IsValidPos(p) {
			if g.GetPos(p) == TOILET_PAPER {
				//	fmt.Println("FOUND:", p)
				n++
			}
		}
	}
	return n
}

////////////////////////////////////////////////////////////

func SolvePart2(input Input) int {
	accessible := make(map[rowcol.Pos]bool)
	for pos, v := range input.grid.Iterator() {
		if v == TOILET_PAPER {
			adjacent := numAdjacentRolls(input.grid, pos)
			if adjacent < 4 {
				accessible[pos] = true
			}
		}
	}
	totalAccessible := len(accessible)
	for len(accessible) > 0 {
		//		fmt.Printf("found %d accessible\n", len(accessible))
		// remove current accessible
		for pos := range accessible {
			input.grid.SetPos(pos, '.')
		}
		accessible = follow(input.grid, accessible)
		totalAccessible += len(accessible)
	}
	return totalAccessible
}

func follow(grid rowcol.Grid[byte], accessible map[rowcol.Pos]bool) map[rowcol.Pos]bool {
	next := make(map[rowcol.Pos]bool)
	for pos := range accessible {
		for _, dir := range directions {
			p := pos.Add(dir)
			if grid.IsValidPos(p) && grid.GetPos(p) == TOILET_PAPER {
				numAccessible := numAdjacentRolls(grid, p)
				if numAccessible < 4 {
					next[p] = true
				}
			}
		}
	}
	return next
}
