package main

// https://adventofcode.com/2024/day/10

import (
	"flag"
	"fmt"

	"github.com/phicode/challenges/lib"
	"github.com/phicode/challenges/lib/rowcol"
)

func main() {
	flag.Parse()
	lib.Timed("Part 1", ProcessPart1, "aoc24/day10/example.txt")
	lib.Timed("Part 1", ProcessPart1, "aoc24/day10/input.txt")

	lib.Timed("Part 2", ProcessPart2, "aoc24/day10/example.txt")
	lib.Timed("Part 2", ProcessPart2, "aoc24/day10/input.txt")

	//lib.Profile(1, "day10.pprof", "Part 2", ProcessPart2, "aoc24/dayXX/input.txt")
}

func ProcessPart1(name string) {
	fmt.Println("Part 1 input:", name)
	input := ReadAndParseInput(name)
	result := SolvePart1(input)
	fmt.Println("Result:", result)
}

func ProcessPart2(name string) {
	fmt.Println("Part 2 input:", name)
	input := ReadAndParseInput(name)
	result := SolvePart2(input)
	fmt.Println("Result:", result)
}

////////////////////////////////////////////////////////////

type Input struct {
	grid rowcol.Grid[byte]
}

func ReadAndParseInput(name string) *Input {
	lines := lib.ReadLines(name)
	return ParseInput(lines)
}

func ParseInput(lines []string) *Input {
	grid := rowcol.NewByteGridFromStrings(lines)
	return &Input{grid}
}

////////////////////////////////////////////////////////////

type Visitor interface {
	Visit(pos rowcol.Pos, value byte) (end bool, result int)
	SetResult(pos rowcol.Pos, result int)
	Reset()
}

func traverse(g rowcol.Grid[byte], pos rowcol.Pos, value byte, v Visitor) int {
	if !g.IsValidPos(pos) || g.GetPos(pos) != value {
		return 0
	}
	if end, result := v.Visit(pos, value); end {
		return result
	}
	result := traverse(g, pos.AddDir(rowcol.Up), value+1, v) +
		traverse(g, pos.AddDir(rowcol.Down), value+1, v) +
		traverse(g, pos.AddDir(rowcol.Left), value+1, v) +
		traverse(g, pos.AddDir(rowcol.Right), value+1, v)
	v.SetResult(pos, result)
	return result
}

func Solve(input *Input, v Visitor) int {
	trailheads := input.grid.FindAll(func(x byte) bool { return x == '0' })
	sum := 0
	for _, th := range trailheads {
		sum += traverse(input.grid, th, byte('0'), v)
		v.Reset()
	}
	return sum
}

////////////////////////////////////////////////////////////

type P1State struct {
	visited rowcol.Grid[bool]
}

func (p *P1State) Visit(pos rowcol.Pos, value byte) (bool, int) {
	if p.visited.GetPos(pos) {
		return true, 0
	}
	p.visited.SetPos(pos, true)
	if value == '9' {
		return true, 1
	}
	return false, 0
}
func (p *P1State) SetResult(pos rowcol.Pos, result int) {}
func (p *P1State) Reset() {
	p.visited.Reset(false)
}

func SolvePart1(input *Input) int {
	p := &P1State{rowcol.NewGrid[bool](input.grid.Size())}
	return Solve(input, p)
}

////////////////////////////////////////////////////////////

type P2State struct {
	ways rowcol.Grid[int]
}

func (p *P2State) Visit(pos rowcol.Pos, value byte) (bool, int) {
	if cached := p.ways.GetPos(pos); cached > 0 {
		return true, cached
	}
	if value == '9' {
		return true, 1
	}
	return false, 0
}
func (p *P2State) SetResult(pos rowcol.Pos, result int) {
	p.ways.SetPos(pos, result)
}
func (p *P2State) Reset() {
	p.ways.Reset(0)
}

func SolvePart2(input *Input) int {
	p := &P2State{rowcol.NewGrid[int](input.grid.Size())}
	return Solve(input, p)
}
