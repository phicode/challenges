package main

// https://adventofcode.com/2023/day/12

import (
	"fmt"
	"strings"
	"time"

	"github.com/phicode/challenges/lib"
)

var VERBOSE = 1

func main() {
	VERBOSE = 0
	ProcessPart1("aoc23/day12/example.txt")
	ProcessPart1("aoc23/day12/input.txt")

	ProcessPart2("aoc23/day12/example.txt")
	ProcessPart2("aoc23/day12/input.txt")
}

func ProcessPart1(name string) {
	fmt.Println("Part 1 input:", name)
	Process(name, false)
}

func ProcessPart2(name string) {
	fmt.Println("Part 1 input:", name)
	Process(name, true)
}

func Process(name string, extend bool) {
	lines := lib.ReadLines(name)
	sequences := ParseSequences(lines)
	var sum int
	t0 := time.Now()
	for _, s := range sequences {
		if extend {
			s = ExtendSequenceToPart2(s)
		}
		c := SolveCombinations(s)

		sum += c
		if VERBOSE >= 1 {
			fmt.Println(s)
			fmt.Println("combinations:", c)
		}
	}
	fmt.Println("T:", time.Since(t0))
	fmt.Println("Sum:", sum)
	fmt.Println()
}

type CacheKey struct {
	Pos  int
	Rock int
}
type Cache map[CacheKey]int

func SolveCombinations(s Sequence) int {
	matcher := PartTypesToMatcher(s.GroupOfDamaged)
	cache := make(Cache)
	return matcher.Match(&s, cache, 0)
}

////////////////////////////////////////////////////////////

type PartType int

const (
	Unknown PartType = iota
	Operational
	Damaged
)

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
	default:
		return ""
	}
}

type Matcher interface {
	Match(s *Sequence, cache Cache, pos int) int
	SetNext(matcher Matcher)
}

// OperationalMatcher tests sequences of N or more matches of operational gears
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

func (m *OperationalMatcher) Match(s *Sequence, cache Cache, pos int) int {
	var sum int
	if m.N == 0 {
		// zero or matched allowed, test zero first
		sum += m.Next.Match(s, cache, pos)
	}
	for m.Accept(s, pos) {
		pos++
		sum += m.Next.Match(s, cache, pos)
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

// RockMatcher matches a rock of size N.
type RockMatcher struct {
	N    int // match this many or more part types
	Rock int
	Next Matcher
}

func MakeRockMatcher(n int, rock int, previous Matcher) Matcher {
	matcher := &RockMatcher{
		N:    n,
		Rock: rock,
		Next: nil,
	}
	previous.SetNext(matcher)
	return matcher
}

func (m *RockMatcher) Match(s *Sequence, cache Cache, pos int) int {
	rem := len(s.PartsTypes) - pos
	if rem < m.N {
		return 0
	}
	for i := 0; i < m.N; i++ {
		if !m.Accept(s, pos+i) {
			return 0
		}
	}
	key := CacheKey{
		Pos:  pos,
		Rock: m.Rock,
	}
	if v, found := cache[key]; found {
		return v
	}
	v := m.Next.Match(s, cache, pos+m.N)
	cache[key] = v
	return v
}

func (m *RockMatcher) Accept(s *Sequence, pos int) bool {
	pt := s.PartsTypes[pos]
	return pt == Damaged || pt == Unknown
}

func (m *RockMatcher) SetNext(next Matcher) { m.Next = next }

type EndMatcher struct {
}

func (m *EndMatcher) Match(s *Sequence, _ Cache, pos int) int {
	if pos == len(s.PartsTypes) {
		return 1
	}
	return 0
}
func (m *EndMatcher) SetNext(_ Matcher) {
	panic("the end cannot continue")
}

func PartTypesToMatcher(groups []int) Matcher {
	start := MakeOperationalMatcher(0, nil)
	end := &EndMatcher{}
	current := start
	for i, v := range groups {
		if i > 0 {
			current = MakeOperationalMatcher(1, current)
		}
		current = MakeRockMatcher(v, i, current)
	}
	current = MakeOperationalMatcher(0, current)
	current.SetNext(end)
	return start
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
