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

	//drawPar2Expand()

	//lib.LogLevel = lib.LogDebug
	lib.Timed("Part 2", ProcessPart2, "aoc24/day20/example.txt")
	//lib.LogLevel = lib.LogInfo
	lib.Timed("Part 2", ProcessPart2, "aoc24/day20/input.txt")

	//lib.Profile(1, "day20.pprof", "Part 2", ProcessPart2, "aoc24/day20/input.txt")
}

func drawPar2Expand() {
	grid := rowcol.NewGrid[byte](50, 50)
	grid.Reset(' ')
	start := rowcol.PosXY(25, 25)
	grid.SetPos(start, 'S')
	for _, dir := range rowcol.Directions {
		for cheatLen := 2; cheatLen <= 20; cheatLen++ {
			for p2len := 0; p2len < cheatLen; p2len++ {
				p1len := cheatLen - p2len
				dir2 := dir.Right()
				end := start.Add(dir.MulS(p1len)).Add(dir2.MulS(p2len))
				assert.True(grid.GetPos(end) == ' ')
				grid.SetPos(end, '#')
			}
		}
	}
	rowcol.PrintByteGrid(grid)
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
	distances.SetPos(start, 0)
	prev, cur := start, start
	next := findNext(grid, prev, cur)
	dist := 1
	for next != end {
		distances.SetPos(next, dist)
		dist++
		prev, cur = cur, next
		next = findNext(grid, prev, cur)
	}
	distances.SetPos(end, dist)
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
	end := pos.AddDir(dir).AddDir(dir)
	if !in.grid.IsValidPos(end) {
		return
	}
	if in.grid.GetPos(end) != '.' {
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
	for pos, distance := range input.distances.Iterator() {
		if distance == -1 {
			continue
		}
		input.findCheatsPart2(pos)
	}
	entries := lib.MapEntries(input.savings)
	slices.SortFunc(entries, func(a, b lib.Entry[int, int]) int {
		return a.Key - b.Key
	})
	over100 := 0
	for _, entry := range entries {
		if entry.Key >= 50 && lib.LogLevel >= lib.LogDebug {
			fmt.Printf("%d cheats with saving of %d\n", entry.Value, entry.Key)
		}
		if entry.Key >= 100 {
			over100 += entry.Value
		}
	}

	return over100
}

func (in Input) findCheatsPart2(start rowcol.Pos) {
	// 4 quadrants
	// left+up
	// up+right
	// right+down
	// down+left
	for _, dir := range rowcol.Directions {
		for cheatLen := 2; cheatLen <= 20; cheatLen++ {
			for p2len := 0; p2len < cheatLen; p2len++ {
				p1len := cheatLen - p2len
				dir2 := dir.Right()
				end := start.Add(dir.MulS(p1len)).Add(dir2.MulS(p2len))
				in.checkPart2Cheat(start, end, cheatLen)
			}
		}
	}
}

func (in Input) checkPart2Cheat(start, end rowcol.Pos, cheatLen int) {
	if !in.grid.IsValidPos(end) {
		return
	}
	if in.grid.GetPos(end) != '.' && in.grid.GetPos(end) != 'E' {
		return
	}
	startDist := in.distances.GetPos(start)
	endDist := in.distances.GetPos(end)
	assert.True(startDist >= 0)
	assert.True(endDist >= 0)
	saving := endDist - startDist - cheatLen
	if saving <= 0 {
		return
	}
	in.savings[saving]++
}
