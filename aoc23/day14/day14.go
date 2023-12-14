package main

// https://adventofcode.com/2023/day/XX

import (
	"fmt"

	"git.bind.ch/phil/challenges/lib"
)

var VERBOSE = 1

func main() {
	ProcessPart1("aoc23/day14/example.txt")
	ProcessPart1("aoc23/day14/input.txt")

	//ProcessPart2("aoc23/day14/example.txt")
	//ProcessPart2("aoc23/day14/input.txt")
}

func ProcessPart1(name string) {
	fmt.Println("Part 1 input:", name)
	lines := lib.ReadLines(name)
	g := ParseInput(lines)
	g.MoveUp()
	if VERBOSE >= 1 {
		g.Print()
	}
	load := g.CalcLoad()
	fmt.Println("Load:", load)

	fmt.Println()
}

func ProcessPart2(name string) {
	fmt.Println("Part 2 input:", name)
	lines := lib.ReadLines(name)
	_ = lines

	fmt.Println()
}

func log(v int, msg string) {
	if v <= VERBOSE {
		fmt.Println(msg)
	}
}

////////////////////////////////////////////////////////////

const (
	Empty     = '.'
	RoundRock = 'O'
	CubeRock  = '#'
)

type Grid struct {
	lib.Grid[byte]
}

func ParseInput(lines []string) *Grid {
	return &Grid{lib.NewByteGridFromStrings(lines)}
}

func (g *Grid) MoveUp() {
	rows, cols := g.Rows(), g.Columns()
	for row := 1; row < rows; row++ {
		for col := 0; col < cols; col++ {
			if g.Data[row][col] == RoundRock {
				g.MoveStoneUp(row, col)
			}
		}
	}
}

func (g *Grid) MoveStoneUp(row, col int) {
	for row > 0 {
		if g.Data[row-1][col] != Empty {
			return
		}
		g.Data[row-1][col] = RoundRock
		g.Data[row][col] = Empty
		row--
	}
}

func (g *Grid) Print() {
	for row, value := range g.Data {
		fmt.Printf("%-2d %s\n", g.LoadFactor(row), string(value))
	}
}

func (g *Grid) LoadFactor(row int) int {
	rows := g.Rows()
	return rows - row
}

func (g *Grid) CalcLoad() int {
	rows := g.Rows()
	var sum int
	for row := 0; row < rows; row++ {
		c := Count(g.Data[row], RoundRock)
		sum += c * g.LoadFactor(row)
	}
	return sum
}

func Count(xs []byte, value byte) int {
	c := 0
	for _, x := range xs {
		if x == value {
			c++
		}
	}
	return c
}
