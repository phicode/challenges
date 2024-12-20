package main

// https://adventofcode.com/2024/day/14

import (
	"flag"
	"fmt"
	"strings"

	"github.com/phicode/challenges/lib"
	"github.com/phicode/challenges/lib/assert"
	"github.com/phicode/challenges/lib/math"
	"github.com/phicode/challenges/lib/rowcol"
)

func main() {
	flag.Parse()
	lib.Timed("Part 1", ProcessPart1Example, "aoc24/day14/example.txt")
	lib.Timed("Part 1", ProcessPart1, "aoc24/day14/input.txt")
	lib.Timed("Part 2", ProcessPart2, "aoc24/day14/input.txt")
	//lib.Profile(1, "day14.pprof", "Part 2", ProcessPart2, "aoc24/day14/input.txt")
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
	result := SolvePart2(input, 101, 103)
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
	if lib.LogLevel >= lib.LogDebug {
		render(input, xmax, ymax)
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

func SolvePart2(input Input, xmax, ymax int) int {
	loop := input.AnalyseLoops(xmax, ymax)

	original := input.copy()

	printed := 0
	for sec := 1; sec <= loop; sec++ {
		input.MoveAll(xmax, ymax)

		// use to find the "solution" visually
		//fmt.Println("Second:", sec)
		//render(input, xmax, ymax)

		if sec == 8006 && lib.LogLevel >= lib.LogDebug {
			render(input, xmax, ymax)
		}

		if input.FindBox(xmax, ymax) {
			return sec
		}
	}
	assert.True(original.eq(&input))
	fmt.Println("Printed:", printed)
	return 0
}

func render(input Input, xmax, ymax int) {
	fmt.Println(strings.Repeat("=", xmax))
	grid := rowcol.NewGrid[byte](ymax, xmax)
	grid.Reset(' ')
	for _, robot := range input.robots {
		grid.SetPos(robot.P, '#')
	}
	rowcol.PrintByteGrid(grid)
}

func (in *Input) AnalyseLoops(xmax int, ymax int) int {
	first := 0
	for i, robot := range in.robots {
		l := robot.FindLoop(xmax, ymax)
		if i == 0 {
			first = l
		}
		assert.True(l == first)
	}
	// all robots loop around in 10403 seconds
	//fmt.Println("loop analysis done:", first)
	return first
}

func (r Robot) FindLoop(xmax, ymax int) int {
	start := r
	cur := r
	sec := 0
	for {
		sec++
		cur = cur.Move(xmax, ymax)
		if cur == start {
			return sec
		}
	}
}

func (in *Input) copy() *Input {
	cpy := &Input{}
	for _, robot := range in.robots {
		cpy.robots = append(cpy.robots, robot)
	}
	return cpy
}
func (in *Input) eq(other *Input) bool {
	for i, robot := range in.robots {
		if other.robots[i] != robot {
			return false
		}
	}
	return true
}

// so this is quite bullshit, the "result" contains a "christmas tree" like:
//
//	###############################
//	#                             #
//	#                             #
//	#                             #
//	#                             #
//	#              #              #
//	#             ###             #
//	#            #####            #
//	#           #######           #
//	#          #########          #
//	#            #####            #
//	#           #######           #
//	#          #########          #
//	#         ###########         #
//	#        #############        #
//	#          #########          #
//	#         ###########         #
//	#        #############        #
//	#       ###############       #
//	#      #################      #
//	#        #############        #
//	#       ###############       #
//	#      #################      #
//	#     ###################     #
//	#    #####################    #
//	#             ###             #
//	#             ###             #
//	#             ###             #
//	#                             #
//	#                             #
//	#                             #
//	#                             #
//	###############################
//
// simplified solution: find the "box", ie: 31 consecutive filled spaces
func (in *Input) FindBox(xmax, ymax int) bool {
	grid := rowcol.NewGrid[bool](ymax, xmax)
	for _, robot := range in.robots {
		grid.SetPos(robot.P, true)
	}
	for pi := range grid.PosIterator() {
		if horiz31AreSet(grid, pi) && vert33AreSet(grid, pi) {
			return true
		}
	}
	return false
}

func horiz31AreSet(grid rowcol.Grid[bool], p rowcol.Pos) bool {
	for i := 0; i < 31; i++ {
		test := p.AddS(0, i)
		if !grid.IsValidPos(test) || !grid.GetPos(test) {
			return false
		}
	}
	return true
}
func vert33AreSet(grid rowcol.Grid[bool], p rowcol.Pos) bool {
	for i := 0; i < 33; i++ {
		test := p.AddS(i, 0)
		if !grid.IsValidPos(test) || !grid.GetPos(test) {
			return false
		}
	}
	return true
}
