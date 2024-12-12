package main

// https://adventofcode.com/2024/day/12

import (
	"flag"
	"fmt"

	"github.com/phicode/challenges/lib"
	"github.com/phicode/challenges/lib/rowcol"
)

func main() {
	flag.Parse()
	lib.Timed("Part 1", ProcessPart1, "aoc24/day12/example.txt")
	lib.Timed("Part 1", ProcessPart1, "aoc24/day12/input.txt")
	//
	//lib.Timed("Part 2", ProcessPart2, "aoc24/day12/example.txt")
	//lib.Timed("Part 2", ProcessPart2, "aoc24/day12/input.txt")

	//lib.Profile(1, "day12.pprof", "Part 2", ProcessPart2, "aoc24/dayXX/input.txt")
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
	grid rowcol.Grid[byte]
}

func ReadAndParseInput(name string) Input {
	lines := lib.ReadLines(name)
	return ParseInput(lines)
}

func ParseInput(lines []string) Input {
	return Input{rowcol.NewByteGridFromStrings(lines)}
}

////////////////////////////////////////////////////////////

type State struct {
	grid    rowcol.Grid[byte]
	visited rowcol.Grid[bool]
}

func (s *State) follow(p rowcol.Pos) (byte, []rowcol.Pos) {
	plant := s.grid.GetPos(p)
	s.visited.SetPos(p, true)
	var ps []rowcol.Pos
	ps = append(ps, p)
	return plant, s.followPlantDirs(plant, ps, p)
}

func (s *State) followPlant(plant byte, ps []rowcol.Pos, p rowcol.Pos) []rowcol.Pos {
	if !s.grid.IsValidPos(p) || s.grid.GetPos(p) != plant {
		return ps
	}
	if s.visited.GetPos(p) {
		return ps
	}
	s.visited.SetPos(p, true)
	ps = append(ps, p)
	return s.followPlantDirs(plant, ps, p)
}

func (s *State) followPlantDirs(plant byte, ps []rowcol.Pos, p rowcol.Pos) []rowcol.Pos {
	ps = s.followPlant(plant, ps, p.AddDir(rowcol.Up))
	ps = s.followPlant(plant, ps, p.AddDir(rowcol.Down))
	ps = s.followPlant(plant, ps, p.AddDir(rowcol.Left))
	ps = s.followPlant(plant, ps, p.AddDir(rowcol.Right))
	return ps
}

func SolvePart1(input Input) int {
	s := State{grid: input.grid, visited: rowcol.NewGrid[bool](input.grid.Size())}
	total := 0
	for pos := range s.visited.PosIterator() {
		if s.visited.GetPos(pos) {
			continue
		}
		plant, ps := s.follow(pos)
		peri := perimeter(ps, s.grid, plant)
		//fmt.Printf("Plant: %c ; area: %d, perimeter: %d\n", plant, len(ps), peri)
		total += peri * len(ps)
	}
	return total
}

func perimeter(ps []rowcol.Pos, grid rowcol.Grid[byte], plant byte) int {
	peri := 0
	for _, p := range ps {
		peri += perimeterContribution(p, grid, plant)
	}
	return peri
}

//func perimeter(ps []rowcol.Pos, grid rowcol.Grid[byte], plant byte) int {
//	_min := rowcol.MinPos(ps)
//	peri := perimeterContribution(_min, grid, plant)
//	//assert.True(peri >= 1 && peri < 4)
//	cur, notDir := next(rowcol.Up, _min, grid, plant)
//	for cur != _min {
//		contr := perimeterContribution(cur, grid, plant)
//		//assert.True(contr > 0 && contr < 4)
//		peri += contr
//		cur, notDir = next(notDir, cur, grid, plant)
//	}
//	return peri
//}

//var nextDirections = []rowcol.Direction{rowcol.Right, rowcol.Down, rowcol.Left, rowcol.Up}

//	func next(notDir rowcol.Direction, pos rowcol.Pos, grid rowcol.Grid[byte], plant byte) (rowcol.Pos, rowcol.Direction) {
//		for _, dir := range nextDirections {
//			test := pos.AddDir(dir)
//			if dir != notDir && grid.IsValidPos(test) && grid.GetPos(test) == plant {
//				return test, dir.Reverse()
//			}
//		}
//		panic("no position to follow found")
//	}
func perimeterContribution(pos rowcol.Pos, grid rowcol.Grid[byte], plant byte) int {
	sameNeighbors := 0
	for _, dir := range rowcol.Directions {
		test := pos.AddDir(dir)
		if grid.IsValidPos(test) && grid.GetPos(test) == plant {
			sameNeighbors++
		}
	}
	// 4 same neighbors => plant is fully surounded => no perimeter contribution
	// 3 plant has 1 neighbors on 3 sided
	return 4 - sameNeighbors
}

////////////////////////////////////////////////////////////

func SolvePart2(input Input) int {
	return 0
}
