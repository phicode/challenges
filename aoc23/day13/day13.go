package main

// https://adventofcode.com/2023/day/XX

import (
	"bytes"
	"fmt"

	"git.bind.ch/phil/challenges/lib"
)

// TODO: timing boilerplate
var VERBOSE = 1

func main() {
	ProcessPart1("aoc23/day13/example.txt")
	ProcessPart1("aoc23/day13/input.txt")
	//
	//ProcessPart2("aoc23/day13/example.txt")
	//ProcessPart2("aoc23/day13/input.txt")
}

func ProcessPart1(name string) {
	fmt.Println("Part 1 input:", name)
	lines := lib.ReadLines(name)
	grids := ParseGrids(lines)
	//fmt.Println("got grids:", len(grids))
	var sum int
	for i, g := range grids {
		mirrorRow := g.FindMirror()
		fmt.Println("grid", i, "mirror line:", mirrorRow)
		t := g.Transpose()
		mirrorCol := t.FindMirror()
		fmt.Println("grid", i, "mirror line t:", mirrorCol)

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
	cols := len(g[0])
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

func (g Grid) Rows() int {
	return len(g)
}

func (g Grid) FindMirror() int {
	// tests if i and i+1 mirror
	for i := 0; i < g.Rows()-1; i++ {
		if g.IsMirrorLine(i) {
			return i + 1
		}
	}
	return -1
}

func (g Grid) IsMirrorLine(i int) bool {
	length := min(i+1, g.Rows()-i-1)
	//fmt.Println("mirror check", i, ", length:", length)
	for offset := 0; offset < length; offset++ {
		//fmt.Println("comparing rows", i-offset, "and", i+1+offset)
		if !bytes.Equal(g[i-offset], g[i+1+offset]) {
			return false
		}
	}
	return true
}
