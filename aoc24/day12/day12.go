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
	//lib.Timed("Part 1", ProcessPart1, "aoc24/day12/input.txt")
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
	for pos := range s.visited.PosIterator() {
		if s.visited.GetPos(pos) {
			continue
		}
		plant, ps := s.follow(pos)
		fmt.Println("Plant:", rune(plant), "area", len(ps))
	}

	return 0
}

////////////////////////////////////////////////////////////

func SolvePart2(input Input) int {
	return 0
}
