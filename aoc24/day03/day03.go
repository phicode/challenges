package main

// https://adventofcode.com/2024/day/3

import (
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"

	"git.bind.ch/phil/challenges/lib"
)

// TODO: timing boilerplate
var VERBOSE = 1

func main() {
	ProcessPart1("aoc24/day03/example.txt")
	ProcessPart1("aoc24/day03/input.txt")

	ProcessPart2("aoc24/day03/example2.txt")
	ProcessPart2("aoc24/day03/input.txt")
}

func ProcessPart1(name string) {
	fmt.Println("Part 1 input:", name)
	input := ParseInput(name)
	total := SolvePart1(input)
	fmt.Println("Total:", total)
}

func ProcessPart2(name string) {
	fmt.Println("Part 2 input:", name)
	input := ParseInputPart2(name)
	total := SolvePart2(input)
	fmt.Println("Total:", total)
}

func log(v int, msg string) {
	if v <= VERBOSE {
		fmt.Println(msg)
	}
}

////////////////////////////////////////////////////////////

type Input struct {
	Muls    []Mul
	DoDonts []DoDont
}

type Mul struct {
	X, Y  int
	start int
}

type DoDont struct {
	Start int
	Do    bool
}

func ParseInput(name string) Input {
	lines := lib.ReadLines(name)
	var rv Input
	for _, line := range lines {
		for len(line) > 0 {
			x, y, rem, ok := parseMul(line)
			if ok {
				line = rem
				rv.Muls = append(rv.Muls, Mul{x, y, 0})
			} else {
				line = line[1:]
			}
		}
	}
	return rv
}

func parseMul(s string) (x, y int, rem string, ok bool) {
	if !strings.HasPrefix(s, "mul(") {
		return 0, 0, "", false
	}
	end := strings.IndexRune(s, ')')
	if end == -1 || end-4 < 3 {
		return 0, 0, "", false
	}
	between := s[4:end]
	parts := strings.Split(between, ",")
	if len(parts) != 2 {
		return 0, 0, "", false
	}
	x, err1 := strconv.Atoi(parts[0])
	y, err2 := strconv.Atoi(parts[1])
	if err1 == nil && err2 == nil {
		return x, y, s[end+1:], true
	}
	return 0, 0, "", false
}

func SolvePart1(input Input) int {
	total := 0
	for _, m := range input.Muls {
		//fmt.Println(m.X, "*", m.Y)
		total += m.X * m.Y
	}
	return total
}

////////////////////////////////////////////////////////////

type ParseContext struct {
	line string
	next int
}

func (ctx *ParseContext) NextMul() (int, int, int, bool) {
	for ctx.next < len(ctx.line) {
		rem := ctx.line[ctx.next:]
		index := strings.Index(rem, "mul(")
		if index == -1 {
			return 0, 0, 0, false
		}
		start := ctx.next + index
		if x, y, next, ok := testmul(ctx.line[start:]); ok {
			ctx.next = start + next
			return x, y, start, true
		}
		ctx.next = start + 1
	}
	return 0, 0, 0, false
}

func testmul(s string) (x, y, next int, ok bool) {
	if !strings.HasPrefix(s, "mul(") {
		return 0, 0, 0, false
	}
	end := strings.IndexRune(s, ')')
	if end == -1 || end-4 < 3 {
		return 0, 0, 0, false
	}
	between := s[4:end]
	parts := strings.Split(between, ",")
	if len(parts) != 2 {
		return 0, 0, 0, false
	}
	x, err1 := strconv.Atoi(parts[0])
	y, err2 := strconv.Atoi(parts[1])
	if err1 == nil && err2 == nil {
		return x, y, end + 1, true
	}
	return 0, 0, 0, false
}

func ParseInputPart2(name string) Input {
	lines := lib.ReadLines(name)
	line := lib.ConcatStrings(lines)
	var rv Input
	ctxt := ParseContext{
		line: line,
	}
	for {
		x, y, start, ok := ctxt.NextMul()
		if !ok {
			break
		}
		rv.Muls = append(rv.Muls, Mul{x, y, start})
	}
	doIndexes := lib.AllStringIndexes(line, "do()")
	dontIndexes := lib.AllStringIndexes(line, "don't()")
	for _, index := range doIndexes {
		rv.DoDonts = append(rv.DoDonts, DoDont{index, true})
	}
	for _, index := range dontIndexes {
		rv.DoDonts = append(rv.DoDonts, DoDont{index, false})
	}
	sort.Slice(rv.DoDonts, func(i, j int) bool { return rv.DoDonts[i].Start < rv.DoDonts[j].Start })
	return rv
}

func SolvePart2(input Input) int {
	total := 0
	for _, m := range input.Muls {
		do := input.FindDoDont(m.start)
		if do {
			total += m.X * m.Y
		}
	}
	return total
}

func (i Input) FindDoDont(start int) bool {
	idx, _ := slices.BinarySearchFunc(i.DoDonts, start, func(dd DoDont, x int) int {
		if dd.Start == x {
			return 0
		}
		if dd.Start < x {
			return -1
		}
		return 1
	})
	// we are searching for a value that does not exist,
	// so binary search returns the "insert" position, which is one index above the 'DoDont' that
	// interests us
	idx--
	if idx < 0 {
		// DO
		return true
	}
	dodont := i.DoDonts[idx]
	return dodont.Do
}
