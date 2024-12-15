package main

// https://adventofcode.com/2024/day/6

import (
	"fmt"

	"github.com/phicode/challenges/lib"
	"github.com/phicode/challenges/lib/assert"
	"github.com/phicode/challenges/lib/rowcol"
)

func main() {
	lib.Timed("Part 1", ProcessPart1, "aoc24/day06/example.txt")
	lib.Timed("Part 1", ProcessPart1, "aoc24/day06/input.txt")

	lib.Timed("Part 2", ProcessPart2, "aoc24/day06/example.txt")
	lib.Timed("Part 2", ProcessPart2, "aoc24/day06/input.txt")
	//lib.Profile(20, "part2.pprof", "Part 2", ProcessPart2, "aoc24/day06/input.txt")
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
	pos, ok := grid.FindFirst(func(v byte) bool { return v == '^' })
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
			dir = dir.Right()
		} else {
			pos = next
			visited.SetPos(pos, true)
		}
	}
}

////////////////////////////////////////////////////////////

func SolvePart2(grid *Grid) int {
	pos, ok := grid.FindFirst(func(v byte) bool { return v == '^' })
	assert.True(ok)
	grid.SetPos(pos, '.')

	nloops := 0
	rows, cols := grid.Size()

	visited := rowcol.NewGrid[DirSet](grid.Size())

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if r == pos.Row && c == pos.Col {
				continue // guard start position
			}
			if grid.Get(r, c) == '.' {
				grid.Set(r, c, '#')
				if IsLoop(grid, pos, visited) {
					nloops++
				}
				visited.Reset(DirSet(0))
				grid.Set(r, c, '.')
			}
		}
	}
	return nloops
}

func IsLoop(grid *Grid, pos rowcol.Pos, visited rowcol.Grid[DirSet]) bool {
	dir := rowcol.Up
	visited.SetPos(pos, NewDirSet(dir))
	for {
		next := pos.Add(rowcol.Pos(dir))
		if !grid.IsValidPos(next) {
			return false // guard is leaving the map -> no loop
		}

		if grid.GetPos(next) != '.' {
			dir = dir.Right()
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
	return set | dirMask(dir)
}

func (set DirSet) IsSet(dir rowcol.Direction) bool {
	return (set & dirMask(dir)) != 0
}

// Left  {Row: 0, Col: -1} = 1*3 + 0 = 3
// Right {Row: 0, Col: +1} = 1*3 + 2 = 5
// Up    {Row: -1, Col: 0} = 0*3 + 1 = 1
// Down  {Row: +1, Col: 0} = 2*3 + 1 = 7
func dirMask(dir rowcol.Direction) DirSet {
	shift := (dir.Row+1)*3 + (dir.Col + 1)
	return DirSet(1) << shift
}
