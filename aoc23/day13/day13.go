package main

// https://adventofcode.com/2023/day/13

import (
	"bytes"
	"fmt"

	"git.bind.ch/phil/challenges/lib"
)

// TODO: timing boilerplate
var VERBOSE = 2

func main() {
	ProcessPart1("aoc23/day13/example.txt")
	ProcessPart1("aoc23/day13/input.txt")

	ProcessPart2("aoc23/day13/example.txt")
	ProcessPart2("aoc23/day13/input.txt")
}

func ProcessPart1(name string) {
	fmt.Println("Part 1 input:", name)
	lines := lib.ReadLines(name)
	grids := ParseGrids(lines)
	var sum int
	for i, g := range grids {
		mirrorRow := g.FindMirror(-1)
		t := g.Transpose()
		mirrorCol := t.FindMirror(-1)
		if VERBOSE >= 2 {
			fmt.Println("grid", i, "mirror line:", mirrorRow)
			fmt.Println("grid", i, "mirror line t:", mirrorCol)
		}

		if (mirrorRow < 0 && mirrorCol < 0) || (mirrorRow >= 0 && mirrorCol >= 0) {
			panic("invalid input")
		}
		if mirrorRow >= 0 {
			sum += mirrorRow * 100
		}
		if mirrorCol >= 0 {
			sum += mirrorCol
		}
	}
	fmt.Println("sum:", sum)
	fmt.Println()
}

func ProcessPart2(name string) {
	fmt.Println("Part 2 input:", name)
	lines := lib.ReadLines(name)
	grids := ParseGrids(lines)
	var sum int
	for i, g := range grids {
		t := g.Transpose()
		//mirrorRow, mirrorCol := FindMirrorPart2(g, t)
		mirrorRow := FindMirrorPart2v2(g)
		mirrorCol := FindMirrorPart2v2(t)

		if VERBOSE >= 2 {
			fmt.Println("grid", i, "mirror line:", mirrorRow)
			fmt.Println("grid", i, "mirror line t:", mirrorCol)
		}

		if (mirrorRow < 0 && mirrorCol < 0) || (mirrorRow >= 0 && mirrorCol >= 0) {
			fmt.Println("Grid nr:", i)
			fmt.Println("Rows")
			g.Print()
			fmt.Println("row line:", g.FindMirror(-1))
			fmt.Println()
			fmt.Println("Columns")
			t.Print()
			fmt.Println("column line:", t.FindMirror(-1))
			//gridm, gridmlen := g.FindMirror()
			//tmir, tmirlen := t.FindMirror()
			//fmt.Println("grid mirror line:", gridm, gridmlen)
			//fmt.Println("transposed mirror line:", tmir, tmirlen)

			fmt.Println("found mirror lines:", mirrorRow, mirrorCol)
			panic(fmt.Errorf("invalid input with grid: %d", i+1))
		}
		if mirrorRow >= 0 {
			sum += mirrorRow * 100
		}
		if mirrorCol >= 0 {
			sum += mirrorCol
		}
	}
	fmt.Println("sum:", sum)
	fmt.Println()
}

func log(v int, msg string) {
	if v <= VERBOSE {
		fmt.Println(msg)
	}
}

////////////////////////////////////////////////////////////

func ParseGrids(lines []string) []Grid {
	var rv []Grid
	var current Grid
	for _, l := range lines {
		if l == "" {
			rv = append(rv, current)
			current = nil
			continue
		}
		current = append(current, []byte(l))
	}
	rv = append(rv, current)
	return rv
}

type Grid [][]byte

func (g Grid) Transpose() Grid {
	cols := g.Cols()
	rows := g.Rows()
	t := make(Grid, cols)
	for i := 0; i < cols; i++ {
		t[i] = make([]byte, rows)
	}
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			t[x][y] = g[y][x]
		}
	}
	return t
}

func (g Grid) Rows() int { return len(g) }
func (g Grid) Cols() int { return len(g[0]) }

func (g Grid) FindMirror(exclude int) int {
	// tests if i and i+1 are mirrors
	found, l := -1, 0
	for i := 0; i < g.Rows()-1; i++ {
		if i == exclude-1 {
			continue
		}
		ok, newLen := g.IsMirrorLine(i)
		if ok {
			if newLen > l {
				l = newLen
				found = i
			}
		}
	}
	if found != -1 {
		return found + 1
	}
	return found
}

func (g Grid) IsMirrorLine(i int) (bool, int) {
	length := min(i+1, g.Rows()-i-1)
	//fmt.Println("mirror check", i, ", length:", length)
	for offset := 0; offset < length; offset++ {
		//fmt.Println("comparing rows", i-offset, "and", i+1+offset)
		if !bytes.Equal(g[i-offset], g[i+1+offset]) {
			return false, 0
		}
	}
	return true, length
}

////////////////////////////////////////////////////////////
// Part2

//func FindMirrorPart2(g, t Grid) (int, int) {
//	// every combination of toggled grid
//	cols := g.Cols()
//	rows := g.Rows()
//	row := -1
//	col := -1
//
//	p1row := g.FindMirror()
//	p1col := t.FindMirror()
//
//	//comb := 0
//	for x := 0; x < cols; x++ {
//		for y := 0; y < rows; y++ {
//
//			g.Toggle(x, y)
//			if mirror := g.FindMirror(); mirror != -1 {
//				if p1row != mirror {
//					if row == -1 {
//						row = mirror
//					}
//					row = min(row, mirror)
//				}
//			}
//			g.Toggle(x, y)
//
//			t.Toggle(y, x)
//			if mirror := t.FindMirror(); mirror != -1 {
//				if p1col != mirror {
//					if col == -1 {
//						col = mirror
//					}
//					col = min(col, mirror)
//				}
//			}
//			t.Toggle(y, x)
//			//comb++
//		}
//	}
//	return row, col
//}

func FindMirrorPart2v2(g Grid) int {
	// every combination of toggled grid
	cols := g.Cols()
	rows := g.Rows()
	p1row := g.FindMirror(-1)
	for x := 0; x < cols; x++ {
		for y := 0; y < rows; y++ {
			g.Toggle(x, y)
			mirror := g.FindMirror(p1row)
			g.Toggle(x, y)
			if mirror != -1 {
				return mirror
			}
		}
	}
	return -1
}

func (g Grid) Toggle(x, y int) {
	v := g[y][x]
	if v == '.' {
		v = '#'
	} else {
		v = '.'
	}
	g[y][x] = v
}

func (g Grid) Print() int {
	for i, r := range g {
		fmt.Printf("%-5d %s\n", i+1, string(r))
	}
	return 0
}
