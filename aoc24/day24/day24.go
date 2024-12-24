package main

// https://adventofcode.com/2024/day/24

import (
	"flag"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/phicode/challenges/lib"
	"github.com/phicode/challenges/lib/assert"
)

func main() {
	flag.Parse()
	lib.Timed("Part 1", ProcessPart1, "aoc24/day24/example.txt")
	lib.Timed("Part 1", ProcessPart1, "aoc24/day24/example2.txt")
	lib.Timed("Part 1", ProcessPart1, "aoc24/day24/input.txt")
	//
	//lib.Timed("Part 2", ProcessPart2, "aoc24/day24/example.txt")
	//lib.Timed("Part 2", ProcessPart2, "aoc24/day24/input.txt")

	//lib.Profile(1, "day24.pprof", "Part 2", ProcessPart2, "aoc24/day24/input.txt")
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
	StartWires map[string]int
	Gates      []Gate
}

func And(a, b int) int { return a & b }
func Or(a, b int) int  { return a | b }
func Xor(a, b int) int { return a ^ b }

type Op func(a, b int) int

var Operations = map[string]Op{
	"AND": And,
	"OR":  Or,
	"XOR": Xor,
}

type Gate struct {
	Op   Op
	A, B string
	Out  string
}

func ReadAndParseInput(name string) Input {
	lines := lib.ReadLines(name)
	return ParseInput(lines)
}

func ParseInput(lines []string) Input {
	divider := slices.Index(lines, "")
	assert.True(divider > 0)
	return Input{
		StartWires: ParseStartWires(lines[:divider]),
		Gates:      ParseGates(lines[divider+1:]),
	}
}

func ParseStartWires(lines []string) map[string]int {
	m := make(map[string]int)
	for _, line := range lines {
		parts := strings.Split(line, ": ")
		assert.True(len(parts) == 2)
		value, err := strconv.Atoi(parts[1])
		assert.True(err == nil)
		m[parts[0]] = value
	}
	return m
}

func ParseGates(lines []string) []Gate {
	var gates []Gate
	for _, line := range lines {
		var a, b, op, out string
		n, err := fmt.Sscanf(line, "%s %s %s -> %s", &a, &op, &b, &out)
		assert.True(n == 4 && err == nil)
		opf, found := Operations[op]
		assert.True(found)
		gates = append(gates, Gate{
			Op:  opf,
			A:   a,
			B:   b,
			Out: out,
		})
	}
	return gates
}

////////////////////////////////////////////////////////////

type State struct {
	values    map[string]int
	byOutput  map[string]Gate
	allwires  []string
	allzwires []string
}

func (s State) pull(wire string) int {
	state, ok := s.values[wire]
	if ok {
		return state
	}
	gate := s.byOutput[wire]
	va := s.pull(gate.A)
	vb := s.pull(gate.B)
	value := gate.Op(va, vb)
	s.values[wire] = value
	return value
}

func (s State) Values() []string {
	var values []string
	for wire, value := range s.values {
		if wire[0] == 'x' || wire[0] == 'y' {
			continue
		}
		values = append(values, fmt.Sprintf("%s: %d", wire, value))
	}
	return values
}

func NewState(input Input) *State {
	values := make(map[string]int)
	allwires := make(map[string]bool)
	allzwires := make(map[string]bool)
	byOutput := make(map[string]Gate)
	for wire, value := range input.StartWires {
		values[wire] = value
		allwires[wire] = true
	}
	for _, gate := range input.Gates {
		addWire(gate.A, allwires, allzwires)
		addWire(gate.B, allwires, allzwires)
		addWire(gate.Out, allwires, allzwires)
		_, found := byOutput[gate.Out]
		assert.False(found)
		byOutput[gate.Out] = gate
	}
	s := &State{
		values:    values,
		byOutput:  byOutput,
		allwires:  lib.MapKeys(allwires),
		allzwires: lib.MapKeys(allzwires),
	}
	slices.Sort(s.allzwires)
	return s
}

func addWire(name string, allwires, allzwires map[string]bool) {
	allwires[name] = true
	if name[0] == 'z' {
		allzwires[name] = true
	}
}

func SolvePart1(input Input) int {
	state := NewState(input)
	out := 0
	for i, zwire := range state.allzwires {
		expectName := fmt.Sprintf("z%02d", i)
		assert.True(zwire == expectName)

		v := state.pull(zwire)
		//fmt.Println(zwire, v)
		out = out | (v << i)
	}

	//allvalues := state.Values()
	//slices.Sort(allvalues)
	//fmt.Println("ALL Wires")
	//for _, v := range allvalues {
	//	fmt.Println("   ", v)
	//}

	return out
}

////////////////////////////////////////////////////////////

func SolvePart2(input Input) int {
	return 0
}
