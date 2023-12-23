package main

// https://adventofcode.com/2023/day/23

import (
	"bytes"
	"fmt"

	"git.bind.ch/phil/challenges/lib"
	"git.bind.ch/phil/challenges/lib/rowcol"
)

var VERBOSE = 1

func main() {
	ProcessPart1("aoc23/day23/example.txt") // 84
	ProcessPart1("aoc23/day23/input.txt")   // 2278

	ProcessPart2("aoc23/day23/example.txt")
	ProcessPart2("aoc23/day23/input.txt")
}

func ProcessPart1(name string) {
	fmt.Println("Part 1 input:", name)
	lines := lib.ReadLines(name)
	g := ParseGraph(lines)
	SolvePart1(g)

	fmt.Println()
}

func ProcessPart2(name string) {
	fmt.Println("Part 2 input:", name)
	lines := lib.ReadLines(name)
	_ = lines

	fmt.Println()
}

////////////////////////////////////////////////////////////

type Grid struct {
	rowcol.Grid[byte]

	Start, End rowcol.Pos
}

func ParseGraph(lines []string) *Grid {
	g := &Grid{Grid: rowcol.NewByteGridFromStrings(lines)}
	g.Start = g.FindStart()
	g.End = g.FindEnd()
	return g
}

func (g *Grid) FindStart() rowcol.Pos {
	col := bytes.IndexByte(g.Grid.Data[0], '.')
	if col == -1 {
		panic("start not found")
	}
	return rowcol.Pos{Row: 0, Col: col}
}

func (g *Grid) FindEnd() rowcol.Pos {
	row := g.Rows() - 1
	col := bytes.IndexByte(g.Grid.Data[row], '.')
	if col == -1 {
		panic("end not found")
	}
	return rowcol.Pos{Row: row, Col: col}
}

func SolvePart1(g *Grid) {
	visited := make(map[rowcol.Pos]bool)
	visited[g.Start] = true
	d := AllSteps(g, g.Start, visited)
	fmt.Println("Distance:", d)
}

func AllSteps(g *Grid, current rowcol.Pos, visited map[rowcol.Pos]bool) int {
	maxd := 0
	for _, dir := range rowcol.Directions {
		next := current.AddDir(dir)
		d := OneStep(g, next, visited)
		maxd = max(maxd, d)
	}
	return maxd
}

func OneStep(g *Grid, pos rowcol.Pos, visited map[rowcol.Pos]bool) int {
	if !g.IsValidPosition(pos.Row, pos.Col) {
		return 0
	}
	if visited[pos] {
		return 0
	}
	if pos == g.End {
		return 1
	}
	v := g.Get(pos.Row, pos.Col)
	if v == '#' {
		return 0
	}
	visited[pos] = true
	defer func() { delete(visited, pos) }()
	if v == '.' {
		return AllSteps(g, pos, visited) + 1
	}
	dir, ok := Directions[v]
	if !ok {
		panic(fmt.Errorf("invalid field: %c", v))
	}
	return OneStep(g, pos.AddDir(dir), visited) + 1
}

var Directions = map[byte]rowcol.Direction{
	'>': rowcol.Right,
	'<': rowcol.Left,
	'^': rowcol.Up,
	'v': rowcol.Down,
}
