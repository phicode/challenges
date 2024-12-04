package main

// https://adventofcode.com/2024/day/4

import (
	"fmt"

	"git.bind.ch/phil/challenges/lib"
	"git.bind.ch/phil/challenges/lib/rowcol"
)

// TODO: timing boilerplate
var VERBOSE = 1

func main() {
	ProcessPart1("aoc24/day04/example.txt")
	ProcessPart1("aoc24/day04/input.txt")

	ProcessPart2("aoc24/day04/example.txt")
	ProcessPart2("aoc24/day04/input.txt")
}

func ProcessPart1(name string) {
	fmt.Println("Part 1 input:", name)
	lines := lib.ReadLines(name)
	grid := ParseInput(lines)
	count := grid.CountXmas()
	fmt.Println("Count:", count)
}

func ProcessPart2(name string) {
	fmt.Println("Part 2 input:", name)
	lines := lib.ReadLines(name)
	grid := ParseInput(lines)
	count := grid.CountXmasP2()
	fmt.Println("Count:", count)
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

func ParseInput(lines []string) *Grid {
	return &Grid{rowcol.NewByteGridFromStrings(lines)}
}

func (g *Grid) CountXmas() int {
	rows, cols := g.Size()
	total := 0
	for c := 0; c < cols; c++ {
		for r := 0; r < rows; r++ {
			for _, dir := range Directions {
				if g.IsXmas(rowcol.Pos{Row: r, Col: c}, dir) {
					total++
				}
			}
		}
	}
	return total
}

func (g *Grid) CountXmasP2() int {
	rows, cols := g.Size()
	total := 0
	for c := 1; c < cols-1; c++ {
		for r := 1; r < rows-1; r++ {
			if g.IsMasCenter(rowcol.Pos{Row: r, Col: c}) {
				total++
			}
		}
	}
	return total
}

var xmas = [4]byte{'X', 'M', 'A', 'S'}
var Directions = []rowcol.Direction{
	rowcol.Left, rowcol.Right, rowcol.Up, rowcol.Down,
	rowcol.UpLeft, rowcol.UpRight,
	rowcol.DownLeft, rowcol.DownRight,
}

func (g *Grid) IsXmas(p rowcol.Pos, d rowcol.Direction) bool {
	for i, val := range xmas {
		if !g.IsValue(val, p.Add(d.MulS(i))) {
			return false
		}
	}
	return true
}

func (g *Grid) IsValue(value byte, p rowcol.Pos) bool {
	return g.IsValidPosition(p.Row, p.Col) &&
		g.Get(p.Row, p.Col) == value
}

func (g *Grid) IsXmasP2(p rowcol.Pos, d rowcol.Direction) bool {
	for i, val := range xmas {
		if !g.IsValue(val, p.Add(d.MulS(i))) {
			return false
		}
	}
	return true
}

func (g *Grid) IsMasCenter(p rowcol.Pos) bool {
	// Coordinates
	// a.b
	// .A.
	// c.d
	a := p.Add(rowcol.Pos(rowcol.UpLeft))
	b := p.Add(rowcol.Pos(rowcol.UpRight))
	c := p.Add(rowcol.Pos(rowcol.DownLeft))
	d := p.Add(rowcol.Pos(rowcol.DownRight))
	if !g.IsValidPosition(a.Row, a.Col) || !g.IsValidPosition(d.Row, d.Col) {
		return false
	}
	if !g.IsValue('A', p) {
		return false
	}
	return isMS(g.Get(a.Row, a.Col), g.Get(d.Row, d.Col)) &&
		isMS(g.Get(b.Row, b.Col), g.Get(c.Row, c.Col))
}

func isMS(a, b byte) bool {
	if a == 'M' {
		return b == 'S'
	}
	return a == 'S' && b == 'M'
}
