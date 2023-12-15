package main

// https://adventofcode.com/2023/day/12

import (
	"fmt"
	"strings"

	"git.bind.ch/phil/challenges/lib"
)

var VERBOSE = 1

func main() {
	//ProcessPart1("aoc23/day12/example2.txt")
	ProcessPart1("aoc23/day12/example.txt")
	//ProcessPart1("aoc23/day12/input.txt")
	//
	//ProcessPart2("aoc23/day12/example.txt")
	//ProcessPart2("aoc23/day12/input.txt")
}

func ProcessPart1(name string) {
	fmt.Println("Part 1 input:", name)
	lines := lib.ReadLines(name)
	sequences := ParseSequences(lines)
	var sum int
	for _, s := range sequences {
		fmt.Println(s)
		c := MatchRec(s, 0, 0, 0)
		fmt.Println("combinations:", c)
		sum += c
	}
	fmt.Println("sum:", sum)
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
		lineParts := strings.Split(l, " ")
		if len(lineParts) != 2 {
			panic("invalid input")
		}
		var s Sequence
		//s.Parts = ParseParts(lineParts[0])
		s.PartsTypes = lib.Map([]rune(lineParts[0]), PartTypeOf)
		s.GroupOfDamaged = ParseIntGroups(lineParts[1])
		rv = append(rv, s)
	}
	return rv
}

func ParseIntGroups(s string) []int {
	is := strings.Split(s, ",")
	return lib.Map(is, lib.ToInt)
}

func ParseParts(l string) []Part {
	var rv []Part
	var current Part
	for i, r := range l {
		t := PartTypeOf(r)
		if i == 0 {
			current.T = t
			current.Len = 1
			continue
		}
		if t == current.T {
			current.Len++
			continue
		}
		// switching part type
		rv = append(rv, current)
		// re-initialize current
		current.T = t
		current.Len = 1
	}
	rv = append(rv, current)
	return rv
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

// return index i of the next part to consider
func Match(s Sequence, i, j int) int {
	if i >= len(s.PartsTypes) || j >= len(s.GroupOfDamaged) {
		return -1
	}
	l := s.GroupOfDamaged[j]

	// the first
	if i != 0 {
		// skip over the next part, which must be operational or unknown
		if s.PartsTypes[i] != Operational && s.PartsTypes[i] != Unknown {
			return -1
		}
		i++
	}
	// skip over any operational parts
	for i < len(s.PartsTypes) && s.PartsTypes[i] == Operational {
		i++
	}

	accumulate := 0
	for i < len(s.PartsTypes) {
		if s.PartsTypes[i] == Unknown || s.PartsTypes[i] == Damaged {
			accumulate++
			i++
		}
		if accumulate == l {
			return i + 1
		}
	}
	return -1
}

func MatchRec(s Sequence, i int, j int, depth int) int {
	i = SkipOperational(s, i)
	if VERBOSE >= 2 {
		fmt.Println(strings.Repeat(" ", depth*2), "searching for group", j, "at index", i)
	}
	var c int
	for i < len(s.PartsTypes) {
		if s.PartsTypes[i] == Operational {
			break
		}
		comb, nexti := MatchRecOne(s, i, j, depth+1)
		c += comb
		i = nexti
	}
	return c
}

// sequence must start valid, if it does the first group is matched
// returns: number of combinations and next search position
func MatchRecOne(s Sequence, i int, j int, depth int) (int, int) {
	n := s.GroupOfDamaged[j]
	accumulate := 0
	i = SkipOperational(s, i)
	startmatch := i
	for i < len(s.PartsTypes) && accumulate < n {
		if s.PartsTypes[i] == Operational {
			accumulate = 0
		} else {
			accumulate++
		}
		if accumulate == 1 {
			startmatch = i
		}
		i++
	}
	if accumulate != n {
		return 0, len(s.PartsTypes) // abort search
	}
	if VERBOSE >= 2 {
		fmt.Println(strings.Repeat(" ", depth*2), "matched group", j, "at index", (i - n))
	}
	if j == len(s.GroupOfDamaged)-1 { // last in the group, test if the remaining sequence is
		if RemainingSequenceNotDamaged(s, i) {
			if VERBOSE >= 2 {
				fmt.Println(strings.Repeat(" ", depth*2), "combination found")
			}
			return 1, len(s.PartsTypes) // combination found
		}
		return 0, len(s.PartsTypes)
	}
	// jump to next possible starting point
	if i < len(s.PartsTypes) {
		// invalid next position
		// the current i cannot be used as a starting position
		if s.PartsTypes[i] == Damaged {
			return 0, len(s.PartsTypes)
		}
	}
	i++
	i = SkipOperational(s, i)

	// test remaining 'i' on remaining groups
	//var combinations int
	//var lastMatchIndex int = -1
	//for i < len(s.PartsTypes) {
	//	//c, idx:= MatchRec(s, i, j+1)
	//	//if lastMatchIndex
	//	_ = lastMatchIndex
	//	combinations += MatchRec(s, i, j+1, depth+1)
	//	i++
	//}

	return MatchRec(s, i, j+1, depth+1), startmatch + 1

	//return combinations
}

func SkipOperational(s Sequence, i int) int {
	for i < len(s.PartsTypes) {
		if s.PartsTypes[i] != Operational {
			return i
		}
		i++
	}
	return i
}

func RemainingSequenceNotDamaged(s Sequence, i int) bool {
	for i < len(s.PartsTypes) {
		if s.PartsTypes[i] == Damaged {
			return false
		}
		i++
	}
	return true
}

func AllMatch(s Sequence, i int) bool {
	for j := 0; j < len(s.GroupOfDamaged); j++ {
		i = Match(s, i, j)
		if i == -1 {
			return false
		}
	}
	// test if the remaining sequence is Unknown or Operational parts
	for i < len(s.PartsTypes) {
		if s.PartsTypes[i] == Damaged {
			return false
		}
	}
	return true
}

//func Follow(s Sequence, i, j int) int {
//	combinations := 0
//	for {
//		// test if all number groups match
//		if !AllMatch(s, i, j) {
//			return combinations
//		}
//		combinations++
//		i++
//	}
//	return combinations
//}
