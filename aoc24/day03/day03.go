package main

// https://adventofcode.com/2024/day/3

import (
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/phicode/challenges/lib"
)

func main() {
	lib.Timed("Part 1", ProcessPart1, "aoc24/day03/example.txt")
	lib.Timed("Part 1", ProcessPart1, "aoc24/day03/input.txt")

	lib.Timed("Part 2", ProcessPart2, "aoc24/day03/example2.txt")
	lib.Timed("Part 2", ProcessPart2, "aoc24/day03/input.txt")
}

func ProcessPart1(name string) {
	fmt.Println("Part 1 input:", name)
	input := ParseInput(name)
	result := SolvePart1(input)
	fmt.Println("Result:", result)
}

func ProcessPart2(name string) {
	fmt.Println("Part 2 input:", name)
	input := ParseInput(name)
	result := SolvePart2(input)
	fmt.Println("Result:", result)
}

////////////////////////////////////////////////////////////

type Input struct {
	Muls    []Mul
	Enabled []Enabled
}

type Mul struct {
	X, Y  int
	start int
}

type Enabled struct {
	Start   int
	Enabled bool
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

func ParseInput(name string) Input {
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
		rv.Enabled = append(rv.Enabled, Enabled{index, true})
	}
	for _, index := range dontIndexes {
		rv.Enabled = append(rv.Enabled, Enabled{index, false})
	}
	sort.Slice(rv.Enabled, func(i, j int) bool { return rv.Enabled[i].Start < rv.Enabled[j].Start })
	return rv
}

////////////////////////////////////////////////////////////

func SolvePart1(input Input) int {
	total := 0
	for _, m := range input.Muls {
		//fmt.Println(m.X, "*", m.Y)
		total += m.X * m.Y
	}
	return total
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
	idx, _ := slices.BinarySearchFunc(i.Enabled, start, func(dd Enabled, x int) int {
		if dd.Start == x {
			return 0
		}
		if dd.Start < x {
			return -1
		}
		return 1
	})
	// we are searching for a value that does not exist,
	// so binary search returns the "insert" position, which is one index above the 'Enabled' that
	// interests us
	idx--
	if idx < 0 {
		// DO
		return true
	}
	dodont := i.Enabled[idx]
	return dodont.Enabled
}
