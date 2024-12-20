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
	lib.Timed("Part 2", ProcessPart2, "aoc24/day15/example_p2.txt")
	lib.Timed("Part 2", ProcessPart2, "aoc24/day15/example.txt")
	lib.Timed("Part 2", ProcessPart2, "aoc24/day15/input.txt")

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

type P1Solver struct {
	rowcol.Grid[byte]
}

func SolvePart1(input Input) int {
	robot, ok := input.grid.FindFirst(func(v byte) bool { return v == Robot })
	assert.True(ok)

	solver := P1Solver{input.grid}
	for _, d := range input.moves {
		robot = solver.Push(robot, d)
	}

	if lib.LogLevel >= lib.LogDebug {
		rowcol.PrintByteGrid(input.grid)
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

	BoxL byte = '['
	BoxR byte = ']'
)

func (in P1Solver) Push(robot rowcol.Pos, d rowcol.Direction) rowcol.Pos {
	next := robot.AddDir(d)
	v := in.GetPos(next)
	switch v {
	case Wall:
		return robot
	case Free:
		in.SetPos(next, Robot)
		in.SetPos(robot, Free)
		return next
	case Box:
		if in.PushBoxes(next, d) {
			in.SetPos(next, Robot)
			in.SetPos(robot, Free)
			return next
		}
		return robot
	default:
		panic("invalid state")
	}
}

func (in P1Solver) PushBoxes(p rowcol.Pos, d rowcol.Direction) bool {
	next := p.AddDir(d)
	v := in.GetPos(next)
	switch v {
	case Wall:
		return false
	case Free:
		in.SetPos(next, Box)
		in.SetPos(p, Free)
		return true
	case Box:
		if in.PushBoxes(next, d) {
			in.SetPos(next, Box)
			in.SetPos(p, Free)
			return true
		}
		return false
	default:
		panic("invalid state")
	}
}

////////////////////////////////////////////////////////////

type P2Solver struct {
	rowcol.Grid[byte]
}

func SolvePart2(input Input) int {
	grid := Part2Expand(input.grid)
	robot, ok := grid.FindFirst(func(v byte) bool { return v == Robot })
	assert.True(ok)

	if lib.LogLevel >= lib.LogDebug {
		rowcol.PrintByteGrid(grid)
	}

	solver := P2Solver{grid}
	for i, d := range input.moves {
		var modified bool
		robot, modified = solver.Push(robot, d)

		if modified && lib.LogLevel >= lib.LogDebug {
			fmt.Printf("Instruction %d: %v\n", (i + 1), d)
			rowcol.PrintByteGrid(grid)
		}
		_ = modified
	}

	total := 0
	for p, v := range grid.Iterator() {
		if v == BoxL {
			total += 100*p.Row + p.Col
		}
	}
	return total
}

func Part2Expand(grid rowcol.Grid[byte]) rowcol.Grid[byte] {
	rows, cols := grid.Size()
	exp := rowcol.NewGrid[byte](rows, cols*2)
	for p, v := range grid.Iterator() {
		a, b := expand(v)
		exp.Set(p.Row, p.Col*2, a)
		exp.Set(p.Row, p.Col*2+1, b)
	}
	return exp
}

func expand(v byte) (byte, byte) {
	switch v {
	case Wall:
		return Wall, Wall
	case Free:
		return Free, Free
	case Box:
		return BoxL, BoxR
	case Robot:
		return Robot, Free
	default:
		panic("invalid input")
	}

}

func (in P2Solver) Push(robot rowcol.Pos, d rowcol.Direction) (rowcol.Pos, bool) {
	next := robot.AddDir(d)
	v := in.GetPos(next)
	switch v {
	case Wall:
		return robot, false
	case Free:
		in.SetPos(next, Robot)
		in.SetPos(robot, Free)
		return next, false
	case BoxL, BoxR:
		if in.CanPushBox(next, d) {
			in.PushBox(next, d)
			assert.True(in.GetPos(next) == Free)
			in.SetPos(next, Robot)
			in.SetPos(robot, Free)
			return next, true
		}
		return robot, false
	default:
		panic("invalid state")
	}
}

func (in P2Solver) CanPushBox(p rowcol.Pos, d rowcol.Direction) bool {
	if in.GetPos(p) == Free {
		return true
	}
	if in.GetPos(p) == Wall {
		return false
	}
	if in.GetPos(p) == BoxR {
		p.Col-- //
	}
	assert.True(in.GetPos(p) == BoxL)

	if isUpOrDown(d) {
		nextL := p.AddDir(d)
		nextR := nextL.AddS(0, 1)
		vL := in.GetPos(nextL)
		vR := in.GetPos(nextR)
		if vL == Wall || vR == Wall {
			return false
		}
		if vL == Free && vR == Free {
			return true
		}
		assert.True(vL == BoxL || vL == BoxR || vL == Free)
		assert.True(vR == BoxL || vR == BoxR || vR == Free)
		return in.CanPushBox(nextL, d) && in.CanPushBox(nextR, d)
	} else { // left or right move
		nextL := p.AddDir(d)
		nextR := nextL.AddS(0, 1)

		if (d == rowcol.Left && in.GetPos(nextL) == Free) ||
			(d == rowcol.Right && in.GetPos(nextR) == Free) {
			return true
		}
		if d == rowcol.Left {
			return in.CanPushBox(nextL, d)
		}
		return in.CanPushBox(nextR, d)
	}
}

func isUpOrDown(d rowcol.Direction) bool {
	return d == rowcol.Up || d == rowcol.Down
}

func (in P2Solver) PushBox(p rowcol.Pos, d rowcol.Direction) {
	if in.GetPos(p) == BoxR {
		p.Col-- //
	}
	pR := p.AddS(0, 1)

	if isUpOrDown(d) {
		nextL := p.AddDir(d)
		nextR := nextL.AddS(0, 1)
		vL := in.GetPos(nextL)
		vR := in.GetPos(nextR)
		assert.True(vL != Wall && vR != Wall)
		if vL == Free && vR == Free {
			in.SetPos(nextL, BoxL)
			in.SetPos(nextR, BoxR)
			in.SetPos(p, Free)
			in.SetPos(pR, Free)
			return
		}
		assert.True(vL == BoxL || vL == BoxR || vL == Free)
		if vL != Free {
			in.PushBox(nextL, d)
			vR = in.GetPos(nextR) // re-read vR as the previous push might have moved changed it
		}
		assert.True(vR == BoxL || vR == BoxR || vR == Free)
		if vR != Free {
			in.PushBox(nextR, d)
		}
		assert.True(in.GetPos(nextL) == Free)
		assert.True(in.GetPos(nextR) == Free)
		in.SetPos(nextL, BoxL)
		in.SetPos(nextR, BoxR)
		in.SetPos(p, Free)
		in.SetPos(pR, Free)
	} else { // left or right move
		nextL := p.AddDir(d)
		nextR := nextL.AddS(0, 1)

		// free to move
		if (d == rowcol.Left && in.GetPos(nextL) == Free) ||
			(d == rowcol.Right && in.GetPos(nextR) == Free) {
			in.SetPos(p, Free)
			in.SetPos(pR, Free)
			in.SetPos(nextL, BoxL)
			in.SetPos(nextR, BoxR)
			return
		}

		if d == rowcol.Left {
			in.PushBox(nextL, d)
			assert.True(in.GetPos(nextL) == Free)
		} else {
			in.PushBox(nextR, d)
			assert.True(in.GetPos(nextR) == Free)
		}

		in.SetPos(p, Free)
		in.SetPos(pR, Free)
		in.SetPos(nextL, BoxL)
		in.SetPos(nextR, BoxR)
	}
}
