package main

// https://adventofcode.com/2023/day/XX

import (
	"fmt"

	"git.bind.ch/phil/challenges/lib"
	"git.bind.ch/phil/challenges/lib/rowcol"
)

var VERBOSE = 1

func main() {
	ProcessPart1("aoc23/day21/example.txt", 6) // 16
	ProcessPart1("aoc23/day21/input.txt", 64)  // 3642

	ProcessPart2("aoc23/day21/example.txt", 6) // 16
	//ProcessPart2("aoc23/day21/example.txt", 10) // 50
	//ProcessPart2("aoc23/day21/example.txt", 50) // 1594
	//ProcessPart2("aoc23/day21/example.txt", 100) // 6536
	//ProcessPart2("aoc23/day21/example.txt", 500) // 167004
	//ProcessPart2("aoc23/day21/example.txt", 1000) // 668697
	//ProcessPart2("aoc23/day21/example.txt", 5000) // 16733044
	//ProcessPart2("aoc23/day21/input.txt", 26501365) // ?
}

func ProcessPart1(name string, steps int) {
	fmt.Println("Part 1 input:", name)
	lines := lib.ReadLines(name)
	g := ParseInput(lines)
	plots := SolvePart1(g, steps)
	fmt.Println("Plots:", plots)

	fmt.Println()
}

func ProcessPart2(name string, steps int) {
	fmt.Println("Part 2 input:", name)
	lines := lib.ReadLines(name)
	g := ParseInput(lines)
	plots := SolvePart2(g, steps)
	fmt.Println("Steps:", steps)
	fmt.Println("Plots:", plots)

	fmt.Println()
}

////////////////////////////////////////////////////////////

type Grid struct {
	rowcol.Grid[byte]
}
type IntGrid struct {
	rowcol.Grid[int]
}

func (g *Grid) Copy() *Grid {
	return &Grid{g.Grid.Copy()}
}

func (g *IntGrid) Copy() *IntGrid {
	return &IntGrid{g.Grid.Copy()}
}

func (g *IntGrid) AddRelative(r, c int, v int) {
	rows, cols := g.Rows(), g.Columns()
	if r < 0 {
		r += rows
	}
	if r >= rows {
		r -= rows
	}
	if c < 0 {
		c += cols
	}
	if c >= cols {
		c -= cols
	}
	current := g.Get(r, c)
	if current == -1 {
		return
	}

	g.Set(r, c, v+current)

}

func (g *IntGrid) Print() {
	for _, row := range g.Data {
		for _, v := range row {
			if v > 9 {
				fmt.Print("+")
			} else if v == -1 {
				fmt.Print("#")
			} else if v == 0 {
				fmt.Print(".")
			} else {
				fmt.Printf("%d", v)
			}
		}
		fmt.Println()
	}
	fmt.Println()
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
	current := CreatePart2Grid(g)
	blank := current.Copy()
	blank.Grid.Map(removeStartInt)

	fmt.Println("Start")
	current.Print()

	for i := 0; i < steps; i++ {
		next := blank.Copy()
		FillInStepPart2(current, next)
		current = next

		current.Print()
	}

	return current.Reduce(0, func(acc, v int) int {
		if v > 0 {
			return acc + v
		}
		return acc
	})
}

func CreatePart2Grid(g *Grid) *IntGrid {
	ig := &IntGrid{rowcol.NewGrid[int](g.Rows(), g.Columns())}
	rows, cols := g.Rows(), g.Columns()
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			v := g.Get(r, c)
			if v == '#' {
				ig.Set(r, c, -1)
			} else if v == 'S' {
				ig.Set(r, c, 1)
			}
		}
	}
	return ig
}

func removeStartInt(input int) int {
	if input == 1 {
		return 0
	}
	return input
}

func FillInStepPart2(current *IntGrid, next *IntGrid) {
	rows, cols := current.Rows(), current.Columns()
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			v := current.Get(r, c)
			if v <= 0 {
				continue
			}
			next.AddRelative(r+1, c, v)
			next.AddRelative(r-1, c, v)
			next.AddRelative(r, c+1, v)
			next.AddRelative(r, c-1, v)
		}
	}
}
