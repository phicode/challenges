package main

// https://adventofcode.com/2023/day/16

import (
	"fmt"
	"time"

	"git.bind.ch/phil/challenges/lib"
)

var VERBOSE = 1

func main() {
	ProcessPart1("aoc23/day16/example.txt")
	ProcessPart1("aoc23/day16/input.txt")

	ProcessPart2("aoc23/day16/example.txt")
	ProcessPart2("aoc23/day16/input.txt")
}

func ProcessPart1(name string) {
	fmt.Println("Part 1 input:", name)
	lines := lib.ReadLines(name)
	grid := ParseGrid(lines)
	t := time.Now()
	grid.FollowAndMark(0, 0, LeftToRight)
	fmt.Println("Marking done", time.Now().Sub(t))
	fmt.Println("Analysing energized fields")
	grid.PrintEnergized()
	e := grid.CountEnergized()
	fmt.Println("Energized:", e)

	fmt.Println()
}

func ProcessPart2(name string) {
	fmt.Println("Part 2 input:", name)
	lines := lib.ReadLines(name)
	grid := ParseGrid(lines)
	e := Part2BruteForce(grid)
	fmt.Println("Energized:", e)

	fmt.Println()
}

func Part2BruteForce(grid *Grid) int {
	rows, cols := grid.Rows(), grid.Columns()
	var e int
	// top row down
	for col := 0; col < cols; col++ {
		grid.FollowAndMark(0, col, TopToBottom)
		e = max(e, grid.CountEnergized())
		grid.ResetEnergized()
	}
	// bottom row up
	for col := 0; col < cols; col++ {
		grid.FollowAndMark(rows-1, col, BottomToTop)
		e = max(e, grid.CountEnergized())
		grid.ResetEnergized()
	}
	// left column to the right
	for row := 0; row < rows; row++ {
		grid.FollowAndMark(row, 0, LeftToRight)
		e = max(e, grid.CountEnergized())
		grid.ResetEnergized()
	}
	// right column to the left
	for row := 0; row < rows; row++ {
		grid.FollowAndMark(row, cols-1, RightToLeft)
		e = max(e, grid.CountEnergized())
		grid.ResetEnergized()
	}
	return e
}

func log(v int, msg string) {
	if v <= VERBOSE {
		fmt.Println(msg)
	}
}

////////////////////////////////////////////////////////////

type Grid struct {
	lib.Grid[byte]
	LaserDirections lib.Grid[Directions]
}

const (
	TopToBottom Directions = 1 << 0
	BottomToTop Directions = 1 << 1
	LeftToRight Directions = 1 << 2
	RightToLeft Directions = 1 << 3
)

const (
	Empty              byte = '.'
	MirrorSlash        byte = '/'
	MirrorBackSlash    byte = '\\'
	SplitterVertical   byte = '|'
	SplitterHorizontal byte = '-'
)

type Directions int

func (d Directions) Set(dir Directions) Directions { return d | dir }
func (d Directions) IsSet(dir Directions) bool     { return d&dir == dir }
func (d Directions) IsHorizontal() bool            { return d == LeftToRight || d == RightToLeft }
func (d Directions) IsVertical() bool              { return d == TopToBottom || d == BottomToTop }
func (d Directions) String() string {
	switch d {
	case TopToBottom:
		return "TopToBottom"
	case BottomToTop:
		return "BottomToTop"
	case LeftToRight:
		return "LeftToRight"
	case RightToLeft:
		return "RightToLeft"
	default:
		return "Unknown"
	}
}

func ParseGrid(lines []string) *Grid {
	grid := lib.NewByteGridFromStrings(lines)
	return &Grid{
		Grid:            lib.NewByteGridFromStrings(lines),
		LaserDirections: lib.NewGrid[Directions](grid.Rows(), grid.Columns()),
	}
}

func (g *Grid) FollowAndMark(row, col int, dir Directions) {
	if !g.IsValidPosition(row, col) {
		// laser leaving the grid is ignored
		return
	}
	v := g.LaserDirections.Get(row, col)
	if v.IsSet(dir) {
		if VERBOSE >= 2 {
			fmt.Printf("(%d,%d) already visited %s\n", row, col, v)
		}
		return // already marked this direction
	}
	v = v.Set(dir)
	g.LaserDirections.Set(row, col, v)
	if VERBOSE >= 2 {
		fmt.Printf("(%d,%d) following %s\n", row, col, dir)
	}
	field := g.Get(row, col)
	switch field {
	case Empty:
		row, col = NextPosition(row, col, dir)
		g.FollowAndMark(row, col, dir)
	case SplitterVertical: // |
		if dir.IsVertical() {
			// beam continues on its way
			row, col = NextPosition(row, col, dir)
			g.FollowAndMark(row, col, dir)
		} else {
			// beam is split
			rowA, colA := NextPosition(row, col, BottomToTop)
			g.FollowAndMark(rowA, colA, BottomToTop)
			rowB, colB := NextPosition(row, col, TopToBottom)
			g.FollowAndMark(rowB, colB, TopToBottom)
		}
	case SplitterHorizontal: // -
		if dir.IsHorizontal() {
			// beam continues on its way
			row, col = NextPosition(row, col, dir)
			g.FollowAndMark(row, col, dir)
		} else {
			// beam is split
			rowA, colA := NextPosition(row, col, LeftToRight)
			g.FollowAndMark(rowA, colA, LeftToRight)
			rowB, colB := NextPosition(row, col, RightToLeft)
			g.FollowAndMark(rowB, colB, RightToLeft)
		}

	case MirrorSlash: // /
		dir = TranslateSlash(dir)
		row, col = NextPosition(row, col, dir)
		g.FollowAndMark(row, col, dir)

	case MirrorBackSlash: // \
		dir = TranslateBashSlash(dir)
		row, col = NextPosition(row, col, dir)
		g.FollowAndMark(row, col, dir)

	default:
		panic(fmt.Errorf("invalid field: %c", field))
	}
}

func TranslateSlash(dir Directions) Directions {
	// /
	switch dir {
	case BottomToTop:
		return LeftToRight
	case TopToBottom:
		return RightToLeft
	case LeftToRight:
		return BottomToTop
	case RightToLeft:
		return TopToBottom
	default:
		panic("invalid direction")
	}
}
func TranslateBashSlash(dir Directions) Directions {
	// \
	switch dir {
	case BottomToTop:
		return RightToLeft
	case TopToBottom:
		return LeftToRight
	case LeftToRight:
		return TopToBottom
	case RightToLeft:
		return BottomToTop
	default:
		panic("invalid direction")
	}
}

func NextPosition(row, col int, dir Directions) (int, int) {
	switch dir {
	case BottomToTop:
		return row - 1, col
	case TopToBottom:
		return row + 1, col
	case LeftToRight:
		return row, col + 1
	case RightToLeft:
		return row, col - 1
	default:
		panic("invalid input")
	}
}

func (g *Grid) CountEnergized() int {
	rows, cols := g.Rows(), g.Columns()
	energized := 0
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if g.LaserDirections.Get(r, c) != 0 {
				energized++
			}
		}
	}
	return energized
}

func (g *Grid) ResetEnergized() {
	rows, cols := g.Rows(), g.Columns()
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			g.LaserDirections.Set(r, c, 0)
		}
	}
}

func (g *Grid) PrintEnergized() {
	rows, cols := g.Rows(), g.Columns()
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if g.LaserDirections.Get(r, c) != 0 {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}
