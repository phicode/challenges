package main

// https://adventofcode.com/2023/day/12

import (
	"fmt"
	"runtime"
	"strings"

	"git.bind.ch/phil/challenges/lib"
)

var VERBOSE = 1

func main() {
	VERBOSE = 0
	ProcessPart1("aoc23/day12/example.txt")
	ProcessPart1("aoc23/day12/input.txt")

	ProcessPart2("aoc23/day12/example.txt")
	//ProcessPart2("aoc23/day12/input.txt")
}

func ProcessPart1(name string) {
	fmt.Println("Part 1 input:", name)
	lines := lib.ReadLines(name)
	sequences := ParseSequences(lines)
	var sum, sumalt int
	for _, s := range sequences {
		c, alt := SolveCombinations(s)
		//if c != alt {
		//	fmt.Println(s)
		//	fmt.Println(c, "!=", alt)
		//	//return
		//}
		sum += c
		sumalt += alt
		if VERBOSE >= 1 {
			fmt.Println(s)
			fmt.Println("combinations:", c)
		}
	}
	fmt.Println("Sum:", sum)
	fmt.Println("Sum alt:", sumalt)
	fmt.Println()
}

func SolveCombinations(s Sequence) (int, int) {
	matcher, n := PartTypesToMatcher(s.GroupOfDamaged)
	c := matcher.Match(&s, 0)
	return *n, c
}

func ProcessPart2(name string) {
	fmt.Println("Part 2 input:", name)
	lines := lib.ReadLines(name)
	sequences := ParseSequences(lines)

	results := make(chan int, 1)
	jobs := make(chan Sequence, 1)

	// job generator
	go func() {
		for _, s := range sequences {
			extended := ExtendSequenceToPart2(s)
			jobs <- extended
		}
	}()

	// solvers
	parallel := max(4, runtime.NumCPU())
	for i := 0; i < parallel; i++ {
		go func() {
			for seq := range jobs {
				c, _ := SolveCombinations(seq)
				if VERBOSE >= 1 {
					fmt.Println(seq)
					fmt.Println("combinations:", c)
				}
				results <- c
			}
		}()
	}

	// aggregator
	var sum int
	for i := 0; i < len(sequences); i++ {
		result := <-results
		sum += result
	}
	fmt.Println("Sum:", sum)
	fmt.Println()
}

func log(v int, msg string) {
	if v <= VERBOSE {
		fmt.Println(msg)
	}
}

////////////////////////////////////////////////////////////

type PartType int

const (
	Unknown PartType = iota
	Operational
	Damaged
	Undefined
)

type Part struct {
	T   PartType
	Len int
}

func (p Part) String() string {
	return fmt.Sprintf("%s*%d", p.T, p.Len)
}

type Sequence struct {
	//Parts          []Part
	PartsTypes     []PartType
	GroupOfDamaged []int
}

func ParseSequences(lines []string) []Sequence {
	var rv []Sequence
	for _, l := range lines {
		rv = append(rv, ParseSequence(l))
	}
	return rv
}
func ParseSequence(s string) Sequence {
	lineParts := strings.Split(s, " ")
	if len(lineParts) != 2 {
		panic("invalid input")
	}
	var seq Sequence
	//s.Parts = ParseParts(lineParts[0])
	seq.PartsTypes = lib.Map([]rune(lineParts[0]), PartTypeOf)
	seq.GroupOfDamaged = ParseIntGroups(lineParts[1])
	return seq
}

func ParseIntGroups(s string) []int {
	is := strings.Split(s, ",")
	return lib.Map(is, lib.ToInt)
}

func PartTypeOf(r rune) PartType {
	switch r {
	case '.':
		return Operational
	case '#':
		return Damaged
	case '?':
		return Unknown
	}
	panic(fmt.Errorf("invalid part type: '%c'", r))
}

func (t PartType) String() string {
	switch t {
	case Operational:
		return "."
	case Damaged:
		return "#"
	case Unknown:
		return "?"
	}
	return ""
}

type Matcher interface {
	Match(s *Sequence, pos int) int
	SetNext(matcher Matcher)
}

// test that a sequence of n or more occurence match
type OperationalMatcher struct {
	N    int // match this many or more part types
	Next Matcher
}

func MakeOperationalMatcher(n int, previous Matcher) Matcher {
	matcher := &OperationalMatcher{
		N:    n,
		Next: nil,
	}
	if previous != nil {
		previous.SetNext(matcher)
	}
	return matcher
}

func (m *OperationalMatcher) Match(s *Sequence, pos int) int {
	var sum int
	if m.N == 0 {
		// zero or matched allowed, test zero first
		sum += m.Next.Match(s, pos)
	}
	for m.Accept(s, pos) {
		pos++
		sum += m.Next.Match(s, pos)
	}
	return sum
}

func (m *OperationalMatcher) Accept(s *Sequence, pos int) bool {
	if pos >= len(s.PartsTypes) {
		return false
	}
	pt := s.PartsTypes[pos]
	return pt == Operational || pt == Unknown
}

func (m *OperationalMatcher) SetNext(next Matcher) { m.Next = next }

// test that a sequence of n matches
type RockMatcher struct {
	N    int // match this many or more part types
	Next Matcher
}

func MakeRockMatcher(n int, previous Matcher) Matcher {
	matcher := &RockMatcher{
		N:    n,
		Next: nil,
	}
	previous.SetNext(matcher)
	return matcher
}

func (m *RockMatcher) Match(s *Sequence, pos int) int {
	rem := len(s.PartsTypes) - pos
	if rem < m.N {
		return 0
	}
	for i := 0; i < m.N; i++ {
		if !m.Accept(s, pos+i) {
			return 0
		}
	}
	return m.Next.Match(s, pos+m.N)
}

func (m *RockMatcher) Accept(s *Sequence, pos int) bool {
	pt := s.PartsTypes[pos]
	return pt == Damaged || pt == Unknown
}

func (m *RockMatcher) SetNext(next Matcher) { m.Next = next }

type EndMatcher struct {
	Count int
}

func (m *EndMatcher) Match(s *Sequence, pos int) int {
	if pos == len(s.PartsTypes) {
		m.Count++
	}
	return 1
}
func (m *EndMatcher) SetNext(_ Matcher) {
	panic("the end cannot continue")
}

func PartTypesToMatcher(groups []int) (Matcher, *int) {
	start := MakeOperationalMatcher(0, nil)
	end := &EndMatcher{}
	current := start
	for i, v := range groups {
		if i > 0 {
			current = MakeOperationalMatcher(1, current)
		}
		current = MakeRockMatcher(v, current)
	}
	current = MakeOperationalMatcher(0, current)
	current.SetNext(end)
	return start, &end.Count
}

func MakeLinearMatcher(groups []int) []Matcher {
	n := len(groups)
	matchers := make([]Matcher, 0, n+n+1)                             // every number, space in between, start and end
	matchers = append(matchers, &OperationalMatcher{N: 0, Next: nil}) // start
	for i, g := range groups {
		if i > 0 {
			matchers = append(matchers, &OperationalMatcher{N: 1, Next: nil})
		}
		matchers = append(matchers, &RockMatcher{N: g, Next: nil})
	}
	// i=0
	matchers = append(matchers, &OperationalMatcher{N: 0, Next: nil}) // end
	return matchers
}

func RunMatchers(ms []Matcher) {

}

func ExtendSequenceToPart2(s Sequence) Sequence {
	numParts := len(s.PartsTypes)
	numGroups := len(s.GroupOfDamaged)
	e := Sequence{
		PartsTypes:     make([]PartType, 0, 5*numParts+4),
		GroupOfDamaged: make([]int, 0, 5*numGroups),
	}
	for i := 0; i < 5; i++ {
		if i > 0 {
			e.PartsTypes = append(e.PartsTypes, Unknown)
		}
		e.PartsTypes = append(e.PartsTypes, s.PartsTypes...)
		e.GroupOfDamaged = append(e.GroupOfDamaged, s.GroupOfDamaged...)
	}
	return e
}
