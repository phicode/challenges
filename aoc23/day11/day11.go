package main

// https://adventofcode.com/2023/day/11

import (
	"fmt"
	"sort"

	"github.com/phicode/challenges/lib"
)

// TODO: timing boilerplate
var VERBOSE = 0

func main() {
	ProcessPart1("aoc23/day11/example.txt")
	ProcessPart1("aoc23/day11/input.txt")

	ProcessPart2("aoc23/day11/example.txt", ExpandStep2Example)
	ProcessPart2("aoc23/day11/input.txt", ExpandStep2)
}

func ProcessPart1(name string) {
	fmt.Println("Part 1 input:", name)
	Solve(name, ExpandDouble)
	fmt.Println()
}

func ProcessPart2(name string, e Expansion) {
	fmt.Println("Part 2 input:", name)
	Solve(name, e)
	fmt.Println()
}

func Solve(name string, e Expansion) {
	lines := lib.ReadLines(name)
	space := ParseSpace(lines)
	space.Expand(e)
	fmt.Println("galaxies:", len(space.Galaxies))
	if VERBOSE >= 1 {
		for i, p := range space.Galaxies {
			fmt.Println("galaxy", (i + 1), ":", p)
		}
	}
	fmt.Println("sum of distance:", space.Distances())
}

////////////////////////////////////////////////////////////

func ParseSpace(lines []string) *Space {
	s := Space{}
	maxy := len(lines) - 1
	for i, line := range lines {
		y := maxy - i
		for x, v := range line {
			if v != '#' {
				continue
			}
			s.Galaxies = append(s.Galaxies, Pos{x, y})
		}
	}
	return &s
}

type Pos struct {
	X, Y int
}

func (p Pos) ManhattanDistance(q Pos) int {
	return abs(p.X-q.X) + abs(p.Y-q.Y)
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

type Space struct {
	Galaxies []Pos
}

func (s *Space) Expand(e Expansion) {
	xs := lib.Map(s.Galaxies, func(p Pos) int { return p.X })
	ys := lib.Map(s.Galaxies, func(p Pos) int { return p.Y })
	sort.Ints(xs)
	sort.Ints(ys)

	// largest to smallest
	for i := len(xs) - 1; i > 0; i-- {
		x1, x2 := xs[i-1], xs[i]
		d := x2 - x1
		if d <= 1 {
			continue
		}
		s.ExpandX(x1+1, e(d-1))
	}
	// largest to smallest
	for i := len(ys) - 1; i > 0; i-- {
		y1, y2 := ys[i-1], ys[i]
		d := y2 - y1
		if d <= 1 {
			continue
		}
		s.ExpandY(y1+1, e(d-1))
	}
}

func (s *Space) ExpandX(startX int, amount int) {
	if VERBOSE >= 2 {
		fmt.Println("expand X", startX, "by", amount)
	}
	for i, p := range s.Galaxies {
		if p.X >= startX {
			s.Galaxies[i].X += amount
		}
	}
}
func (s *Space) ExpandY(startY int, amount int) {
	if VERBOSE >= 2 {
		fmt.Println("expand Y", startY, "by", amount)
	}
	for i, p := range s.Galaxies {
		if p.Y >= startY {
			s.Galaxies[i].Y += amount
		}
	}
}

func (s *Space) VisitGalaxyPairs(f func(a, b Pos)) {
	n := len(s.Galaxies)
	for i, p := range s.Galaxies {
		for j := i + 1; j < n; j++ {
			f(p, s.Galaxies[j])
		}
	}
}

func (s *Space) Distances() int {
	var sum int
	f := func(a, b Pos) {
		sum += a.ManhattanDistance(b)
	}
	s.VisitGalaxyPairs(f)
	return sum
}

// determines how much space to add for a given amount of empty space
type Expansion func(int) int

func ExpandDouble(x int) int {
	return x
}
func ExpandStep2Example(x int) int {
	return x*10 - x
}
func ExpandStep2(x int) int {
	return x*1_000_000 - x
}
