package main

// https://adventofcode.com/2022/day/17

import (
	"bufio"
	"fmt"
	"os"
	"slices"
)

func main() {
	Process("aoc22/dayXX/example.txt")
	Process("aoc22/dayXX/input.txt")
}

func Process(name string) {
	fmt.Println("input:", name)
	lines := ReadInput(name)
	_ = lines

	fmt.Println()
}

var (
	Rocks = []Rock{
		RockDash,
		RockPlus,
		RoclInvL,
		RockPipe,
		RockBlock,
	}
	RockDash = NewRock(
		LineS("####"),
	)
	RockPlus = NewRock(
		LineS(".#."),
		LineS("###"),
		LineS(".#."),
	)
	RoclInvL = NewRock(
		LineS("..#"),
		LineS("..#"),
		LineS("###"),
	)
	RockPipe = NewRock(
		LineS("#"),
		LineS("#"),
		LineS("#"),
		LineS("#"),
	)
	RockBlock = NewRock(
		LineS("##"),
		LineS("##"),
	)
)

func NewRock(ls ...Line) Rock {
	r := Rock(ls)
	slices.Reverse(r)
	return r
}

func ReadInput(name string) []string {
	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	s := bufio.NewScanner(f)

	var lines []string
	for s.Scan() {
		line := s.Text()
		lines = append(lines, line)
	}

	if err := s.Err(); err != nil {
		panic(s.Err())
	}
	return lines
}

////////////////////////////////////////////////////////////

type Cave struct {
	blub [][]byte
}

// stores a bitfield representing a line
// ...##.
// 000110
type Line int

func LineS(s string) Line {
	return LineB([]byte(s))
}
func LineB(content []byte) Line {
	var l Line
	for _, v := range content {
		l = l << 1
		if v == '#' {
			l |= 1
		}
	}
	return l
}

func (l Line) Intersects(b Line) bool {
	return l&b != 0
}

func (l Line) Left() Line {
	return l << 1
}
func (l Line) Right() Line {
	return l >> 1
}
func (l Line) CanMoveLeft() {

}

type Rock []Line

func (r Rock) Copy() Rock {
	rv := make([]Line, len(r))
	for i, l := range r {
		rv[i] = l
	}
	return rv
}
func (r Rock) Intersects(lines []Line) bool {
	l := len(r)
	ll := len(lines)
	for i := 0; i < l; i++ {
		// i=0 -> topmost row of rock
		idx := ll - i - 1
		if idx < 0 {
			return true
		}
		if r[i].Intersects(lines[idx]) {
			return true
		}
	}
	return false
}
