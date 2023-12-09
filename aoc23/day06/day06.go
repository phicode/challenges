package main

// https://adventofcode.com/2023/day/XX

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"git.bind.ch/phil/challenges/lib"
)

func main() {
	ProcessPart1("aoc23/day06/example.txt")
	ProcessPart1("aoc23/day06/input.txt")

	ProcessPart2("aoc23/day06/example.txt")
	ProcessPart2("aoc23/day06/input.txt")
}

func ProcessPart1(name string) {
	fmt.Println("input:", name)
	lines := lib.ReadLines(name)
	input := ParseInput(lines)
	input.Run()
	fmt.Println()
}

func ProcessPart2(name string) {
	fmt.Println("input:", name)
	lines := lib.ReadLines(name)
	time, distance := ParseInput2(lines)
	better := SolveSample(time, distance)
	fmt.Println("Better:", better)
	fmt.Println()
}

////////////////////////////////////////////////////////////

type Input struct {
	Time     []int
	Distance []int
}

func ParseInput(lines []string) Input {
	if len(lines) != 2 {
		panic("invalid input")
	}
	t, ok := strings.CutPrefix(lines[0], "Time:")
	if !ok {
		panic("invalid time line")
	}
	d, ok := strings.CutPrefix(lines[1], "Distance:")
	if !ok {
		panic("invalid distance line")
	}
	in := Input{
		lib.ExtractInts(t),
		lib.ExtractInts(d),
	}
	if len(in.Time) != len(in.Distance) || len(in.Time) == 0 {
		panic("input lists are not of equal length")
	}
	return in
}

func ParseInput2(lines []string) (int, int) {
	if len(lines) != 2 {
		panic("invalid input")
	}
	t, ok := strings.CutPrefix(lines[0], "Time:")
	if !ok {
		panic("invalid time line")
	}
	d, ok := strings.CutPrefix(lines[1], "Distance:")
	if !ok {
		panic("invalid distance line")
	}
	for strings.IndexByte(t, ' ') != -1 {
		t = strings.ReplaceAll(t, " ", "")
	}
	for strings.IndexByte(d, ' ') != -1 {
		d = strings.ReplaceAll(d, " ", "")
	}
	ti, terr := strconv.Atoi(t)
	di, derr := strconv.Atoi(d)
	if terr != nil || derr != nil {
		panic("integer conversion failed")
	}
	return ti, di
}

func distance(time, hold int) int {
	return (time - hold) * hold
}

func (in *Input) Run() {
	product := 1
	for i := 0; i < len(in.Time); i++ {
		time := in.Time[i]
		dist := in.Distance[i]
		better := SolveSample(time, dist)
		product *= better
	}
	fmt.Println("Product:", product)
}

func SolveSample(time, dist int) int {
	debug(1, fmt.Sprintf("solving time=%d, distance=%d", time, dist))
	better := func(hold int) bool {
		return distance(time, hold) > dist
	}
	if better(0) != false || better(time) != false {
		panic("weird input")
	}
	index := searchBetter(time, better)
	fb := searchFirstBetter(index, better)
	lb := searchLastBetter(time, index, better)
	return lb - fb + 1
}

func searchBetter(time int, better func(hold int) bool) int {
	// test if the midpoint is better, so that the lower index can be search in the lower part of the
	// time-range and the upper index can be searched for in the upper part to the time-range.
	test := time / 2
	if better(test) {
		return test
	}
	panic("no better value found")
}

func searchFirstBetter(index int, better func(hold int) bool) int {
	i := sort.Search(index, better)
	if i == index {
		panic("search failed")
	}
	return i
}

func searchLastBetter(time, index int, better func(hold int) bool) int {
	// search finds the first index where the function holds true
	// search for the first non-better index
	last := func(hold int) bool {
		return hold > index && !better(hold)
	}
	i := sort.Search(time, last)
	if i == time {
		panic("search failed")
	}
	// found the first non-better index
	// go one back to have the last better index
	return i - 1
}

var DEBUG = 1

func debug(v int, msg string) {
	if v <= DEBUG {
		fmt.Println(msg)
	}
}
