package main

// https://adventofcode.com/2023/day/18

import (
	"fmt"
	"regexp"
	"strconv"

	"git.bind.ch/phil/challenges/lib"
	"git.bind.ch/phil/challenges/lib/rowcol"
)

var VERBOSE = 0

func main() {
	ProcessPart1("aoc23/day18/example.txt") // 62
	ProcessPart1("aoc23/day18/input.txt")   // 67891
	ProcessPart2("aoc23/day18/example.txt") // 952408144115
	ProcessPart2("aoc23/day18/input.txt")   // 94116351948493
}

func ProcessPart1(name string) {
	fmt.Println("Part 1 input:", name)
	lines := lib.ReadLines(name)
	instrs := ParseInstructions(lines)
	ProcessPicksTheorem(instrs)
	fmt.Println()
}

func ProcessPart2(name string) {
	fmt.Println("Part 2 input:", name)
	lines := lib.ReadLines(name)
	instrs := ParsePart2Instructions(lines)
	ProcessPicksTheorem(instrs)
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

func BuildPolygon(instrs []Instruction) []rowcol.Pos {
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

// https://en.wikipedia.org/wiki/Pick%27s_theorem
func ProcessPicksTheorem(instrs []Instruction) {
	poly := BuildPolygon(instrs)
	area := Area(poly)
	peri := Perimeter(instrs)
	if VERBOSE >= 1 {
		fmt.Println("area:", area)
		fmt.Println("peri:", peri)
	}
	totalArea := area + (peri / 2) - 1
	fmt.Println("Total Area:", totalArea)
}

func Perimeter(instrs []Instruction) int {
	var perimeter int // start position
	curDir := instrs[len(instrs)-1].Direction
	for _, instr := range instrs {
		perimeter += instr.Distance
		if curDir.Right() == instr.Direction {
			perimeter++
		} else {
			perimeter--
		}
		curDir = instr.Direction
	}
	return perimeter
}
