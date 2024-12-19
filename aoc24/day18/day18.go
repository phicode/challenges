package main

// https://adventofcode.com/2024/day/18

import (
	"flag"
	"fmt"

	"github.com/phicode/challenges/lib"
	"github.com/phicode/challenges/lib/assert"
	"github.com/phicode/challenges/lib/graphs"
	"github.com/phicode/challenges/lib/rowcol"
)

func main() {
	flag.Parse()
	lib.Timed("Part 1", ProcessPart1Example, "aoc24/day18/example.txt")
	lib.Timed("Part 1", ProcessPart1Real, "aoc24/day18/input.txt")

	lib.Timed("Part 2", ProcessPart2, "aoc24/day18/example.txt")
	lib.Timed("Part 2", ProcessPart2, "aoc24/day18/input.txt")

	//lib.Profile(1, "day18.pprof", "Part 2", ProcessPart2, "aoc24/day18/input.txt")
}

func ProcessPart1Example(name string) {
	fmt.Println("Part 1 input:", name)
	input := ReadAndParseInput(name)
	result := SolvePart1(input, 7, 12)
	fmt.Println("Result:", result)
}

func ProcessPart1Real(name string) {
	fmt.Println("Part 1 input:", name)
	input := ReadAndParseInput(name)
	result := SolvePart1(input, 71, 1024)
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
	Coordinates []rowcol.Pos
}

func ReadAndParseInput(name string) Input {
	lines := lib.ReadLines(name)
	return ParseInput(lines)
}

func ParseInput(lines []string) Input {
	var cs []rowcol.Pos
	for _, line := range lines {
		var c rowcol.Pos
		n, err := fmt.Sscanf(line, "%d,%d", &c.Col, &c.Row)
		assert.True(n == 2 && err == nil)
		cs = append(cs, c)
	}
	return Input{
		Coordinates: cs,
	}
}

////////////////////////////////////////////////////////////

func SolvePart1(input Input, gridSize, numBytes int) int {
	assert.True(numBytes <= len(input.Coordinates))
	grid := rowcol.NewGrid[byte](gridSize, gridSize)
	grid.Reset('.')
	for i := 0; i < numBytes; i++ {
		grid.SetPos(input.Coordinates[i], '#')
	}
	return ShortestPath(grid)
}

func ShortestPath(grid rowcol.Grid[byte]) int {
	var nodes []rowcol.Pos
	for p, v := range grid.Iterator() {
		if v == '.' {
			nodes = append(nodes, p)
		}
	}

	var start = func(a rowcol.Pos) bool { return a.Row == 0 && a.Col == 0 }
	var neigh = func(p rowcol.Pos) []rowcol.Pos {
		return freeNeighbors(grid, p)
	}
	var paths map[rowcol.Pos]*graphs.Node[rowcol.Pos]
	paths = graphs.Dijkstra(nodes, start, neigh)
	rows, cols := grid.Size()
	end := rowcol.Pos{rows - 1, cols - 1}
	endNode := paths[end]
	assert.True(endNode != nil)
	return len(endNode.GetPath()) - 1
}

func freeNeighbors(grid rowcol.Grid[byte], p rowcol.Pos) []rowcol.Pos {
	var rv []rowcol.Pos
	for _, dir := range rowcol.Directions {
		next := p.AddDir(dir)
		if grid.IsValidPos(next) && grid.GetPos(next) == '.' {
			rv = append(rv, next)
		}
	}
	return rv
}

////////////////////////////////////////////////////////////

func SolvePart2(input Input) int {
	return 0
}
