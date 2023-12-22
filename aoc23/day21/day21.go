package main

// https://adventofcode.com/2023/day/21

import (
	"fmt"

	"git.bind.ch/phil/challenges/lib"
	"git.bind.ch/phil/challenges/lib/rowcol"
)

var VERBOSE = 1

func main() {
	ProcessPart1("aoc23/day21/example.txt", 6) // 16
	ProcessPart1("aoc23/day21/input.txt", 64)  // 3642

	//ProcessPart2("aoc23/day21/example.txt", 6) // 16
	//ProcessPart2("aoc23/day21/example.txt", 10) // 50 - validated
	//ProcessPart2("aoc23/day21/example.txt", 50)  // 1594 - validated
	//ProcessPart2("aoc23/day21/example.txt", 100) // 6536 - validated
	//ProcessPart2("aoc23/day21/example.txt", 500) // 167004 - validated
	//ProcessPart2("aoc23/day21/example.txt", 1000) // 668697 - validated
	//ProcessPart2("aoc23/day21/example.txt", 5000)   // 16733044 - not validated
	ProcessPart2("aoc23/day21/input.txt", 26501365) //  608603023105276
}

func ProcessPart1(name string, steps int) {
	fmt.Println("Part 1 input:", name)
	lines := lib.ReadLines(name)
	g := ParseInput(lines)
	positions := SolvePart1(g, steps)
	fmt.Println("Positions:", positions)

	fmt.Println()
}

func ProcessPart2(name string, steps int) {
	fmt.Println("Part 2 input:", name)
	lines := lib.ReadLines(name)
	g := ParseInput(lines)
	positions := SolvePart2(g, steps)
	fmt.Println("Steps:", steps)
	fmt.Println("Positions:", positions)

	fmt.Println()
}

////////////////////////////////////////////////////////////

type Grid struct {
	rowcol.Grid[byte]
}

func (g *Grid) Copy() *Grid {
	return &Grid{g.Grid.Copy()}
}

func (g *Grid) MarkIfPossible(r, c int) {
	if !g.IsValidPosition(r, c) {
		return
	}
	if g.Get(r, c) != '#' {
		g.Set(r, c, 'O')
	}
}

func (g *Grid) Print() {
	for _, row := range g.Data {
		fmt.Println(string(row))
	}
	fmt.Println()
}

func ParseInput(lines []string) *Grid {
	return &Grid{rowcol.NewByteGridFromStrings(lines)}
}

func SolvePart1(g *Grid, steps int) int {
	blank := g.Copy()
	blank.Grid.Map(removeStart)
	g.Map(startToMarked)
	current := g

	for i := 0; i < steps; i++ {
		next := blank.Copy()
		FillInStep(current, next)
		current = next
		if steps < 10 {
			current.Print()
		}
	}

	return current.Count(func(v byte) bool { return v == 'O' })
}

func FillInStep(current *Grid, next *Grid) {
	rows, cols := current.Rows(), current.Columns()
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			v := current.Get(r, c)
			if v != 'O' {
				continue
			}
			next.MarkIfPossible(r+1, c)
			next.MarkIfPossible(r-1, c)
			next.MarkIfPossible(r, c+1)
			next.MarkIfPossible(r, c-1)
		}
	}
}

func removeStart(input byte) byte {
	if input == 'S' {
		return '.'
	}
	return input
}
func startToMarked(input byte) byte {
	if input == 'S' {
		return 'O'
	}
	return input
}

////////////////////////////////////////////////////////////

func SolvePart2(g *Grid, steps int) int {
	var ps map[rowcol.Pos]bool = make(map[rowcol.Pos]bool)
	start := FindStart(g)
	ps[start] = true

	//var positions = []int{1}
	var points []int

	for i := 1; i <= steps; i++ {
		ps = FillInStepPart2(g, ps)
		if VERBOSE >= 2 {
			fmt.Printf("after step %d: %d positions\n", i, len(ps))
		}

		// record number of possibilities each time a complete grid is traversed
		if i%g.Rows() == steps%g.Rows() {
			points = append(points, len(ps))
			fmt.Println("point", len(points), len(ps), "after", i, "steps")
		}
		if len(points) == 3 {
			break
		}
	}
	// the original code is not my solutions, but based on a reddit post.
	// I only analysed the formula for fitting the quadratic equation:
	//
	// f(n), where n is is the number of grids traversed
	// find a, b, c of a quadratic equation with values:
	// f(0) = y0 = 3776
	// f(1) = y1 = 33652
	// f(2) = y2 = 93270
	//
	// ax^2 + bx + c = 0
	//
	// a*0^2 + b*0 + c = y0
	// a*1^2 + b*1 + c = y1
	// a*2^2 + b*2 + c = y2
	//
	// a*0^2 + b*0 + c = y0
	// 0    +  0   + c = y0
	//               c = y0
	//
	// a*1^2 + b*1 + c  = y1
	// a     + b   + y0 = y1
	//         b        = y1 - a - y0
	//
	// a*2^2 + b*2             + c  = y2
	// 4a    + 2 (y1 - a - y0) + y0 = y2
	// 4a    + 2y1 - 2a - 2y0  + y0 = y2
	// 2a    + 2y1      -  y0       = y2
	//  a  = (y2 - 2y1 + y0) / 2
	n := steps / g.Rows()
	y0 := points[0]
	y1 := points[1]
	y2 := points[2]
	a := (y2 - 2*y1 + y0) / 2
	b := y1 - a - y0
	c := y0
	return a*n*n + b*n + c
}

var offR = [4]int{0, 0, 1, -1}
var offC = [4]int{1, -1, 0, 0}

func FillInStepPart2(g *Grid, ps map[rowcol.Pos]bool) map[rowcol.Pos]bool {
	next := make(map[rowcol.Pos]bool)
	rows, cols := g.Rows(), g.Columns()
	for pos, _ := range ps {
		for i := 0; i < 4; i++ {
			nextpos := rowcol.Pos{Row: pos.Row + offR[i], Col: pos.Col + offC[i]}
			testpos := rowcol.Pos{
				Row: mod(nextpos.Row, rows),
				Col: mod(nextpos.Col, cols),
			}
			if g.Get(testpos.Row, testpos.Col) != '#' {
				next[nextpos] = true
			}
		}
	}
	return next
}

func mod(x, y int) int {
	m := x % y
	if m < 0 {
		return m + y
	}
	return m
}

func FindStart(g *Grid) rowcol.Pos {
	rows, cols := g.Rows(), g.Columns()
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if g.Get(r, c) == 'S' {
				return rowcol.Pos{Row: r, Col: c}
			}
		}
	}
	panic("start not found")
}
