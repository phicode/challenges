package main

// https://adventofcode.com/2022/day/17

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"slices"
)

var debug = 1

const reduce = false
const findRepetitions = true

func main() {
	fmt.Println("Example debug")
	Process("aoc22/day17/example.txt", 11)
	debug = 0
	fmt.Println("Example")
	Process("aoc22/day17/example.txt", 2022)

	fmt.Println("Step 1")
	Process("aoc22/day17/input.txt", 2022)

	fmt.Println("Step 2")
	Process("aoc22/day17/input.txt", 1000000000000)
}

func Process(name string, rocksToPlace int) {
	fmt.Println("input:", name)
	lines := ReadInput(name)

	c := Run(lines, rocksToPlace)
	fmt.Println("height:", c.Height())
	fmt.Println("removed:", c.Removed)
	fmt.Println("total:", c.Height()+c.Removed+c.HeightSkipped)

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
	Full = LineS("#######")
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
	RocksAt       []int // rock index was placed at this line index
	NextRockIndex int
	ActiveRock    Rock
	RockOffset    int // how far of the floor the rock is
	RocksToPlace  int
	RocksPlaced   int
	Removed       int
	HeightSkipped int
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
		//fmt.Println("padding", pad, "lines")
		for i := 0; i < pad; i++ {
			c.Grow()
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

func (l Line) IsFull() bool {
	return l == Full
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

func Run(lines []string, rocksToPlace int) *Cave {
	if len(lines) > 1 {
		panic("too many lines")
	}
	input := []rune(lines[0])
	c := &Cave{
		RocksToPlace: rocksToPlace,
	}

	for i := 0; !c.Finished(); i++ {
		moveIdx := i % len(input)
		move := input[moveIdx]
		if len(c.ActiveRock) == 0 {
			c.Fill()
			c.SpawnRock()
		}
		title := "after move right"
		if move == '<' {
			title = "after move left"
		}
		c.ApplyMove(move)
		c.Print(title, 2)

		c.ApplyGravity()
		c.Print("after gravity", 2)

		// space saving stuff
		//if reduce {
		//	if i > 0 && i%100_000 == 0 {
		//		c.Reduce()
		//		if i%1_000_000 == 0 {
		//			fmt.Printf("i=%d rocks-placed: %d\n", i, c.RocksPlaced)
		//		}
		//	}
		//}
		if findRepetitions {
			if c.RocksPlaced == 10_000 {
				c.FindRepetitions()
			}
		}
	}
	return c
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
	if c.RockOffset == -1 {
		c.RockOffset++
		c.FixRock()
	}
	if c.Intersects() {
		c.RockOffset++
		c.FixRock()
	}
}

func (c *Cave) Print(title string, level int) {
	if level > debug {
		return
	}
	fmt.Println(title)
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
			c.Grow()
		}
		c.Lines[idx] = c.Lines[idx].Merge(c.ActiveRock[i])
		c.RocksAt[idx] = c.NextRockIndex - 1
	}
	c.ActiveRock = nil
	c.RocksPlaced++
}

func (c *Cave) Grow() {
	c.Lines = append(c.Lines, Line(0))
	c.RocksAt = append(c.RocksAt, 0)
}

func (c *Cave) Finished() bool {
	return c.RocksPlaced >= c.RocksToPlace
}

func (c *Cave) SpawnRock() {
	c.ActiveRock = Rocks[c.NextRockIndex%len(Rocks)].Copy()
	// move to 2 spaces away from left wall
	c.ActiveRock = c.ActiveRock.MoveRight()
	c.ActiveRock = c.ActiveRock.MoveRight()
	c.NextRockIndex++
	c.RockOffset = len(c.Lines)
	c.Print("spawn rock", 1)
}

func (c *Cave) Height() int {
	// index of the highest line which is not empty
	h := len(c.Lines) - 1
	for h >= 0 {
		if !c.Lines[h].IsEmpty() {
			break
		}
		h--
	}
	return h + 1
}

func (c *Cave) Reduce() {
	full := c.FindFull()
	if full < 0 {
		return
	}
	c.Removed += full + 1
	copy(c.Lines[:], c.Lines[full+1:])
	kept := len(c.Lines) - full - 1
	c.Lines = c.Lines[:kept]
}

func (c *Cave) FindFull() int {
	for i := len(c.Lines) - 1; i >= 0; i-- {
		if c.Lines[i].IsFull() {
			return i
		}
	}
	return -1
}

func (c *Cave) FindRepetitions() {
	fmt.Println("LOOP FINDER")
	var fullIndxs []int
	for i, v := range c.Lines {
		if v.IsFull() {
			fullIndxs = append(fullIndxs, i)
		}
	}
	for i, idx := range fullIndxs {
		// find 2 loops
		for j := i + 1; j < len(fullIndxs); j++ {
			start := idx
			end := fullIndxs[j]

			start2 := end
			end2 := end + (end - start)
			if c.IsLoop(start, end) && c.IsLoop(start2, end2) {
				rocksDiff1 := c.RocksAt[end] - c.RocksAt[start]
				rocksDiff2 := c.RocksAt[end2] - c.RocksAt[start2]
				//fmt.Printf("loop found, start=%d, end=%d, length=%d\n", start, end, end-start)
				fmt.Printf("loop found, start=%d, end=%d, length=%d, d1=%d, d2=%d\n", start, end, end-start, rocksDiff1, rocksDiff2)
				n := end - start
				rocksRemaining := c.RocksToPlace - c.RocksPlaced
				skips := rocksRemaining / rocksDiff1
				plusRocks := skips * rocksDiff1
				plusHeight := skips * n
				fmt.Printf("can skip %d loops for %d heigh, %d rocks\n", skips, plusHeight, plusRocks)
				c.RocksPlaced += plusRocks
				c.HeightSkipped += plusHeight
				return
			}
		}
	}
}

// start inclusive, end exclusive
func (c *Cave) IsLoop(start, end int) bool {
	l := end - start
	if end+l >= len(c.Lines) {
		return false
	}
	for i := 0; i < l; i++ {
		if c.Lines[start+i] != c.Lines[end+i] {
			return false
		}
	}
	return true
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
