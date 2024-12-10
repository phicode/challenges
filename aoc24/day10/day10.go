package main

// https://adventofcode.com/2024/day/10

import (
	"flag"
	"fmt"

	"github.com/phicode/challenges/lib"
	"github.com/phicode/challenges/lib/rowcol"
)

func main() {
	flag.Parse()
	lib.Timed("Part 1", ProcessPart1, "aoc24/day10/example.txt")
	lib.Timed("Part 1", ProcessPart1, "aoc24/day10/input.txt")

	//	lib.Timed("Part 2", ProcessPart2, "aoc24/day10/example.txt")
	//	lib.Timed("Part 2", ProcessPart2, "aoc24/day10/input.txt")

	//lib.Profile(1, "day10.pprof", "Part 2", ProcessPart2, "aoc24/dayXX/input.txt")
}

func ProcessPart1(name string) {
	fmt.Println("Part 1 input:", name)
	input := ParseInput(name)
	result := SolvePart1(input)
	fmt.Println("Result:", result)
}

func ProcessPart2(name string) {
	fmt.Println("Part 2 input:", name)
	input := ParseInput(name)
	result := SolvePart2(input)
	fmt.Println("Result:", result)
}

////////////////////////////////////////////////////////////

type Input struct {
	grid rowcol.Grid[byte]
}

func ParseInput(name string) *Input {
	lines := lib.ReadLines(name)
	lines = lib.RemoveEmptyLines(lines)
	grid := rowcol.NewByteGridFromStrings(lines)
	return &Input{grid}
}

////////////////////////////////////////////////////////////

func SolvePart1(input *Input) int {
	visited := rowcol.NewGrid[bool](input.grid.Size())
	trailheads := input.grid.FindAll(func(x byte) bool { return x == '0' })
	sum := 0
	for _, th := range trailheads {
		sum += markDirs(input.grid, visited, th, byte('0'))
		visited.Reset(false)
	}
	return sum
}

func markDirs(g rowcol.Grid[byte], visited rowcol.Grid[bool], p rowcol.Pos, value byte) int {
	if g.GetPos(p) != value {
		return 0
	}
	if visited.GetPos(p) {
		return 0
	}
	sum := mark(g, visited, p, rowcol.Up, value+1) +
		mark(g, visited, p, rowcol.Down, value+1) +
		mark(g, visited, p, rowcol.Left, value+1) +
		mark(g, visited, p, rowcol.Right, value+1)
	visited.SetPos(p, true)
	return sum
}

func mark(g rowcol.Grid[byte], visited rowcol.Grid[bool], p rowcol.Pos, dir rowcol.Direction, follow byte) int {
	p = p.Add(rowcol.Pos(dir))
	if !visited.IsValidPos(p) {
		return 0
	}
	if visited.GetPos(p) {
		return 0
	}
	if g.GetPos(p) != follow {
		return 0
	}
	if follow == '9' {
		visited.SetPos(p, true)
		return 1
	}
	return markDirs(g, visited, p, follow)
}

////////////////////////////////////////////////////////////

func SolvePart2(input *Input) int {
	ways := rowcol.NewGrid[int](input.grid.Size())
	trailheads := input.grid.FindAll(func(x byte) bool { return x == '0' })
	sum := 0
	for _, th := range trailheads {
		sum += markDirs2(input.grid, ways, th, byte('0'))
		visited.Reset(false)
	}
	return sum
}

func markDirs2(g rowcol.Grid[byte], ways rowcol.Grid[int], p rowcol.Pos, value byte) int {
	if g.GetPos(p) != value {
		return 0
	}
	if cached := ways.GetPos(p); cached > 0 {
		return cached
	}
	sum := mark2(g, visited, p, rowcol.Up, value+1) +
		mark2(g, visited, p, rowcol.Down, value+1) +
		mark2(g, visited, p, rowcol.Left, value+1) +
		mark2(g, visited, p, rowcol.Right, value+1)
	visited.SetPos(p, true)
	return sum
}

func mark2(g rowcol.Grid[byte], ways rowcol.Grid[int], p rowcol.Pos, dir rowcol.Direction, follow byte) int {
	p = p.Add(rowcol.Pos(dir))
	if !visited.IsValidPos(p) {
		return 0
	}
	if visited.GetPos(p) {
		return 0
	}
	if g.GetPos(p) != follow {
		return 0
	}
	if follow == '9' {
		visited.SetPos(p, true)
		return 1
	}
	return markDirs2(g, visited, p, follow)
}
