package main

// https://adventofcode.com/2024/day/20

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
	lib.Timed("Part 1", ProcessPart1, "aoc24/day20/example.txt")
	lib.Timed("Part 1", ProcessPart1, "aoc24/day20/input.txt")
	//
	//lib.Timed("Part 2", ProcessPart2, "aoc24/day20/example.txt")
	//lib.Timed("Part 2", ProcessPart2, "aoc24/day20/input.txt")

	//lib.Profile(1, "day20.pprof", "Part 2", ProcessPart2, "aoc24/day20/input.txt")
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
	grid      rowcol.Grid[byte]
	distances rowcol.Grid[int]

	// key: distance saved
	// value: number of cheats found
	savings map[int]int
}

func ReadAndParseInput(name string) Input {
	lines := lib.ReadLines(name)
	return ParseInput(lines)
}

func ParseInput(lines []string) Input {
	grid := rowcol.NewByteGridFromStrings(lines)
	distances := createDistanceMap(grid)
	return Input{
		grid:      grid,
		distances: distances,
		savings:   make(map[int]int),
	}
}

func createDistanceMap(grid rowcol.Grid[byte]) rowcol.Grid[int] {
	start := grid.MustFindFirst(func(v byte) bool { return v == 'S' })
	end := grid.MustFindFirst(func(v byte) bool { return v == 'E' })
	distances := rowcol.NewGrid[int](grid.Size())
	distances.Reset(-1)
	prev, cur := start, start
	next := findNext(grid, prev, cur)
	dist := 1
	for next != end {
		distances.SetPos(next, dist)
		dist++
		prev, cur = cur, next
		next = findNext(grid, prev, cur)
	}
	return distances
}

func findNext(grid rowcol.Grid[byte], prev, current rowcol.Pos) rowcol.Pos {
	for _, dir := range rowcol.Directions {
		c := current.AddDir(dir)
		if grid.GetPos(c) != '#' && c != prev {
			return c
		}
	}
	panic("no next position found")
}

////////////////////////////////////////////////////////////

func SolvePart1(input Input) int {
	for pos, distance := range input.distances.Iterator() {
		if distance == -1 {
			continue
		}
		input.findCheats(pos)
	}
	entries := lib.MapEntries(input.savings)
	slices.SortFunc(entries, func(a, b lib.Entry[int, int]) int {
		return a.Key - b.Key
	})
	over100 := 0
	for _, entry := range entries {
		//fmt.Printf("%d cheats with saving of %d\n", entry.Value, entry.Key)
		if entry.Key >= 100 {
			over100 += entry.Value
		}
	}

	return over100
}

func (in Input) findCheats(pos rowcol.Pos) {
	for _, dir := range rowcol.Directions {
		in.findCheatsDir(pos, dir)
	}
}

func (in Input) findCheatsDir(pos rowcol.Pos, dir rowcol.Direction) {
	wall := pos.AddDir(dir)
	end := wall.AddDir(dir)
	d := in.distances
	if !d.IsValidPos(end) {
		return
	}
	if in.grid.GetPos(wall) != '#' {
		return
	}
	startDist := in.distances.GetPos(pos)
	endDist := in.distances.GetPos(end)
	assert.True(startDist >= 0)
	if endDist == -1 {
		return
	}
	saving := endDist - startDist - 2
	if saving < 0 {
		return
	}
	in.savings[saving]++
}

////////////////////////////////////////////////////////////

func SolvePart2(input Input) int {
	return 0
}
