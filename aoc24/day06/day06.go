package main

// https://adventofcode.com/2024/day/6

import (
	"fmt"

	"git.bind.ch/phil/challenges/lib"
	"git.bind.ch/phil/challenges/lib/assert"
	"git.bind.ch/phil/challenges/lib/rowcol"
)

var VERBOSE = 1

func main() {
	lib.Timed("Part 1", ProcessPart1, "aoc24/day06/example.txt")
	lib.Timed("Part 1", ProcessPart1, "aoc24/day06/input.txt")

	lib.Timed("Part 2", ProcessPart2, "aoc24/day06/example.txt")
	lib.Timed("Part 2", ProcessPart2, "aoc24/day06/input.txt")
}

func ProcessPart1(name string) {
	fmt.Println("Part 1 input:", name)
	grid := ParseInput(name)
	result := SolvePart1(grid)
	fmt.Println("Result:", result)
}

func ProcessPart2(name string) {
	fmt.Println("Part 2 input:", name)
	grid := ParseInput(name)
	result := SolvePart2(grid)
	fmt.Println("Result:", result)
}

func log(v int, msg string) {
	if v <= VERBOSE {
		fmt.Println(msg)
	}
}

////////////////////////////////////////////////////////////

type Grid struct {
	rowcol.Grid[byte]
}

func ParseInput(name string) *Grid {
	lines := lib.ReadLines(name)
	return &Grid{rowcol.NewByteGridFromStrings(lines)}
}

////////////////////////////////////////////////////////////

func SolvePart1(grid *Grid) int {
	pos, ok := grid.Find(func(v byte) bool { return v == '^' })
	assert.True(ok)
	visited := rowcol.NewGrid[bool](grid.Size())

	grid.SetPos(pos, '.')
	visited.SetPos(pos, true)
	dir := rowcol.Up
	for {
		next := pos.Add(rowcol.Pos(dir))
		if !grid.IsValidPos(next) {
			return rowcol.Reduce(&visited, 0, func(acc int, v bool) int {
				if v {
					acc++
				}
				return acc
			})
		}

		if grid.GetPos(next) != '.' {
			dir = TurnRight(dir)
		} else {
			pos = next
			visited.SetPos(pos, true)
		}
	}
}

func TurnRight(dir rowcol.Direction) rowcol.Direction {
	switch dir {
	case rowcol.Up:
		return rowcol.Right
	case rowcol.Right:
		return rowcol.Down
	case rowcol.Down:
		return rowcol.Left
	case rowcol.Left:
		return rowcol.Up
	}
	panic("invalid direction")
}

////////////////////////////////////////////////////////////

func SolvePart2(grid *Grid) int {
	return 0
}
