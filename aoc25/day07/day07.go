package main

// https://adventofcode.com/2025/day/7

import (
	"flag"
	"fmt"

	"github.com/phicode/challenges/lib"
	"github.com/phicode/challenges/lib/rowcol"
)

func main() {
	flag.Parse()
	lib.Timed("Part 1", ProcessPart1, "aoc25/day07/example.txt")
	lib.Timed("Part 1", ProcessPart1, "aoc25/day07/input.txt")

	lib.Timed("Part 2", ProcessPart2, "aoc25/day07/example.txt")
	lib.Timed("Part 2", ProcessPart2, "aoc25/day07/input.txt")

	//lib.Profile(1, "day07.pprof", "Part 2", ProcessPart2, "aoc25/day07/input.txt")
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

const (
	START    = 'S'
	SPLITTER = '^'
	FREE     = '.'
	BEAM     = '|'
)

type Input struct {
	grid rowcol.Grid[byte]
}

func ReadAndParseInput(name string) Input {
	lines := lib.ReadLines(name)
	return ParseInput(lines)
}

func ParseInput(lines []string) Input {
	return Input{
		grid: rowcol.NewByteGridFromStrings(lines),
	}
}

////////////////////////////////////////////////////////////

func SolvePart1(input Input) int {
	grid := input.grid
	current := grid.MustFindFirst(func(x byte) bool { return x == START })
	splits := follow(grid, current)
	return splits
}

func follow(g rowcol.Grid[byte], current rowcol.Pos) int {
	if !g.IsValidPos(current) {
		return 0
	}
	splits := 0
	for {
		next := current.AddDir(rowcol.Down)
		if !g.IsValidPos(next) {
			return splits
		}
		cell := g.GetPos(next)
		switch cell {
		case BEAM:
			return splits // already accounted for
		case FREE:
			current = next
			g.SetPos(next, BEAM)
			continue
		case SPLITTER:
			splits++
			left := next.AddDir(rowcol.Left)
			right := next.AddDir(rowcol.Right)
			return splits + follow(g, left) + follow(g, right)
		default:
			panic(fmt.Errorf("invalid cell value: %v", cell))
		}
	}
}

////////////////////////////////////////////////////////////

func SolvePart2(input Input) int {
	grid := input.grid
	current := grid.MustFindFirst(func(x byte) bool { return x == START })
	mem := rowcol.NewGrid[int](grid.Size())
	splits := follow2(grid, mem, current.AddDir(rowcol.Down))
	return splits
}

func follow2(g rowcol.Grid[byte], mem rowcol.Grid[int], current rowcol.Pos) int {
	if !g.IsValidPos(current) {
		return 1
	}
	if m := mem.GetPos(current); m > 0 {
		return m
	}
	cell := g.GetPos(current)
	switch cell {
	case FREE:
		m := follow2(g, mem, current.AddDir(rowcol.Down))
		mem.SetPos(current, m)
		return m
	case SPLITTER:
		left := current.AddDir(rowcol.Left)
		right := current.AddDir(rowcol.Right)
		m := follow2(g, mem, left) + follow2(g, mem, right)
		mem.SetPos(current, m)
		return m
	default:
		panic(fmt.Errorf("invalid cell value: %v", cell))
	}
}
