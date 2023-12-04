package main

// https://adventofcode.com/2022/day/17

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"slices"
)

func main() {
	Process("aoc22/day17/example.txt")
	//Process("aoc22/day17/input.txt")
}

func Process(name string) {
	fmt.Println("input:", name)
	lines := ReadInput(name)

	Run(lines)

	fmt.Println()
}

var (
	Rocks = []Rock{
		RockDash,
		RockPlus,
		RockInvL,
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
	RockInvL = NewRock(
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
	for r.CanMoveLeft() {
		r = r.MoveLeft()
	}
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
	Lines         []Line
	NextRockIndex int
	ActiveRock    Rock
	RockOffset    int // index
}

func (c *Cave) Render(w *bytes.Buffer) {
	Render(w, c.Lines)
}

// ensures 3 empty lines are on top of
func (c *Cave) Fill() {
	pad := 3
	for i := len(c.Lines) - 1; i >= 0 && pad > 0; i-- {
		if c.Lines[i].IsEmpty() {
			pad--
		} else {
			break
		}
	}
	if pad > 0 {
		fmt.Println("padding", pad, "lines")
		for i := 0; i < pad; i++ {
			c.Lines = append(c.Lines, Line(0))
		}
	}
}

////////////////////////////////////////////////////////////

// stores a bitfield representing a line
// ...##..
// 0001100
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

func (l Line) IsEmpty() bool {
	return l == 0
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
func (l Line) CanMoveLeft() bool {
	// position 7 is free
	return l&0b0100_0000 == 0
}
func (l Line) CanMoveRight() bool {
	return l&1 == 0
}

func (l Line) String() string {
	var b bytes.Buffer
	l.Render(&b)
	return b.String()
}

func (l Line) Render(w *bytes.Buffer) {
	w.WriteByte('|')
	for i := 6; i >= 0; i-- {
		mask := 1 << i
		if (int(l) & mask) != 0 {
			w.WriteByte('#')
		} else {
			w.WriteByte('.')
		}
	}
	w.WriteString("|\n")
}

func (l Line) Merge(b Line) Line {
	return l | b
}

type Rock []Line

func (r Rock) MoveLeft() Rock {
	b := make([]Line, len(r))
	for i, l := range r {
		b[i] = l.Left()
	}
	return b
}

func (r Rock) MoveRight() Rock {
	b := make([]Line, len(r))
	for i, l := range r {
		b[i] = l.Right()
	}
	return b
}

func (r Rock) CanMoveLeft() bool {
	for _, l := range r {
		if !l.CanMoveLeft() {
			return false
		}
	}
	return true
}
func (r Rock) CanMoveRight() bool {
	for _, l := range r {
		if !l.CanMoveRight() {
			return false
		}
	}
	return true
}

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

func (r Rock) Render(w *bytes.Buffer) {
	for i := len(r) - 1; i >= 0; i-- {
		r[i].Render(w)
	}
}

////////////////////////////////////////////////////////////

func Run(lines []string) {
	if len(lines) > 1 {
		panic("too many lines")
	}
	input := []rune(lines[0])
	c := &Cave{}

	for i := 0; ; i++ {
		move := input[i%len(input)]
		if len(c.ActiveRock) == 0 {
			c.Fill()
			fmt.Println("spawning rock")
			c.ActiveRock = Rocks[c.NextRockIndex%len(Rocks)].Copy()
			// move to to the right
			c.ActiveRock = c.ActiveRock.MoveRight()
			c.ActiveRock = c.ActiveRock.MoveRight()
			c.NextRockIndex++
			c.RockOffset = len(c.Lines)
			c.Print()
		}
		if move == '>' {
			fmt.Println("wind right")
		} else {
			fmt.Println("wind left")
		}
		c.ApplyMove(move)
		c.Print()

		fmt.Println("gravity")
		c.ApplyGravity()
		c.Print()

	}
}

func Render(w *bytes.Buffer, m []Line) {
	for i := len(m) - 1; i >= 0; i-- {
		m[i].Render(w)
	}
	w.WriteString("+-------+\n")
}

func (c *Cave) ApplyMove(move rune) {
	switch move {
	case '>':
		if c.ActiveRock.CanMoveRight() {
			before := c.ActiveRock
			c.ActiveRock = c.ActiveRock.MoveRight()
			if c.Intersects() {
				c.ActiveRock = before
			}
		}
	case '<':
		if c.ActiveRock.CanMoveLeft() {
			before := c.ActiveRock
			c.ActiveRock = c.ActiveRock.MoveLeft()
			if c.Intersects() {
				c.ActiveRock = before
			}
		}
	default:
		panic("invalid move")
	}
}

func (c *Cave) ApplyGravity() {
	c.RockOffset--
	if c.RockOffset == 0 {
		c.FixRock()
		c.ActiveRock = nil
	}
	if c.Intersects() {
		c.RockOffset++
		c.FixRock()
		c.ActiveRock = nil
	}
}

func (c *Cave) Print() {
	var b bytes.Buffer
	m := Merge(c.Lines, c.ActiveRock, c.RockOffset)
	Render(&b, m)
	b.WriteTo(os.Stderr)
	fmt.Println()
}

func (c *Cave) Intersects() bool {
	for i, x := range c.ActiveRock {
		idx := i + c.RockOffset
		if idx >= len(c.Lines) {
			return false
		}
		if c.Lines[idx].Intersects(x) {
			return true
		}
	}
	return false
}

func (c *Cave) FixRock() {
	for i := 0; i < len(c.ActiveRock); i++ {
		idx := i + c.RockOffset
		if len(c.Lines) <= idx {
			c.Lines = append(c.Lines, Line(0))
		}
		c.Lines[idx] = c.Lines[idx].Merge(c.ActiveRock[i])
	}
}

func Merge(a, b []Line, offset int) []Line {
	al := len(a)
	bl := len(b)
	m := make([]Line, max(al, bl+offset))
	for i, x := range a {
		m[i] = x
	}
	for i, x := range b {
		m[i+offset] = m[i+offset].Merge(x)
	}
	return m
}
