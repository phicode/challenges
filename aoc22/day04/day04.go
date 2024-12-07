package main

// https://adventofcode.com/2022/day/4

import (
	"bufio"
	"fmt"
	"os"

	"github.com/phicode/challenges/lib/assets"
)

func main() {
	Process("aoc22/day04/example.txt")
	Process("aoc22/day04/input.txt")
}

func Process(name string) {
	fmt.Println("input:", name)
	pairs := ReadInput(name)
	fc := 0
	ol := 0
	for _, p := range pairs {
		if p.HasFullyContained() {
			fc++
		}
		if p.Overlaps() {
			ol++
		}
	}
	fmt.Println("fully contained", fc)
	fmt.Println("overlaps", ol)

	fmt.Println()
}

func ReadInput(name string) []Pair {
	f, err := os.Open(assets.MustFind(name))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	s := bufio.NewScanner(f)

	var pairs []Pair
	for s.Scan() {
		var x, y Range
		n, err := fmt.Sscanf(s.Text(), "%d-%d,%d-%d", &x.A, &x.B, &y.A, &y.B)
		if err != nil || n != 4 {
			panic(fmt.Errorf("invalid input: %q", s.Text()))
		}
		pairs = append(pairs, Pair{x, y})
	}

	if err := s.Err(); err != nil {
		panic(s.Err())
	}
	return pairs
}

type Pair struct {
	X Range
	Y Range
}

func (p Pair) HasFullyContained() bool {
	return p.X.FullyContained(p.Y) || p.Y.FullyContained(p.X)
}
func (p Pair) Overlaps() bool {
	return p.X.Overlaps(p.Y)
}

type Range struct {
	A int
	B int
}

func (x Range) FullyContained(y Range) bool {
	return x.A <= y.A && x.B >= y.B
}

func (x Range) Overlaps(y Range) bool {
	noOverlap := x.B < y.A || x.A > y.B
	return !noOverlap
}
