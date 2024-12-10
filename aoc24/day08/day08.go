package main

// https://adventofcode.com/2024/day/08

import (
	"flag"
	"fmt"

	"github.com/phicode/challenges/lib"
	"github.com/phicode/challenges/lib/rowcol"
)

func main() {
	flag.Parse()
	lib.Timed("Part 1", ProcessPart1, "aoc24/day08/example.txt")
	lib.Timed("Part 1", ProcessPart1, "aoc24/day08/input.txt")

	lib.Timed("Part 2", ProcessPart2, "aoc24/day08/example2.txt")
	lib.Timed("Part 2", ProcessPart2, "aoc24/day08/example.txt")
	lib.Timed("Part 2", ProcessPart2, "aoc24/day08/input.txt")

	//lib.Profile(1, "day08.pprof", "Part 2", ProcessPart2, "aoc24/day08/input.txt")
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
	grid        rowcol.Grid[byte]
	frequencies []Frequency
}

type Frequency struct {
	f        byte
	antennas []rowcol.Pos
}

func ParseInput(name string) *Input {
	lines := lib.ReadLines(name)
	grid := rowcol.NewByteGridFromStrings(lines)

	byFreq := make(map[byte][]rowcol.Pos)
	grid.VisitWithPos(func(v byte, p rowcol.Pos) {
		if v == '.' || v == '#' {
			return
		}
		antennas := byFreq[v]
		antennas = append(antennas, p)
		byFreq[v] = antennas
	})
	var freqs []Frequency
	for freq, antennas := range byFreq {
		freqs = append(freqs, Frequency{freq, antennas})
	}

	return &Input{grid: grid, frequencies: freqs}
}

////////////////////////////////////////////////////////////

func SolvePart1(input *Input) int {
	rows, cols := input.grid.Size()
	antinodes := rowcol.NewGrid[byte](rows, cols)
	antinodes.Reset('.')

	for _, freq := range input.frequencies {
		for i, antennaA := range freq.antennas {
			for j := i + 1; j < len(freq.antennas); j++ {
				antennaB := freq.antennas[j]
				p, q := Antinodes(antennaA, antennaB)
				antinodes.SetIfValidPos(p, '#')
				antinodes.SetIfValidPos(q, '#')
			}
		}
		lib.Log(lib.LogDebug, "after frequency", freq.f)
		Print(antinodes)
	}
	return antinodes.Count(func(v byte) bool { return v == '#' })
}

func Antinodes(a, b rowcol.Pos) (rowcol.Pos, rowcol.Pos) {
	d := a.Sub(b)
	p := a.Add(d)
	q := b.Sub(d)
	return p, q
}

func Print(g rowcol.Grid[byte]) {
	for _, row := range g.Data {
		lib.Log(lib.LogDebug, string(row))
	}
}

////////////////////////////////////////////////////////////

func SolvePart2(input *Input) int {
	rows, cols := input.grid.Size()
	antinodes := rowcol.NewGrid[byte](rows, cols)
	antinodes.Reset('.')

	for _, freq := range input.frequencies {
		for i, antennaA := range freq.antennas {
			for j := i + 1; j < len(freq.antennas); j++ {
				antennaB := freq.antennas[j]
				antinodes.SetPos(antennaA, '#')
				antinodes.SetPos(antennaB, '#')
				d := antennaA.Sub(antennaB)

				p := antennaA.Add(d)
				for antinodes.IsValidPos(p) {
					antinodes.SetPos(p, '#')
					p = p.Add(d)
				}

				q := antennaB.Sub(d)
				for antinodes.IsValidPos(q) {
					antinodes.SetPos(q, '#')
					q = q.Sub(d)
				}
			}
		}
		lib.Log(lib.LogDebug, "after frequency", freq.f)
		Print(antinodes)
	}
	return antinodes.Count(func(v byte) bool { return v == '#' })
}
