package main

// https://adventofcode.com/2025/day/5

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/phicode/challenges/lib"
)

func main() {
	flag.Parse()
	lib.Timed("Part 1", ProcessPart1, "aoc25/day05/example.txt")
	lib.Timed("Part 1", ProcessPart1, "aoc25/day05/input.txt")

	lib.Timed("Part 2", ProcessPart2, "aoc25/day05/example.txt")
	lib.Timed("Part 2", ProcessPart2, "aoc25/day05/input.txt")

	//lib.Profile(1, "day05.pprof", "Part 2", ProcessPart2, "aoc25/day05/input.txt")
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

type Range struct {
	a, b int
}

func (r Range) contains(i int) bool {
	return i >= r.a && i <= r.b
}

func (r Range) isZero() bool {
	return r.a == 0 && r.b == 0
}

func (r Range) String() string {
	return fmt.Sprintf("%d-%d", r.a, r.b)
}

type Input struct {
	ranges []Range
	ids    []int
}

func ReadAndParseInput(name string) Input {
	lines := lib.ReadLines(name)
	return ParseInput(lines)
}

func ParseInput(lines []string) Input {
	var rv Input
	separator := false
	for _, line := range lines {
		if line == "" {
			separator = true
			continue
		}
		if separator {
			id, err := strconv.Atoi(line)
			if err != nil {
				panic(err)
			}
			rv.ids = append(rv.ids, id)
		} else {
			rv.ranges = append(rv.ranges, parseRange(line))
		}
	}
	return rv
}

func parseRange(s string) Range {
	var r Range
	n, err := fmt.Sscanf(s, "%d-%d", &r.a, &r.b)
	if n != 2 || err != nil {
		panic(fmt.Errorf("failed to parse range: %v", err))
	}
	if r.a > r.b {
		panic(fmt.Errorf("invalid range %d - %d", r.a, r.b))
	}
	return Range{
		a: r.a,
		b: r.b,
	}
}

func (i Input) inAnyRange(v int) bool {
	for _, r := range i.ranges {
		if r.contains(v) {
			return true
		}
	}
	return false
}

////////////////////////////////////////////////////////////

func SolvePart1(input Input) int {
	fresh := 0
	for _, id := range input.ids {
		if input.inAnyRange(id) {
			fresh++
		}
	}
	return fresh
}

////////////////////////////////////////////////////////////

func SolvePart2(input Input) int {
	rs := mergeAll(input.ranges)
	validateNoOverlap(rs)
	total := 0
	for _, r := range rs {
		total += r.b - r.a + 1
	}
	return total
}

func validateNoOverlap(rs []Range) {
	for i, a := range rs {
		for j, b := range rs[i+1:] {
			if overlap(a, b) {
				panic(fmt.Errorf("overlap found: (%d, %d) = %v, %v", i, i+j+1, a, b))
			}
		}
	}
}

func overlap(a, b Range) bool {
	if a.b < b.a || a.a > b.b {
		return false
	}
	return true
}

func mergeAll(rs []Range) []Range {
	var rv []Range
	for i, r := range rs {
		if r.isZero() { // already merged
			continue
		}

		for j := i + 1; j < len(rs); j++ {
			if rs[j].isZero() {
				continue
			}
			if merged, m := merge(r, rs[j]); merged {
				//fmt.Printf("merged %v + %v = %v\n", r, rs[j], m)
				r = m
				rs[j] = Range{} // zero out merged entry
				// if we merge any entry, we have to search for new candidates from the beginning
				j = i
			}
		}
		rv = append(rv, r)
	}
	return rv
}

func merge(a, b Range) (bool, Range) {
	if a.b < b.a || a.a > b.b {
		return false, Range{}
	}
	return true, Range{
		a: min(a.a, b.a),
		b: max(a.b, b.b),
	}
}
