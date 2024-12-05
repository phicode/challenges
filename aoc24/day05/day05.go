package main

// https://adventofcode.com/2024/day/5

import (
	"fmt"
	"strings"

	"git.bind.ch/phil/challenges/lib"
	"git.bind.ch/phil/challenges/lib/assert"
)

var VERBOSE = 1

func main() {
	lib.Timed("Part 1", ProcessPart1, "aoc24/day05/example.txt")
	lib.Timed("Part 1", ProcessPart1, "aoc24/day05/input.txt")

	lib.Timed("Part 2", ProcessPart2, "aoc24/day05/example.txt")
	lib.Timed("Part 2", ProcessPart2, "aoc24/day05/input.txt")
}

func ProcessPart1(name string) {
	fmt.Println("Part 1 input:", name)
	input := ParseInput(name)
	result := SolvePart1(input)
	fmt.Println("Result:", result)
}

func ProcessPart2(name string) {
	fmt.Println("Part 2 input:", name)
	input := ParseInput(name)
	result := SolvePart2(input)
	fmt.Println("Result:", result)
}

func log(v int, msg string) {
	if v <= VERBOSE {
		fmt.Println(msg)
	}
}

////////////////////////////////////////////////////////////

// A before B
type Rule struct {
	A, B int
}

// two numbers where A < B
type RuleKey struct {
	A, B int
}

func NewKey(a, b int) RuleKey {
	if a < b {
		return RuleKey{a, b}
	}
	return RuleKey{b, a}
}

type Manual struct {
	Pages []int
}

type Input struct {
	Rules      []Rule
	RulesIndex map[RuleKey]Rule
	Manuals    []Manual
}

func (in Input) AddRule(r Rule) {
	in.Rules = append(in.Rules, r)
	in.RulesIndex[NewKey(r.A, r.B)] = r
}

func ParseInput(name string) Input {
	lines := lib.ReadLines(name)
	rv := Input{RulesIndex: make(map[RuleKey]Rule)}
	rules := true
	for _, line := range lines {
		if line == "" {
			rules = false
			continue
		}
		if rules {
			rv.AddRule(ParseRule(line))
		} else {
			rv.Manuals = append(rv.Manuals, ParseManual(line))
		}
	}
	return rv
}

func ParseRule(line string) Rule {
	var a, b int
	n, err := fmt.Sscanf(line, "%d|%d", &a, &b)
	assert.True(n == 2 && err == nil)
	return Rule{A: a, B: b}
}

func ParseManual(line string) Manual {
	parts := strings.Split(line, ",")
	return Manual{Pages: lib.Map(parts, lib.ToInt)}
}

////////////////////////////////////////////////////////////

func SolvePart1(input Input) int {
	total := 0
	for _, man := range input.Manuals {
		valid := input.IsValid(man)
		if valid {
			total += man.Pages[len(man.Pages)/2]
		}
	}
	return total
}

func (in Input) IsValid(m Manual) bool {
	for i, a := range m.Pages {
		for _, b := range m.Pages[i+1:] {
			r, found := in.FindRule(a, b)
			if found && r.IsInvalid(a, b) {
				return false
			}
		}
	}
	return true
}

func (in Input) FindRule(a int, b int) (Rule, bool) {
	key := NewKey(a, b)
	r, found := in.RulesIndex[key]
	return r, found
}

func (r Rule) IsInvalid(a, b int) bool {
	return r.A != a || r.B != b
}

////////////////////////////////////////////////////////////

func SolvePart2(input Input) int {
	total := 0
	for _, man := range input.Manuals {
		valid := input.IsValid(man)
		if !valid {
			input.Reorder(man)
			total += man.Pages[len(man.Pages)/2]
		}
	}
	return total
}

func (in Input) Reorder(m Manual) {
	// fmt.Println("Invalid:", m)
	// naive approach: find violating rules and swap the two offending places
	p := m.Pages
	l := len(p)
	for i := 0; i < l; i++ {
		for j := i + 1; j < l; j++ {
			a, b := p[i], p[j]
			r, found := in.FindRule(a, b)
			if found && r.IsInvalid(a, b) {
				//fmt.Println("  swapping:", a, b, "@", i, j)
				p[i], p[j] = p[j], p[i]
			}
		}
	}
	//fmt.Println("  result:", m)
}
