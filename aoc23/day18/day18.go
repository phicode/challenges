package main

// https://adventofcode.com/2023/day/18

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"git.bind.ch/phil/challenges/lib"
	"git.bind.ch/phil/challenges/lib/rowcol"
)

var VERBOSE = 1

func main() {
	ProcessPart1_FirstSolution("aoc23/day18/example.txt") // 62
	ProcessPart1("aoc23/day18/example.txt")               // 62
	ProcessPart1("aoc23/day18/input.txt")                 // 67891

	ProcessPart2("aoc23/day18/example.txt") // 952408144115
	ProcessPart2("aoc23/day18/input.txt")
}

func ProcessPart1(name string) {
	fmt.Println("Part 1 input:", name)
	lines := lib.ReadLines(name)
	instrs := ParseInstructions(lines)
	Process(instrs)
}

func ProcessPart1_FirstSolution(name string) {
	fmt.Println("Part 1 input:", name)
	lines := lib.ReadLines(name)
	instrs := ParseInstructions(lines)
	if VERBOSE >= 1 {
		for _, instr := range instrs {
			fmt.Println(instr)
		}
	}

	lagoon := CreateLagoon(instrs)
	lagoon.DigPath(instrs)
	if VERBOSE >= 1 {
		lagoon.Print()
	}

	lagoon.DigInterior()
	if VERBOSE >= 1 {
		lagoon.Print()
	}

	area := lagoon.CountDugOut()
	fmt.Println("Area:", area)
	fmt.Println()
}

func ProcessPart2(name string) {
	fmt.Println("Part 2 input:", name)
	lines := lib.ReadLines(name)
	instrs := ParsePart2Instructions(lines)
	Process(instrs)
}

func Process(instrs []Instruction) {
	if VERBOSE >= 2 {
		for _, instr := range instrs {
			fmt.Println(instr)
		}
	}
	polygon := BuildPolygon(instrs)
	area := Area(polygon)
	fmt.Println("Area:", area)
	fmt.Println()
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

func ParsePart2Instructions(lines []string) []Instruction {
	var rv []Instruction
	for _, l := range lines {
		rv = append(rv, ParsePart2Instruction(l))
	}
	return rv
}

var pattern = regexp.MustCompile(`[UDLR] [0-9]+ \(#([a-z0-9]+)\)`)

func ParsePart2Instruction(l string) Instruction {
	results := pattern.FindStringSubmatch(l)
	if len(results) != 2 { // the entire pattern plus the matched group
		panic(fmt.Errorf("unexpected match: %v", results))
	}
	instruction := results[1]
	directionRune := instruction[5]
	distanceHex := instruction[:5]
	var direction rowcol.Direction
	switch directionRune {
	case '0':
		direction = rowcol.Right
	case '1':
		direction = rowcol.Down
	case '2':
		direction = rowcol.Left
	case '3':
		direction = rowcol.Up
	default:
		panic("invalid direction in hex code")
	}
	distance, err := strconv.ParseInt(distanceHex, 16, 63)
	if err != nil {
		panic(err)
	}
	return Instruction{
		Direction: direction,
		Distance:  int(distance),
		Color:     0,
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
	if i.Color == 0 {
		return fmt.Sprintf("%s %d", i.Direction, i.Distance)
	} else {
		return fmt.Sprintf("%s %d (#%x)", i.Direction, i.Distance, i.Color)
	}
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
	if VERBOSE >= 1 {
		fmt.Println("lagoon size:", size)
		fmt.Println("lagoon minimum:", posMin)
		fmt.Println("lagoon maximum:", posMax)
	}

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

func BuildPath(instrs []Instruction) []rowcol.Pos {
	rv := make([]rowcol.Pos, 0, len(instrs)+1)
	start := rowcol.Pos{}
	rv = append(rv, start)
	for _, instr := range instrs {
		vec := instr.Direction.MulS(instr.Distance)
		end := start.Add(vec)
		rv = append(rv, end)
		start = end
	}
	return rv
}

func BuildPolygon(instrs []Instruction) []rowcol.Pos {
	path := BuildPath(instrs)
	path = path[:len(path)-1]
	l := len(path)
	var rv []rowcol.Pos
	prev := path[l-1]
	for i := 0; i < l; i++ {
		current := path[i]
		next := path[(i+1)%l]
		dirA := Direction(prev, current)
		dirB := Direction(current, next)
		pc := ToPointCoordinate(current, dirA, dirB)
		if VERBOSE >= 2 {
			fmt.Printf("%v(%s + %s) = %v\n", current, dirA, dirB, pc)
		}
		rv = append(rv, pc)

		prev = current
	}
	rv = append(rv, rv[0])
	return rv
}

var (
	PointCoordinateTopLeft     = rowcol.Pos{Col: 0, Row: 0}
	PointCoordinateBottomLeft  = rowcol.Pos{Col: 0, Row: 1}
	PointCoordinateTopRight    = rowcol.Pos{Col: 1, Row: 0}
	PointCoordinateBottomRight = rowcol.Pos{Col: 1, Row: 1}
)

func ToPointCoordinate(current rowcol.Pos, a, b rowcol.Direction) rowcol.Pos {
	if a == rowcol.Left {
		if b == rowcol.Up {
			// (-1,0),(0,-1) => (0,1)
			return current.Add(PointCoordinateBottomLeft)
		}
		if b == rowcol.Down {
			return current.Add(PointCoordinateBottomRight)
		}
		panic("invalid state")
	}
	if a == rowcol.Right {
		if b == rowcol.Up {
			return current.Add(PointCoordinateTopLeft)
		}
		if b == rowcol.Down {
			return current.Add(PointCoordinateTopRight)
		}
		panic("invalid state")
	}
	if a == rowcol.Up {
		if b == rowcol.Left {
			return current.Add(PointCoordinateBottomLeft)
		}
		if b == rowcol.Right {
			return current.Add(PointCoordinateTopLeft)
		}
		panic("invalid state")
	}
	if a == rowcol.Down {
		if b == rowcol.Left {
			return current.Add(PointCoordinateBottomRight)
		}
		if b == rowcol.Right {
			return current.Add(PointCoordinateTopRight)
		}
		panic("invalid state")
	}
	panic(fmt.Errorf("invalid state, a=%v, b=%v", a, b))
}

func Direction(from rowcol.Pos, to rowcol.Pos) rowcol.Direction {
	dir := to.Sub(from)
	if dir.Row != 0 && dir.Col != 0 {
		panic("invalid state, row or col must be zero")
	}
	dir.Row = clamp(dir.Row, -1, 1)
	dir.Col = clamp(dir.Col, -1, 1)
	return rowcol.Direction(dir)
}

func clamp(value, _min, _max int) int {
	value = max(value, _min)
	value = min(value, _max)
	return value
}

func Area(ps []rowcol.Pos) int {
	area := Shoelace(ps)
	if area < 0 {
		area = -area
	}
	return area / 2
}

// negative results: points are clockwise
// positive results: points are counter clockwise
func Shoelace(ps []rowcol.Pos) int {
	var rv int
	a := ps[len(ps)-1]
	for _, b := range ps {
		rv += a.Col * b.Row
		rv -= b.Col * a.Row
		a = b
	}
	return rv
}

func IsCounterClockwise(ps []rowcol.Pos) bool {
	return Shoelace(ps) > 0
}

func PrintPolygon(polygon []rowcol.Pos) {
	l := len(polygon)
	for i := 1; i < l; i++ {
		a, b := polygon[i-1], polygon[i]
		fmt.Println(a, "-->", b)
	}
	fmt.Println("counter clockwise:", IsCounterClockwise(polygon))
}
