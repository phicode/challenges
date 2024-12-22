package main

// https://adventofcode.com/2024/day/21

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/phicode/challenges/lib"
	"github.com/phicode/challenges/lib/assert"
	"github.com/phicode/challenges/lib/math"
	"github.com/phicode/challenges/lib/rowcol"
)

func main() {
	flag.Parse()
	lib.Timed("Part 1", ProcessPart1, "aoc24/day21/example.txt")
	//lib.Timed("Part 1", ProcessPart1, "aoc24/day21/input.txt")
	//
	//lib.Timed("Part 2", ProcessPart2, "aoc24/day21/example.txt")
	//lib.Timed("Part 2", ProcessPart2, "aoc24/day21/input.txt")

	//lib.Profile(1, "day21.pprof", "Part 2", ProcessPart2, "aoc24/day21/input.txt")
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
	Codes []string
}

type NumPad struct {
	V rune
	P rowcol.Pos
}
type DirPad struct {
	V rune
	P rowcol.Pos
}

func ReadAndParseInput(name string) Input {
	lines := lib.ReadLines(name)
	return ParseInput(lines)
}

func ParseInput(lines []string) Input {
	return Input{Codes: lines}
}

////////////////////////////////////////////////////////////

func SolvePart1(input Input) int {
	total := 0
	for _, code := range input.Codes {
		np := NumPad{'A', NumPadPosition('A')}
		dp1 := DirPad{'A', DirPadPosition('A')}
		dp2 := DirPad{'A', DirPadPosition('A')}

		outnp := np.PressMultiple([]rune(code))
		outdp1 := dp1.PressMultiple(outnp)
		outdp2 := dp2.PressMultiple(outdp1)
		fmt.Println(string(outdp2))
		complexity := CodeComplexity(code, string(outdp2))
		fmt.Printf("  len: %d, complexity: %d\n", len(outdp2), complexity)
		total += complexity
	}
	return total
}

func SolvePart1Combinations(input Input) int {
	total := 0
	for _, code := range input.Codes {
		np := NumPad{'A', NumPadPosition('A')}
		//for _,
		//np.PressCombinations()

		dp1 := DirPad{'A', DirPadPosition('A')}
		dp2 := DirPad{'A', DirPadPosition('A')}

		outnp := np.PressMultiple([]rune(code))
		outdp1 := dp1.PressMultiple(outnp)
		outdp2 := dp2.PressMultiple(outdp1)
		fmt.Println(string(outdp2))
		complexity := CodeComplexity(code, string(outdp2))
		fmt.Printf("  len: %d, complexity: %d\n", len(outdp2), complexity)
		total += complexity
	}
	return total
}

func NumericCodePart(code string) int {
	l := len(code)
	assert.True(code[l-1] == 'A')
	code = code[:l-1]
	if code[0] == '0' {
		code = code[1:]
	}
	num, err := strconv.Atoi(code)
	assert.True(err == nil)
	return num
}

func (np *NumPad) PressMultiple(code []rune) []rune {
	var rv []rune
	for _, press := range code {
		rv = append(rv, np.Press(press)...)
	}
	return rv
}

func (dp *DirPad) PressMultiple(code []rune) []rune {
	var rv []rune
	for _, press := range code {
		rv = append(rv, dp.Press(press)...)
	}
	return rv
}

var invalidNumPadPos = rowcol.Pos{Row: 3, Col: 0}
var invalidDirPadPos = rowcol.Pos{Row: 0, Col: 0}

func (np *NumPad) PressCombinations(c rune) Combinations {
	var cs Combinations
	pos := NumPadPosition(c)
	cur := np.P
	move := pos.Sub(cur)
	if move.Row == 0 {
		cs.AddLeftRight(move.Col)
	} else if move.Col == 0 {
		cs.AddUpDown(move.Row)
	} else {
		// combination of up-down/left-right moves
		// max x=2, y=3
		lrDir := rowcol.Right
		if move.Col < 0 {
			lrDir = rowcol.Left
		}
		udDir := rowcol.Down
		if move.Row < 0 {
			udDir = rowcol.Up
		}
		lrN := math.AbsInt(move.Col)
		udN := math.AbsInt(move.Row)
		cs.AddX(cur, invalidNumPadPos, lrDir, udDir, lrN, udN)
	}

	cs.Add('A')
	np.P = pos
	np.V = c
	return cs
}

type Combinations struct {
	Cs [][]rune
}

func (cs *Combinations) Add(c rune) {
	if len(cs.Cs) == 0 {
		cs.Cs = [][]rune{
			[]rune{c},
		}
		return
	}
	for i, _ := range cs.Cs {
		cs.Cs[i] = append(cs.Cs[i], c)
	}
}

func (cs *Combinations) Copy() Combinations {
	return Combinations{
		Cs: append([][]rune(nil), cs.Cs...),
	}
}

func (cs *Combinations) Combine(b Combinations) {
	cs.Cs = append(cs.Cs, b.Cs...)
}

func (cs *Combinations) AddLeftRight(n int) {
	c := '>'
	if n < 0 {
		c = '<'
		n = -n
	}
	assert.True(n > 0)
	for i := 0; i < n; i++ {
		cs.Add(c)
	}
}

func (cs *Combinations) AddUpDown(n int) {
	c := 'v'
	if n < 0 {
		c = '^'
		n = -n
	}
	assert.True(n > 0)
	for i := 0; i < n; i++ {
		cs.Add(c)
	}
}

func (cs *Combinations) AddX(pos, not rowcol.Pos, lrDir, udDir rowcol.Direction, lrN, udN int) {
	if lrN == 0 && udN == 0 {
		return
	}
	posLr := pos.AddDir(lrDir)
	posUd := pos.AddDir(udDir)

	if posLr == not || lrN == 0 {
		// only ud moves
		if udN == 0 {
			return
		}
		cs.Add(rune(udDir.PrintChar()))
		cs.AddX(posUd, not, lrDir, udDir, lrN, udN-1)
		return
	}

	if posUd == not || udN == 0 {
		// only lr moves
		if lrN == 0 {
			return
		}
		cs.Add(rune(lrDir.PrintChar()))
		cs.AddX(posLr, not, lrDir, udDir, lrN-1, udN)
		return
	}

	assert.True(lrN > 0 && udN > 0)
	b := cs.Copy()
	cs.Add(rune(udDir.PrintChar()))
	cs.AddX(posUd, not, lrDir, udDir, lrN, udN-1)
	b.Add(rune(lrDir.PrintChar()))
	b.AddX(posLr, not, lrDir, udDir, lrN-1, udN)
	cs.Combine(b)
}

func (np *NumPad) Press(c rune) []rune {
	pos := NumPadPosition(c)
	cur := np.P
	move := pos.Sub(cur)
	moves := make([]rune, 0, 5)
	if move.Row < 0 { // move up first
		moves = AddUpDownMoves(moves, move.Row)
		moves = AddLeftRightMoves(moves, move.Col)
	} else {
		moves = AddLeftRightMoves(moves, move.Col)
		moves = AddUpDownMoves(moves, move.Row)
	}
	moves = append(moves, 'A')
	np.P = pos
	np.V = c
	return moves
}

func (dp *DirPad) Press(c rune) []rune {
	pos := DirPadPosition(c)
	cur := dp.P
	move := pos.Sub(cur)
	moves := make([]rune, 0, 5)
	if move.Row >= 0 { // move down first
		moves = AddUpDownMoves(moves, move.Row)
		moves = AddLeftRightMoves(moves, move.Col)
	} else if move.Col >= 0 { // move right first
		moves = AddLeftRightMoves(moves, move.Col)
		moves = AddUpDownMoves(moves, move.Row)
	} else {
		moves = AddUpDownMoves(moves, move.Row)
		moves = AddLeftRightMoves(moves, move.Col)
	}
	moves = append(moves, 'A')
	dp.P = pos
	dp.V = c
	return moves
}

func AddUpDownMoves(moves []rune, n int) []rune {
	c := 'v'
	if n < 0 {
		c = '^'
		n = -n
	}
	// navigate up first
	for i := 0; i < n; i++ {
		moves = append(moves, c)
	}
	return moves
}

func AddLeftRightMoves(moves []rune, n int) []rune {
	c := '>'
	if n < 0 {
		c = '<'
		n = -n
	}
	// navigate up first
	for i := 0; i < n; i++ {
		moves = append(moves, c)
	}
	return moves
}

func NumPadPosition(c rune) rowcol.Pos {
	switch c {
	case '7':
		return rowcol.Pos{Row: 0, Col: 0}
	case '8':
		return rowcol.Pos{Row: 0, Col: 1}
	case '9':
		return rowcol.Pos{Row: 0, Col: 2}

	case '4':
		return rowcol.Pos{Row: 1, Col: 0}
	case '5':
		return rowcol.Pos{Row: 1, Col: 1}
	case '6':
		return rowcol.Pos{Row: 1, Col: 2}

	case '1':
		return rowcol.Pos{Row: 2, Col: 0}
	case '2':
		return rowcol.Pos{Row: 2, Col: 1}
	case '3':
		return rowcol.Pos{Row: 2, Col: 2}

	case '0':
		return rowcol.Pos{Row: 3, Col: 1}
	case 'A':
		return rowcol.Pos{Row: 3, Col: 2}
	default:
		panic("invalid numpad rune")
	}
}

func DirPadPosition(c rune) rowcol.Pos {
	switch c {
	case '^':
		return rowcol.Pos{Row: 0, Col: 1}
	case 'A':
		return rowcol.Pos{Row: 0, Col: 2}

	case '<':
		return rowcol.Pos{Row: 1, Col: 0}
	case 'v':
		return rowcol.Pos{Row: 1, Col: 1}
	case '>':
		return rowcol.Pos{Row: 1, Col: 2}
	default:
		panic("invalid numpad rune")
	}
}

func CodeComplexity(incode, outcode string) int {
	return len(outcode) * NumericCodePart(incode)
}

////////////////////////////////////////////////////////////

func SolvePart2(input Input) int {
	return 0
}
