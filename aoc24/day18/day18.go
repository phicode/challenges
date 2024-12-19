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

	lib.Timed("Part 2", ProcessPart2Example, "aoc24/day18/example.txt")
	lib.Timed("Part 2", ProcessPart2Input, "aoc24/day18/input.txt")

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

func ProcessPart2Example(name string) {
	fmt.Println("Part 2 input:", name)
	input := ReadAndParseInput(name)
	result := SolvePart2(input, 7, 12)
	fmt.Println("Result:", result)
}

func ProcessPart2Input(name string) {
	fmt.Println("Part 2 input:", name)
	input := ReadAndParseInput(name)
	result := SolvePart2(input, 71, 1024)
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
	l, ok := ShortestPath(grid)
	assert.True(ok)
	return l
}

func ShortestPath(grid rowcol.Grid[byte]) (int, bool) {
	nodes := findReachable(grid, rowcol.Pos{Row: 0, Col: 0})

	var start = func(a rowcol.Pos) bool { return a.Row == 0 && a.Col == 0 }
	var neigh = func(p rowcol.Pos) []rowcol.Pos {
		return freeNeighbors(grid, p)
	}
	var paths map[rowcol.Pos]*graphs.Node[rowcol.Pos]
	paths = graphs.Dijkstra(nodes, start, neigh)
	rows, cols := grid.Size()
	end := rowcol.Pos{Row: rows - 1, Col: cols - 1}
	endNode := paths[end]
	if endNode == nil {
		return 0, false
	}
	return len(endNode.GetPath()) - 1, true
}

func findReachable(grid rowcol.Grid[byte], start rowcol.Pos) []rowcol.Pos {
	reachable := make(map[rowcol.Pos]bool)

	_findReachable(grid, reachable, start)
	return lib.MapKeys(reachable)
}

func _findReachable(grid rowcol.Grid[byte], reachable map[rowcol.Pos]bool, p rowcol.Pos) {
	reachable[p] = true
	for _, dir := range rowcol.Directions {
		next := p.AddDir(dir)
		if reachable[next] {
			continue
		}
		if grid.IsValidPos(next) && grid.GetPos(next) == '.' {
			_findReachable(grid, reachable, next)
		}
	}
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

func SolvePart2(input Input, gridSize, numBytes int) string {
	assert.True(numBytes <= len(input.Coordinates))
	grid := rowcol.NewGrid[byte](gridSize, gridSize)
	grid.Reset('.')
	for i := 0; i < numBytes; i++ {
		grid.SetPos(input.Coordinates[i], '#')
	}
	_, ok := ShortestPath(grid)
	assert.True(ok)
	numCoords := len(input.Coordinates)
	for i := numBytes; i < numCoords; i++ {
		c := input.Coordinates[i]
		grid.SetPos(c, '#')
		if _, ok := ShortestPath(grid); !ok {
			return fmt.Sprintf("%d,%d", c.Col, c.Row)
		}

	}
	panic("no solution found")
}
