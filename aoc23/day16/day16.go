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
	grid.FollowAndMark(0, 0, Right)
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
		grid.FollowAndMark(0, col, Down)
		e = max(e, grid.CountEnergized())
		grid.ResetEnergized()
	}
	// bottom row up
	for col := 0; col < cols; col++ {
		grid.FollowAndMark(rows-1, col, Up)
		e = max(e, grid.CountEnergized())
		grid.ResetEnergized()
	}
	// left column to the right
	for row := 0; row < rows; row++ {
		grid.FollowAndMark(row, 0, Right)
		e = max(e, grid.CountEnergized())
		grid.ResetEnergized()
	}
	// right column to the left
	for row := 0; row < rows; row++ {
		grid.FollowAndMark(row, cols-1, Left)
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
	LaserDirection lib.Grid[Direction]
}

const (
	Down  Direction = 1 << 0
	Up    Direction = 1 << 1
	Right Direction = 1 << 2
	Left  Direction = 1 << 3
)

const (
	Empty              byte = '.'
	MirrorSlash        byte = '/'
	MirrorBackSlash    byte = '\\'
	SplitterVertical   byte = '|'
	SplitterHorizontal byte = '-'
)

type Direction int

func (d Direction) Add(dir Direction) Direction { return d | dir }
func (d Direction) Contains(dir Direction) bool { return d&dir == dir }
func (d Direction) IsHorizontal() bool          { return d == Right || d == Left }
func (d Direction) IsVertical() bool            { return d == Down || d == Up }
func (d Direction) String() string {
	switch d {
	case Down:
		return "Down"
	case Up:
		return "Up"
	case Right:
		return "Right"
	case Left:
		return "Left"
	default:
		return "Unknown"
	}
}

func ParseGrid(lines []string) *Grid {
	grid := lib.NewByteGridFromStrings(lines)
	return &Grid{
		Grid:           lib.NewByteGridFromStrings(lines),
		LaserDirection: lib.NewGrid[Direction](grid.Rows(), grid.Columns()),
	}
}

func (g *Grid) FollowAndMark(row, col int, dir Direction) {
	if !g.IsValidPosition(row, col) {
		// laser leaving the grid is ignored
		return
	}
	ld := g.LaserDirection.Get(row, col)
	if ld.Contains(dir) {
		if VERBOSE >= 2 {
			fmt.Printf("(%d,%d) already visited %s\n", row, col, ld)
		}
		return // already marked this direction
	}
	ld = ld.Add(dir)
	g.LaserDirection.Set(row, col, ld)
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
			rowA, colA := NextPosition(row, col, Up)
			g.FollowAndMark(rowA, colA, Up)
			rowB, colB := NextPosition(row, col, Down)
			g.FollowAndMark(rowB, colB, Down)
		}
	case SplitterHorizontal: // -
		if dir.IsHorizontal() {
			// beam continues on its way
			row, col = NextPosition(row, col, dir)
			g.FollowAndMark(row, col, dir)
		} else {
			// beam is split
			rowA, colA := NextPosition(row, col, Right)
			g.FollowAndMark(rowA, colA, Right)
			rowB, colB := NextPosition(row, col, Left)
			g.FollowAndMark(rowB, colB, Left)
		}

	case MirrorSlash: // /
		dir = RedirectSlash(dir)
		row, col = NextPosition(row, col, dir)
		g.FollowAndMark(row, col, dir)

	case MirrorBackSlash: // \
		dir = RedirectBashSlash(dir)
		row, col = NextPosition(row, col, dir)
		g.FollowAndMark(row, col, dir)

	default:
		panic(fmt.Errorf("invalid field: %c", field))
	}
}

func RedirectSlash(dir Direction) Direction {
	// /
	switch dir {
	case Up:
		return Right
	case Down:
		return Left
	case Right:
		return Up
	case Left:
		return Down
	default:
		panic("invalid direction")
	}
}

func RedirectBashSlash(dir Direction) Direction {
	// \
	switch dir {
	case Up:
		return Left
	case Down:
		return Right
	case Right:
		return Down
	case Left:
		return Up
	default:
		panic("invalid direction")
	}
}

func NextPosition(row, col int, dir Direction) (int, int) {
	switch dir {
	case Up:
		return row - 1, col
	case Down:
		return row + 1, col
	case Right:
		return row, col + 1
	case Left:
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
			if g.LaserDirection.Get(r, c) != 0 {
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
			g.LaserDirection.Set(r, c, 0)
		}
	}
}

func (g *Grid) PrintEnergized() {
	rows, cols := g.Rows(), g.Columns()
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if g.LaserDirection.Get(r, c) != 0 {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}
