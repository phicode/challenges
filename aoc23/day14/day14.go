package main

// https://adventofcode.com/2023/day/XX

import (
	"fmt"
	"time"

	"git.bind.ch/phil/challenges/lib"
)

var VERBOSE = 1

func main() {
	ProcessPart1("aoc23/day14/example.txt")
	ProcessPart1("aoc23/day14/input.txt")

	ProcessPart2("aoc23/day14/example.txt", 3)
	VERBOSE = 0
	ProcessPart2("aoc23/day14/input.txt", 1_000_000_000)
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

func ProcessPart2(name string, cycles int) {
	fmt.Println("Part 2 input:", name)
	lines := lib.ReadLines(name)
	g := ParseInput(lines)
	t := time.Now()
	var loads []int
	for c := 0; c < cycles; c++ {
		g.SpinCycle()
		loads = append(loads, g.CalcLoad())
		if len(loads) > 100_0000 && FindLoop(loads) {
			fmt.Println("loop found")
		}
		if VERBOSE >= 1 {
			fmt.Println("After", c+1, "cycles:")
			g.Print()
			fmt.Println()
		}
		if c%100_000 == 0 {
			percent := 100.0 * float64(c) / float64(cycles)
			fmt.Printf("\r%2f %%, %v", percent, time.Now().Sub(t))
		}
	}
	fmt.Println("")
	load := g.CalcLoad()
	fmt.Println("Load:", load)

	fmt.Println()
}

func FindLoop(loads []int) bool {
	maxlen:=

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

func (g *Grid) MoveDown() {
	rows, cols := g.Rows(), g.Columns()
	for row := rows - 2; row >= 0; row-- {
		for col := 0; col < cols; col++ {
			if g.Data[row][col] == RoundRock {
				g.MoveStoneDown(row, col)
			}
		}
	}
}
func (g *Grid) MoveLeft() {
	rows, cols := g.Rows(), g.Columns()
	for row := 0; row < rows; row++ {
		for col := 1; col < cols; col++ {
			if g.Data[row][col] == RoundRock {
				g.MoveStoneLeft(row, col)
			}
		}
	}
}
func (g *Grid) MoveRight() {
	rows, cols := g.Rows(), g.Columns()
	for row := 0; row < rows; row++ {
		for col := cols - 2; col >= 0; col-- {
			if g.Data[row][col] == RoundRock {
				g.MoveStoneRight(row, col)
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
func (g *Grid) MoveStoneDown(row, col int) {
	rows := g.Rows()
	for row < rows-1 {
		if g.Data[row+1][col] != Empty {
			return
		}
		g.Data[row+1][col] = RoundRock
		g.Data[row][col] = Empty
		row++
	}
}
func (g *Grid) MoveStoneLeft(row, col int) {
	for col > 0 {
		if g.Data[row][col-1] != Empty {
			return
		}
		g.Data[row][col-1] = RoundRock
		g.Data[row][col] = Empty
		col--
	}
}
func (g *Grid) MoveStoneRight(row, col int) {
	cols := g.Columns()
	for col < cols-1 {
		if g.Data[row][col+1] != Empty {
			return
		}
		g.Data[row][col+1] = RoundRock
		g.Data[row][col] = Empty
		col++
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

func (g *Grid) SpinCycle() {
	// north
	g.MoveUp()
	// west
	g.MoveLeft()
	// south
	g.MoveDown()
	// east
	g.MoveRight()
}


// start inclusive, end exclusive
func IsLoop(xs []int , start, end int) bool {
	l := end - start
	if end+l >= len(xs) {
		return false
	}
	for i := 0; i < l; i++ {
		if xs[start+i] != xs[end+i] {
			return false
		}
	}
	return true
}