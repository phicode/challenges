package main

// https://adventofcode.com/2022/day/18

import (
	"fmt"

	"git.bind.ch/phil/challenges/lib"
)

var VERBOSE = 1

func main() {
	// 64
	ProcessPart1("aoc22/day18/example.txt")
	ProcessPart1("aoc22/day18/input.txt")

	ProcessPart2("aoc22/day18/example.txt")
	ProcessPart2("aoc22/day18/input.txt")
}

func ProcessPart1(name string) {
	fmt.Println("Part 1 input:", name)
	lines := lib.ReadLines(name)
	ps := ParsePositions(lines)
	area := FindSurfaceArea(ps)
	fmt.Println("Area:", area)

	fmt.Println()
}

func ProcessPart2(name string) {
	fmt.Println("Part 2 input:", name)
	lines := lib.ReadLines(name)
	_ = lines

	fmt.Println()
}

func log(v int, msg string) {
	if v <= VERBOSE {
		fmt.Println(msg)
	}
}

////////////////////////////////////////////////////////////

type Position struct {
	X, Y, Z int
}

func ParsePositions(lines []string) []Position {
	var rv []Position
	for _, l := range lines {
		var p Position
		n, err := fmt.Sscanf(l, "%d,%d,%d", &p.X, &p.Y, &p.Z)
		if n != 3 || err != nil {
			panic(fmt.Errorf("invalid input; n=%d, err=%w", n, err))
		}
		rv = append(rv, p)
	}
	return rv
}

func FindSurfaceArea(ps []Position) int {
	m := make(map[Position]struct{})
	for _, p := range ps {
		m[p] = struct{}{}
	}
	var area int
	for _, p := range ps {
		area += Exposed(m, p, 1, 0, 0)
		area += Exposed(m, p, -1, 0, 0)
		area += Exposed(m, p, 0, 1, 0)
		area += Exposed(m, p, 0, -1, 0)
		area += Exposed(m, p, 0, 0, 1)
		area += Exposed(m, p, 0, 0, -1)
	}
	return area
}

func Exposed(m map[Position]struct{}, p Position, addX, addY, addZ int) int {
	test := Position{p.X + addX, p.Y + addY, p.Z + addZ}
	if _, found := m[test]; found {
		return 0
	}
	return 1
}

func Render(ps []Position, a, b func(position Position) int) string {
	//as := lib.Map(ps, a)
	//bs := lib.Map(ps, b)
	//amax := max(as[0], as[1:]...) // columns
	//bmax := max(bs[0], bs[1:]...) // lines
	return ""
}
func GetX(p Position) int { return p.X }
func GetY(p Position) int { return p.Y }
func GetZ(p Position) int { return p.Z }
