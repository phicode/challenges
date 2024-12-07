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

func (g *Grid) Print() {
	for _, row := range g.Data {
		fmt.Println(string(row))
	}
	fmt.Println()
}

func ParseInput(name string) *Grid {
	lines := lib.ReadLines(name)
	return &Grid{rowcol.NewByteGridFromStrings(lines)}
}

////////////////////////////////////////////////////////////

func SolvePart1(grid *Grid) int {
	pos, ok := grid.Find(func(v byte) bool { return v == '^' })
	assert.True(ok)
	grid.SetPos(pos, '.')

	visited := rowcol.NewGrid[bool](grid.Size())
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
	pos, ok := grid.Find(func(v byte) bool { return v == '^' })
	assert.True(ok)
	grid.SetPos(pos, '.')

	nloops := 0
	rows, cols := grid.Size()
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if r == pos.Row && c == pos.Col {
				continue // guard start position
			}
			if grid.Get(r, c) == '.' {
				grid.Set(r, c, '#')
				//grid.Print()
				if IsLoop(grid, pos) {
					nloops++
				}
				grid.Set(r, c, '.')
			}
		}
	}
	return nloops
}

func IsLoop(grid *Grid, pos rowcol.Pos) bool {
	visited := rowcol.NewGrid[DirSet](grid.Size())

	dir := rowcol.Up
	visited.SetPos(pos, NewDirSet(dir))
	for {
		next := pos.Add(rowcol.Pos(dir))
		if !grid.IsValidPos(next) {
			return false // guard is leaving the map -> no loop
		}

		if grid.GetPos(next) != '.' {
			dir = TurnRight(dir)
			continue
		}
		pos = next
		set := visited.GetPos(pos)
		if set.IsSet(dir) {
			return true
		}
		set = set.Add(dir)
		visited.SetPos(pos, set)
	}
}

type DirSet byte

func NewDirSet(dir rowcol.Direction) DirSet {
	var set DirSet
	set = set.Add(dir)
	return set
}

func (set DirSet) Add(dir rowcol.Direction) DirSet {
	for i, d := range rowcol.Directions {
		if d == dir {
			mask := DirSet(1) << uint(i)
			return set | mask
		}
	}
	panic("invalid direction")
}

func (set DirSet) IsSet(dir rowcol.Direction) bool {
	for i, d := range rowcol.Directions {
		if d == dir {
			mask := DirSet(1) << uint(i)
			return (set & mask) != 0
		}
	}
	panic("invalid direction")
}
