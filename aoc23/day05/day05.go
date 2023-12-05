package main

// https://adventofcode.com/2023/day/XX

import (
	"fmt"
	"sort"
	"strings"

	"git.bind.ch/phil/challenges/lib"
)

var debug = false

func main() {
	ProcessStep1("aoc23/day05/example.txt")
	ProcessStep1("aoc23/day05/input.txt")

	ProcessStep2("aoc23/day05/example.txt")
	ProcessStep2("aoc23/day05/input.txt")
}

func ProcessStep1(name string) {
	fmt.Println("input:", name)
	lines := lib.ReadLines(name)

	input := ParseInput(lines)
	input.Run()

	fmt.Println()
}

func ProcessStep2(name string) {
	fmt.Println("input:", name)
	lines := lib.ReadLines(name)

	input := ParseInput(lines)
	input.RunStep2()

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
		currentMap.Ranges = append(currentMap.Ranges, r)
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
		sort.Sort(SortBySource(m.Ranges))
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

func (in *Input) RunStep2() {
	var lowest int
	n := len(in.Seeds) / 2
	for i := 0; i < n; i += 2 {
		v := in.TranslateSeedStep2(in.Seeds[i], in.Seeds[i+1])
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

func (in *Input) TranslateSeedStep2(start, num int) int {
	t := "seed"
	ranges := []Range{{start, num}}
	var low int
	for t != "location" {
		t, low, ranges = in.TranslateRanges(t, ranges)
	}
	return low
}

func (in *Input) Translate(t string, value int) (string, int) {
	m := in.FindMap(t)
	out := m.Translate(value)
	if debug {
		fmt.Printf("%s -> %s: %d -> %d\n", m.From, m.To, value, out)
	}
	return m.To, out
}

func (in *Input) TranslateRanges(t string, ranges []Range) (string, int, []Range) {
	var outranges []Range
	m := in.FindMap(t)
	var lowest int
	for i, r := range ranges {
		l, out := m.TranslateRange(r)
		outranges = append(outranges, out...)
		if i == 0 {
			lowest = l
		} else {
			lowest = min(lowest, l)
		}
	}
	return m.To, lowest, outranges
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

func (r Range) Overlaps(b Range) bool {
	return !(r.Start > b.Max() || b.Start > r.Max())
}

type Translation struct {
	Range
	Dst int
}

func (r Range) Matches(x int) bool {
	return x >= r.Start && x < r.Start+r.N
}

func (r Translation) Translate(x int) int {
	offset := x - r.Start
	return r.Dst + offset
}

type Map struct {
	From   string
	To     string
	Ranges []Translation // sorted by Source field
}

func (m *Map) Translate(x int) int {
	l := len(m.Ranges)
	if x < m.Ranges[0].Start {
		return x
	}
	last := m.Ranges[l-1]
	if x > last.Max() {
		return x
	}

	for _, r := range m.Ranges {
		if r.Matches(x) {
			return r.Translate(x)
		}
	}

	// binary search version, not really needed for such small input sets
	//searchFunc := func(idx int) bool { return m.Ranges[idx].Src >= x }
	//pos := sort.Search(l, searchFunc)
	//// pos is either an exact match, or one above
	//if pos != l {
	//	if m.Ranges[pos].Matches(x) {
	//		return m.Ranges[pos].Translate(x)
	//	}
	//}
	//if pos > 0 {
	//	if m.Ranges[pos-1].Matches(x) {
	//		return m.Ranges[pos-1].Translate(x)
	//	}
	//}

	return x
}

type TranslateProgress struct {
	Ranges []Range
}

// handle the part of Range `x` that comes before Translation `t`.
func (x Range) Translate(t Translation) []Range {
	if x.Max() < t.Start {
		// the entire range is before the translation range
		// nothing changes
		return []Range{x}
	}
	if x.Start > t.Max() {
		// the entire range is after the translation range
		// nothing changes
		return []Range{x}
	}
	var rv []Range
	// the part of x before t
	if x.Start < t.Start && x.Max() >= t.Start {
		before := Range{x.Start, t.Start - x.Start}
		rv = append(rv, before)
	}
	// the part of x after t
	if x.Start <= t.Max() && x.Max() > t.Max() {
		after := Range{t.Max() + 1, x.Max() - t.Max()}
		rv = append(rv, after)
	}
	// the overlapping part
	// first extract the overlapping part
	// then apply the tranlsation as dictated by 't'
	panic("todo")
}

func (m *Map) TranslateRange(x Range) (int, []Range) {
	l := len(m.Ranges)

	//tp := TranslateProgress{}
	//tp.handleBefore(x, m.Ranges[0])
	//for i, r := range m.Ranges {
	//	tp.handleRange(x, r)
	//	if i+1 < l {
	//		tp.handleBetween(x, r, m.Ranges[i+1])
	//	}
	//}
	//tp.handleAfter(x, m.Ranges[l-1])

	//out := make([]Range, 0, 16)

	// check if range goes beyond left border
	if x.Start < m.Ranges[0].Start {
		if !x.Overlaps(m.Ranges[0]) {
			return
		}
		return x
	}
	//last := m.Ranges[l-1]
	//if x > last.Max() {
	//	return x
	//}
	//
	//for _, r := range m.Ranges {
	//	if r.Matches(x) {
	//		return r.Translate(x)
	//	}
	//}
}

////////////////////////////////////////////////////////////

type SortBySource []Translation

var _ sort.Interface = SortBySource(nil)

func (s SortBySource) Len() int           { return len(s) }
func (s SortBySource) Less(i, j int) bool { return s[i].Start < s[j].Start }
func (s SortBySource) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
