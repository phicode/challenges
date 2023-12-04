package main

// https://adventofcode.com/2023/day/3

import (
	"fmt"

	"git.bind.ch/phil/challenges/lib"
)

func main() {
	ProcessStep1("aoc23/day03/example.txt")
	ProcessStep1("aoc23/day03/input.txt")

	ProcessStep2("aoc23/day03/example.txt")
	ProcessStep2("aoc23/day03/input.txt")
}

func ProcessStep1(name string) {
	fmt.Println("input:", name)
	lines := lib.ReadLines(name)
	grid := Grid(lines)
	s := NewSolver(grid)
	grid.VisitSymbols(s.SymbolVisitor)
	fmt.Println("found number:", len(s.numbers))
	fmt.Println("sum of numbers:", s.NumberSum())

	fmt.Println()
}

func ProcessStep2(name string) {
	fmt.Println("input:", name)
	lines := lib.ReadLines(name)
	_ = lines

	grid := Grid(lines)
	s := NewSolverStep2(grid)
	grid.VisitSymbols(s.SymbolVisitor)
	fmt.Println("sum:", s.sum)

	fmt.Println()
}

////////////////////////////////////////////////////////////

type Grid []string
type Pos struct {
	X, Y int
}

func (g Grid) Width() int {
	return len(g[0])
}

func (g Grid) Height() int {
	return len(g)
}

func (g Grid) Size() (w, h int) {
	return g.Width(), g.Height()
}

func (g Grid) VisitSymbols(f func(g Grid, x, y int)) {
	w, h := g.Size()
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			if IsSymbol(g[y][x]) {
				f(g, x, y)
			}
		}
	}
}

func (g Grid) Get(x, y int) uint8 {
	return g[y][x]
}
func (g Grid) GetDigit(x, y int) int {
	v := g[y][x]
	if !IsDigit(v) {
		panic("not a digit")
	}
	return int(v - '0')
}

// return the extracted number and the starting position
func (g Grid) ExtractNumber(x, y int) (int, Pos) {
	// find start position
	for (x-1) >= 0 && IsDigit(g.Get(x-1, y)) {
		x--
	}
	start := Pos{x, y}

	// accumulate number
	var number int = g.GetDigit(x, y)
	w := g.Width()
	for x+1 < w && IsDigit(g.Get(x+1, y)) {
		x++
		number *= 10
		number += g.GetDigit(x, y)
	}
	return number, start
}

func IsSymbol(c uint8) bool {
	return c != '.' && !IsDigit(c)
}

func IsDigit(c uint8) bool {
	return c >= '0' && c <= '9'
}

type Solver struct {
	grid           Grid
	startPositions map[Pos]bool
	numbers        []int
}

func NewSolver(g Grid) *Solver {
	return &Solver{
		grid:           g,
		startPositions: make(map[Pos]bool),
	}
}

func (s *Solver) SymbolVisitor(g Grid, x, y int) {
	// look left
	s.Test(x-1, y-1)
	s.Test(x-1, y)
	s.Test(x-1, y+1)

	// look right
	s.Test(x+1, y-1)
	s.Test(x+1, y)
	s.Test(x+1, y+1)

	// top and bottom
	s.Test(x, y-1)
	s.Test(x, y+1)
}

func (s *Solver) Test(x, y int) {
	w, h := s.grid.Size()
	if x < 0 || x >= w || y < 0 || y >= h {
		return
	}
	if IsDigit(s.grid.Get(x, y)) {
		s.ExtractNumber(x, y)
	}
}

func (s *Solver) ExtractNumber(x, y int) {
	num, start := s.grid.ExtractNumber(x, y)
	if _, found := s.startPositions[start]; found {
		return
	}
	s.startPositions[start] = true
	s.numbers = append(s.numbers, num)
}

func (s *Solver) NumberSum() int {
	var sum int
	for _, v := range s.numbers {
		sum += v
	}
	return sum
}

////////////////////////////////////////////////////////////

type SolverStep2 struct {
	grid Grid
	sum  int
}

func NewSolverStep2(g Grid) *SolverStep2 {
	return &SolverStep2{
		grid: g,
	}
}

func (s *SolverStep2) SymbolVisitor(g Grid, x, y int) {
	// only consider gears
	if g.Get(x, y) != '*' {
		return
	}

	acc := make(map[Pos]int)

	// look left
	s.Test(acc, x-1, y-1)
	s.Test(acc, x-1, y)
	s.Test(acc, x-1, y+1)

	// look right
	s.Test(acc, x+1, y-1)
	s.Test(acc, x+1, y)
	s.Test(acc, x+1, y+1)

	// top and bottom
	s.Test(acc, x, y-1)
	s.Test(acc, x, y+1)

	if len(acc) == 2 {
		product := 1
		for _, num := range acc {
			product *= num
		}
		s.sum += product
	}
}

func (s *SolverStep2) Test(acc map[Pos]int, x, y int) {
	w, h := s.grid.Size()
	if x < 0 || x >= w || y < 0 || y >= h {
		return
	}
	if IsDigit(s.grid.Get(x, y)) {
		num, pos := s.grid.ExtractNumber(x, y)
		acc[pos] = num
	}
}
