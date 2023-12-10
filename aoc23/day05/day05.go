package main

// https://adventofcode.com/2023/day/5

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"git.bind.ch/phil/challenges/lib"
)

var debug = 0

func main() {
	ProcessPart1("aoc23/day05/example.txt")
	ProcessPart1("aoc23/day05/input.txt")

	ProcessPart2("aoc23/day05/example.txt")
	t := time.Now()
	ProcessPart2("aoc23/day05/input.txt")
	e := time.Now().Sub(t)
	fmt.Println(e)
}

func ProcessPart1(name string) {
	fmt.Println("input:", name)
	lines := lib.ReadLines(name)

	input := ParseInput(lines)
	input.Run()

	fmt.Println()
}

func ProcessPart2(name string) {
	fmt.Println("input:", name)
	lines := lib.ReadLines(name)

	input := ParseInput(lines)
	input.RunProcessPart2()

	fmt.Println()
}

////////////////////////////////////////////////////////////

func ParseInput(ls []string) *Input {
	input := &Input{}
	var currentMap *Map
	for i, l := range ls {
		if i == 0 {
			ints, found := strings.CutPrefix(l, "seeds: ")
			if !found {
				panic("invalid input")
			}
			input.Seeds = lib.ExtractInts(ints)
			continue
		}
		if l == "" {
			currentMap = nil
			continue
		}
		if currentMap == nil {
			mapping, found := strings.CutSuffix(l, " map:")
			parts := strings.Split(mapping, "-to-")
			if !found || len(parts) != 2 {
				panic(fmt.Errorf("invalid input: %q, mapping: %q", l, mapping))
			}
			currentMap = &Map{
				From: parts[0],
				To:   parts[1],
			}
			input.Maps = append(input.Maps, currentMap)
			continue
		}
		var r Translation
		if n, err := fmt.Sscanf(l, "%d %d %d", &r.Dst, &r.Range.Start, &r.Range.N); n != 3 || err != nil {
			panic(fmt.Errorf("invalid range input, n=%d, err=%w", n, err))
		}
		currentMap.Translations = append(currentMap.Translations, r)
	}
	input.SortRanges()
	return input
}

type Input struct {
	Seeds []int
	Maps  []*Map
}

func (in *Input) SortRanges() {
	for _, m := range in.Maps {
		sort.Sort(TranslationBySource(m.Translations))
	}
}

func (in *Input) Run() {
	var lowest int
	for i, seed := range in.Seeds {
		v := in.TranslateSeed(seed)
		if i == 0 {
			lowest = v
		} else {
			lowest = min(lowest, v)
		}
	}
	fmt.Println("Lowest location:", lowest)
}

func (in *Input) TranslateSeed(value int) int {
	t := "seed"
	for t != "location" {
		t, value = in.Translate(t, value)
	}
	return value
}

func (in *Input) Translate(t string, value int) (string, int) {
	m := in.FindMap(t)
	out := m.Translate(value)
	if debug >= 2 {
		fmt.Printf("%s -> %s: %d -> %d\n", m.From, m.To, value, out)
	}
	return m.To, out
}

func (in *Input) FindMap(t string) *Map {
	for _, m := range in.Maps {
		if m.From == t {
			return m
		}
	}
	panic(fmt.Errorf("no map found for %q", t))
}

type Range struct {
	Start int
	N     int
}

// inclusive maximum value
func (t Range) Max() int {
	return t.Start + t.N - 1
}

type Translation struct {
	Range
	Dst int
}

func (r Range) Matches(x int) bool {
	return x >= r.Start && x <= r.Max()
}
func (r Range) MatchesRange(b Range) bool {
	return !(b.Max() < r.Start || r.Max() < b.Start)
}

func (r Translation) Translate(x int) int {
	offset := r.Dst - r.Start
	return x + offset
}

type Map struct {
	From         string
	To           string
	Translations []Translation // sorted by Source field
}

func (m *Map) Translate(x int) int {
	l := len(m.Translations)
	if x < m.Translations[0].Start {
		return x
	}
	last := m.Translations[l-1]
	if x > last.Max() {
		return x
	}

	for _, t := range m.Translations {
		if t.Matches(x) {
			return t.Translate(x)
		}
	}

	// binary search version, not really needed for such small input sets
	//searchFunc := func(idx int) bool { return m.Translations[idx].Src >= x }
	//pos := sort.Search(l, searchFunc)
	//// pos is either an exact match, or one above
	//if pos != l {
	//	if m.Translations[pos].Matches(x) {
	//		return m.Translations[pos].Translate(x)
	//	}
	//}
	//if pos > 0 {
	//	if m.Translations[pos-1].Matches(x) {
	//		return m.Translations[pos-1].Translate(x)
	//	}
	//}

	return x
}

////////////////////////////////////////////////////////////

type TranslationBySource []Translation

var _ sort.Interface = TranslationBySource(nil)

func (s TranslationBySource) Len() int           { return len(s) }
func (s TranslationBySource) Less(i, j int) bool { return s[i].Start < s[j].Start }
func (s TranslationBySource) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type RangeBySource []Range

var _ sort.Interface = RangeBySource(nil)

func (s RangeBySource) Len() int           { return len(s) }
func (s RangeBySource) Less(i, j int) bool { return s[i].Start < s[j].Start }
func (s RangeBySource) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

////////////////////////////////////////////////////////////
// PART2

func (in *Input) RunProcessPart2() {
	var low int
	n := len(in.Seeds)
	for i := 0; i < n-1; i += 2 {
		v := in.TranslateSeedProcessPart2(in.Seeds[i], in.Seeds[i+1])
		if i == 0 {
			low = lowest(v)
		} else {
			low = min(lowest(v), low)
		}
		// validation that all ranges together make up the same range as we started with
		var sum int
		for _, r := range v {
			sum += r.N
		}
		if sum != in.Seeds[i+1] {
			panic("invalid state")
		}
	}
	fmt.Println("Lowest location:", low)
}

func (in *Input) TranslateSeedProcessPart2(start, num int) []Range {
	t := "seed"
	values := []Range{{start, num}}
	for t != "location" {
		t, values = in.TranslateRanges(t, values)
	}
	return values
}

func lowest(values []Range) int {
	low := values[0].Start
	for _, v := range values {
		low = min(low, v.Start)
	}
	if debug >= 1 {
		fmt.Printf("lowest of %d: %d\n", len(values), low)
	}
	return low
}

func (in *Input) TranslateRanges(t string, inputs []Range) (string, []Range) {
	var rv []Range
	m := in.FindMap(t)
	for _, input := range inputs {
		out := m.TranslateRange(input)
		rv = append(rv, out...)
	}
	return m.To, rv
}

func (m *Map) TranslateRange(input Range) []Range {
	var rv []Range
	var remainder = []Range{input}
	for _, t := range m.Translations {
		if !t.MatchesRange(input) {
			continue
		}
		var newremainder []Range
		for _, x := range remainder {
			result, rem := x.Translate(t)
			rv = append(rv, result)
			newremainder = append(newremainder, rem...)

		}
		remainder = newremainder
	}
	rv = append(rv, remainder...)
	return rv
}

func (t Translation) String() string {
	return fmt.Sprintf("%d-%d:%d", t.Start, t.Max(), t.Dst-t.Start)
}

// translates range 'x', returning the result of that transformation and any remaining, non-translated, ranges.
func (x Range) Translate(t Translation) (Range, []Range) {
	if x.Max() < t.Start || x.Start > t.Max() {
		panic("invalid state")
	}
	var rem []Range
	// the part of x before t
	if x.Start < t.Start && x.Max() >= t.Start {
		before := Range{x.Start, t.Start - x.Start}
		rem = append(rem, before)
	}
	// the part of x after t
	if x.Start <= t.Max() && x.Max() > t.Max() {
		after := Range{t.Max() + 1, x.Max() - t.Max()}
		rem = append(rem, after)
	}
	// the overlapping part
	// first extract the overlapping part
	start := max(x.Start, t.Start)
	end := min(x.Max(), t.Max())
	overlap := Range{start, end - start + 1}
	// then apply the translation as dictated by 't'
	overlap.Start = t.Translate(overlap.Start)
	return overlap, rem
}
