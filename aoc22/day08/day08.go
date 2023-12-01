package main

// https://adventofcode.com/2022/day/8

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	Process("aoc22/day08/example.txt")
	Process("aoc22/day08/input.txt")
}

func Process(name string) {
	fmt.Println("input:", name)
	lines := ReadInput(name)

	g := NewGrid(lines)
	g.print()
	fmt.Println("num visible trees:", g.NumVisible())
	fmt.Println("max scenic score :", g.MaxScenicScore())

	fmt.Println()
}

func ReadInput(name string) []string {
	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	s := bufio.NewScanner(f)

	var lines []string
	for s.Scan() {
		line := s.Text()
		lines = append(lines, line)
	}

	if err := s.Err(); err != nil {
		panic(s.Err())
	}
	return lines
}

////////////////////////////////////////////////////////////

type grid struct {
	rows    int
	columns int
	data    []int
}

func NewGrid(lines []string) *grid {
	rows := len(lines)
	columns := len(lines[0])
	data := make([]int, rows*columns)
	i := 0
	for _, line := range lines {
		if len(line) != columns {
			panic("invalid line length")
		}
		for _, c := range line {
			if c < '0' || c > '9' {
				panic("invalid input")
			}
			data[i] = int(c - '0')
			i++
		}
	}
	return &grid{
		rows:    rows,
		columns: columns,
		data:    data,
	}
}

func (g *grid) idx(x, y int) int {
	return y*g.columns + x
}

func (g *grid) print() {
	for y := 0; y < g.rows; y++ {
		for x := 0; x < g.rows; x++ {
			fmt.Printf("%d", g.data[g.idx(x, y)])
		}
		fmt.Println()
	}
}

func (g *grid) NumVisible() int {
	v := g.columns*2 + (g.rows-2)*2
	for y := 1; y < g.rows-1; y++ {
		for x := 1; x < g.rows-1; x++ {
			if g.isVisible(x, y) {
				v++
			}
		}
	}
	return v
}

func (g *grid) MaxScenicScore() int {
	ss := -1
	for y := 0; y < g.rows; y++ {
		for x := 0; x < g.rows; x++ {
			score := g.score(x, y)
			ss = max(ss, score)
		}
	}
	return ss
}

func (g *grid) isVisible(x int, y int) bool {
	h := g.data[g.idx(x, y)]
	blocked := 0
	// horizontal left
	for ty := 0; ty < y; ty++ {
		if g.data[g.idx(x, ty)] >= h {
			blocked++
			break
		}
	}

	// horizontal right
	for ty := y + 1; ty < g.rows; ty++ {
		if g.data[g.idx(x, ty)] >= h {
			blocked++
			break
		}
	}

	// vertical top
	for tx := 0; tx < x; tx++ {
		if g.data[g.idx(tx, y)] >= h {
			blocked++
			break
		}
	}

	// vertical bottom
	for tx := x + 1; tx < g.columns; tx++ {
		if g.data[g.idx(tx, y)] >= h {
			blocked++
			break
		}
	}
	return blocked < 4
}

func (g *grid) score(x int, y int) int {
	h := g.data[g.idx(x, y)]

	var a, b, c, d int

	// horizontal left
	for ty := y - 1; ty >= 0; ty-- {
		a++
		if g.data[g.idx(x, ty)] >= h {
			break
		}
	}

	// horizontal right
	for ty := y + 1; ty < g.rows; ty++ {
		b++
		if g.data[g.idx(x, ty)] >= h {
			break
		}
	}

	// vertical top
	for tx := x - 1; tx >= 0; tx-- {
		c++
		if g.data[g.idx(tx, y)] >= h {
			break
		}
	}

	// vertical bottom
	for tx := x + 1; tx < g.columns; tx++ {
		d++
		if g.data[g.idx(tx, y)] >= h {
			break
		}
	}
	return a * b * c * d
}
