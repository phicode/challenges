package main

// https://adventofcode.com/2024/day/14

import (
	"flag"
	"fmt"

	"github.com/phicode/challenges/lib"
	"github.com/phicode/challenges/lib/assert"
	"github.com/phicode/challenges/lib/math"
	"github.com/phicode/challenges/lib/rowcol"
)

func main() {
	flag.Parse()
	lib.Timed("Part 1", ProcessPart1Example, "aoc24/day14/example.txt")
	lib.Timed("Part 1", ProcessPart1, "aoc24/day14/input.txt")

	lib.Timed("Part 2", ProcessPart2, "aoc24/day14/example.txt")
	lib.Timed("Part 2", ProcessPart2, "aoc24/day14/input.txt")

	//lib.Profile(1, "day14.pprof", "Part 2", ProcessPart2, "aoc24/dayXX/input.txt")
}

func ProcessPart1Example(name string) {
	fmt.Println("Part 1 input:", name)
	input := ReadAndParseInput(name)
	result := SolvePart1(input, 11, 7)
	fmt.Println("Result:", result)
}

func ProcessPart1(name string) {
	fmt.Println("Part 1 input:", name)
	input := ReadAndParseInput(name)
	result := SolvePart1(input, 101, 103)
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
	robots []Robot
}

type Robot struct {
	P rowcol.Pos
	V rowcol.Pos
}

func ReadAndParseInput(name string) Input {
	lines := lib.ReadLines(name)
	return ParseInput(lines)
}

func ParseInput(lines []string) Input {
	var robots []Robot
	for _, line := range lines {
		var r Robot
		n, err := fmt.Sscanf(line, "p=%d,%d v=%d,%d",
			&r.P.Col, &r.P.Row, &r.V.Col, &r.V.Row)
		assert.True(n == 4 && err == nil)
		robots = append(robots, r)
	}
	return Input{robots}
}

////////////////////////////////////////////////////////////

func SolvePart1(input Input, xmax, ymax int) int {
	for second := 1; second <= 100; second++ {
		input.MoveAll(xmax, ymax)
	}
	tl, tr, bl, br := input.CountQuadrants(xmax, ymax)
	return tl * tr * bl * br
}

func (in *Input) MoveAll(xmax, ymax int) {
	for i, robot := range in.robots {
		in.robots[i] = robot.Move(xmax, ymax)
	}
}

func (r Robot) Move(xmax, ymax int) Robot {
	p := r.P.Add(r.V)
	p.Row = math.ModUnsigned(p.Row, ymax)
	p.Col = math.ModUnsigned(p.Col, xmax)
	return Robot{
		P: p,
		V: r.V,
	}
}

func (in *Input) CountQuadrants(xmax, ymax int) (int, int, int, int) {
	assert.True(xmax%2 == 1)
	assert.True(ymax%2 == 1)
	// top-left, top-right, bottom-left, bottom-right
	var tl, tr, bl, br int
	midX := xmax / 2
	midY := ymax / 2
	for _, robot := range in.robots {
		xdiff := compare(robot.P.Col, midX)
		ydiff := compare(robot.P.Row, midY)
		if xdiff == 0 || ydiff == 0 {
			continue // in the middle -> dont count
		}
		switch {
		case xdiff == -1 && ydiff == -1:
			tl++
		case xdiff == -1 && ydiff == 1:
			bl++
		case xdiff == 1 && ydiff == -1:
			tr++
		case xdiff == 1 && ydiff == 1:
			br++
		default:
			panic("invalid state")
		}
	}
	return tl, tr, bl, br
}

func compare(a int, b int) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

////////////////////////////////////////////////////////////

func SolvePart2(input Input) int {
	return 0
}
