package main

// https://adventofcode.com/2024/day/15

import (
	"flag"
	"fmt"
	"slices"

	"github.com/phicode/challenges/lib"
	"github.com/phicode/challenges/lib/assert"
	"github.com/phicode/challenges/lib/rowcol"
)

func main() {
	flag.Parse()
	lib.Timed("Part 1", ProcessPart1, "aoc24/day15/example_small.txt")
	lib.Timed("Part 1", ProcessPart1, "aoc24/day15/example.txt")
	lib.Timed("Part 1", ProcessPart1, "aoc24/day15/input.txt")
	//
	//lib.Timed("Part 2", ProcessPart2, "aoc24/day15/example.txt")
	//lib.Timed("Part 2", ProcessPart2, "aoc24/day15/input.txt")

	//lib.Profile(1, "day15.pprof", "Part 2", ProcessPart2, "aoc24/dayXX/input.txt")
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
	grid  rowcol.Grid[byte]
	moves []rowcol.Direction
}

func ReadAndParseInput(name string) Input {
	lines := lib.ReadLines(name)
	return ParseInput(lines)
}

func ParseInput(lines []string) Input {
	// find first newline, everything before is the input grid
	dividerIndex := slices.Index(lines, "")
	assert.True(dividerIndex > 0)
	grid := rowcol.NewByteGridFromStrings(lines[:dividerIndex])
	instructionLines := lines[dividerIndex+1:]
	var moves []rowcol.Direction
	for _, instrLine := range instructionLines {
		for _, move := range instrLine {
			moves = append(moves, rowcol.ParseDirectionByte(byte(move)))
		}
	}
	return Input{grid, moves}
}

////////////////////////////////////////////////////////////

func SolvePart1(input Input) int {
	robot, ok := input.grid.FindFirst(func(v byte) bool { return v == Robot })
	assert.True(ok)

	for _, d := range input.moves {
		robot = input.Push(robot, d)
	}

	if lib.LogLevel >= lib.LogDebug {
		rowcol.PrintByteGrid(&input.grid)
	}

	total := 0
	for p, v := range input.grid.Iterator() {
		if v == Box {
			total += 100*p.Row + p.Col
		}
	}

	return total
}

const (
	Wall  byte = '#'
	Free  byte = '.'
	Box   byte = 'O'
	Robot byte = '@'
)

func (in Input) Push(robot rowcol.Pos, d rowcol.Direction) rowcol.Pos {
	next := robot.AddDir(d)
	v := in.grid.GetPos(next)
	switch v {
	case Wall:
		return robot
	case Free:
		in.grid.SetPos(next, Robot)
		in.grid.SetPos(robot, Free)
		return next
	case Box:
		if in.PushBoxes(next, d) {
			in.grid.SetPos(next, Robot)
			in.grid.SetPos(robot, Free)
			return next
		}
		return robot
	default:
		panic("invalid state")
	}
}

func (in Input) PushBoxes(p rowcol.Pos, d rowcol.Direction) bool {
	next := p.AddDir(d)
	v := in.grid.GetPos(next)
	switch v {
	case Wall:
		return false
	case Free:
		in.grid.SetPos(next, Box)
		in.grid.SetPos(p, Free)
		return true
	case Box:
		if in.PushBoxes(next, d) {
			in.grid.SetPos(next, Box)
			in.grid.SetPos(p, Free)
			return true
		}
		return false
	default:
		panic("invalid state")
	}
}

////////////////////////////////////////////////////////////

func SolvePart2(input Input) int {
	return 0
}
