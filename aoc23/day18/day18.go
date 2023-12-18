package main

// https://adventofcode.com/2023/day/18

import (
	"fmt"
	"strings"

	"git.bind.ch/phil/challenges/lib"
	"git.bind.ch/phil/challenges/lib/rowcol"
)

// TODO: timing boilerplate
var VERBOSE = 1

func main() {
	ProcessPart1("aoc23/day18/example.txt") // 62
	ProcessPart1("aoc23/day18/input.txt")   // 67891

	//ProcessPart2("aoc23/day18/example.txt")
	//ProcessPart2("aoc23/day18/input.txt")
}

func ProcessPart1(name string) {
	fmt.Println("Part 1 input:", name)
	lines := lib.ReadLines(name)
	instrs := ParseInstructions(lines)
	for _, instr := range instrs {
		fmt.Println(instr)
	}

	lagoon := CreateLagoon(instrs)
	lagoon.DigPath(instrs)
	lagoon.Print()

	lagoon.DigInterior()
	lagoon.Print()

	volume := lagoon.CountDugOut()
	fmt.Println("Volume:", volume)

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

func ParseInstructions(lines []string) []Instruction {
	var rv []Instruction
	for _, l := range lines {
		rv = append(rv, ParseInstruction(l))
	}
	return rv
}

func ParseInstruction(l string) Instruction {
	var dir rune
	var dist int
	var color int
	if n, err := fmt.Sscanf(l, "%c %d (#%x)", &dir, &dist, &color); n != 3 || err != nil {
		panic(fmt.Errorf("invalid instruction: n=%d, err=%w", n, err))
	}
	d, found := Directions[dir]
	if !found {
		panic(fmt.Errorf("invalid direction: %c", dir))
	}
	return Instruction{
		Direction: d,
		Distance:  dist,
		Color:     color,
	}
}

var Directions = map[rune]rowcol.Direction{
	'L': rowcol.Left,
	'R': rowcol.Right,
	'U': rowcol.Up,
	'D': rowcol.Down,
}

type Instruction struct {
	Direction rowcol.Direction
	Distance  int
	Color     int
}

func (i Instruction) String() string {
	return fmt.Sprintf("%s %d (#%x)", i.Direction, i.Distance, i.Color)
}

type Lagoon struct {
	data      rowcol.Grid[Field]
	gridShift rowcol.Pos
}

func (l *Lagoon) ToGridCoordinates(r, c int) (int, int) {
	return r - l.gridShift.Row, c - l.gridShift.Col
}
func (l *Lagoon) ToGridCoordinatesP(p rowcol.Pos) rowcol.Pos {
	return rowcol.Pos{
		Row: p.Row - l.gridShift.Row,
		Col: p.Col - l.gridShift.Col,
	}
}
func (l *Lagoon) ToWorldCoordinatesP(r, c int) rowcol.Pos {
	return rowcol.Pos{Row: r + l.gridShift.Row, Col: c + l.gridShift.Col}
}

func (l *Lagoon) Get(r, c int) Field {
	r, c = l.ToGridCoordinates(r, c)
	return l.data.Get(r, c)
}
func (l *Lagoon) Set(r, c int, v Field) {
	r, c = l.ToGridCoordinates(r, c)
	l.data.Set(r, c, v)
}
func (l *Lagoon) IsValidPosition(r, c int) bool {
	r, c = l.ToGridCoordinates(r, c)
	return l.data.IsValidPosition(r, c)
}
func (l *Lagoon) Rows() int    { return l.data.Rows() }
func (l *Lagoon) Columns() int { return l.data.Columns() }

func (l *Lagoon) MarkDug(x rowcol.Pos, dir rowcol.Direction, color int) {
	f := l.Get(x.Row, x.Col)
	f.DugOut = true
	f.DigDirection = dir
	f.Color = color
	l.Set(x.Row, x.Col, f)
}

type Field struct {
	DigDirection rowcol.Direction
	Color        int
	DugOut       bool
}

const (
	ColorStart    = 0xFFFFFF
	ColorInterior = 0xFFFFFE
)

func CreateLagoon(instrs []Instruction) *Lagoon {
	posMin, posMax, size := LagoonSize(instrs)
	fmt.Println("lagoon size:", size)
	fmt.Println("lagoon minimum:", posMin)
	fmt.Println("lagoon maximum:", posMax)

	return &Lagoon{
		data:      rowcol.NewGrid[Field](size.Row, size.Col),
		gridShift: posMin,
	}
}

func LagoonSize(instrs []Instruction) (rowcol.Pos, rowcol.Pos, rowcol.Pos) {
	current := rowcol.Pos{}
	posMin := rowcol.Pos{}
	posMax := rowcol.Pos{}
	for _, instr := range instrs {
		current = current.Add(instr.Direction.MulS(instr.Distance))
		posMin.Row = min(posMin.Row, current.Row)
		posMin.Col = min(posMin.Col, current.Col)
		posMax.Row = max(posMax.Row, current.Row)
		posMax.Col = max(posMax.Col, current.Col)
	}
	size := posMax.Sub(posMin).Add(rowcol.Pos{Row: 1, Col: 1})
	return posMin, posMax, size
}

func (l *Lagoon) DigPath(instrs []Instruction) {
	current := rowcol.Pos{}
	l.MarkDug(current, rowcol.Stand, ColorStart)
	for _, instr := range instrs {
		for i := 0; i < instr.Distance; i++ {
			current = current.AddDir(instr.Direction)
			l.MarkDug(current, instr.Direction, instr.Color)
		}
	}
}

func (l *Lagoon) DigInterior() {
	rows, cols := l.Rows(), l.Columns()
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			pos := l.ToWorldCoordinatesP(r, c)
			l.DigInteriorAt(pos)
		}
	}
}

func (l *Lagoon) DigInteriorAt(pos rowcol.Pos) {
	field := l.Get(pos.Row, pos.Col)
	if !field.DugOut {
		return
	}
	if field.DigDirection == rowcol.Stand {
		return
	}
	pos = pos.AddDir(field.DigDirection.Right())
	l.FollowDigInterior(pos, 0)
}

func (l *Lagoon) FollowDigInterior(p rowcol.Pos, depth int) {
	if !l.IsValidPosition(p.Row, p.Col) {
		return
	}
	f := l.Get(p.Row, p.Col)
	if f.DugOut {
		return
	}
	if VERBOSE >= 2 {
		fmt.Println(strings.Repeat("    ", depth), "following interior:", p)
	}
	l.MarkDug(p, rowcol.Stand, ColorInterior) // TODO
	// flood fill interior
	for _, dir := range rowcol.Directions {
		l.FollowDigInterior(p.AddDir(dir), depth+1)
	}
}

func (l *Lagoon) Print() {
	// this method operates in grid coordinates !
	rows, cols := l.Rows(), l.Columns()
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			field := l.data.Get(r, c)
			state := "."
			if field.DugOut {
				state = "#"
			}
			fmt.Print(state)
		}
		fmt.Println()
	}
	fmt.Println()
}

func (l *Lagoon) CountDugOut() int {
	// this method operates in grid coordinates !
	var dugout int
	rows, cols := l.Rows(), l.Columns()
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			field := l.data.Get(r, c)
			if field.DugOut {
				dugout++
			}
		}
	}
	return dugout
}
