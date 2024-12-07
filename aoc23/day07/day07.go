package main

// https://adventofcode.com/2023/day/7

import (
	"fmt"
	"slices"
	"sort"

	"github.com/phicode/challenges/lib"
)

var DEBUG = 0

func main() {
	ProcessPart1("aoc23/day07/example.txt")
	ProcessPart1("aoc23/day07/input.txt")

	//DEBUG = 1
	ProcessPart2("aoc23/day07/example.txt")
	ProcessPart2("aoc23/day07/input.txt")
}

func ProcessPart1(name string) {
	fmt.Println("Part 1 input:", name)
	lines := lib.ReadLines(name)

	hands := ParseHands(lines)
	sort.Sort(HandByStrength(hands))
	var sum int
	for i, h := range hands {
		value := (i + 1) * h.Bid
		sum += value
		debug(1, fmt.Sprintf("Hand %d=%q, strength=%s, value=%d", i+1, string(h.Cards), h.Str, value))
	}
	fmt.Println("Sum:", sum)

	fmt.Println()
}

func ProcessPart2(name string) {
	fmt.Println("Part 2 input:", name)
	lines := lib.ReadLines(name)

	hands := ParseHands(lines)
	sort.Sort(HandByJokerStrength(hands))
	var sum int
	for i, h := range hands {
		value := (i + 1) * h.Bid
		sum += value
		debug(1, fmt.Sprintf("Hand %d=%q, strength=%s, value=%d", i+1, string(h.Cards), h.JokerStr, value))
	}
	fmt.Println("Sum:", sum)

	fmt.Println()
}

func debug(v int, msg string) {
	if v <= DEBUG {
		fmt.Println(msg)
	}
}

////////////////////////////////////////////////////////////

var (
	Cards = []rune{
		'2', '3', '4', '5', '6', '7', '8', '9', 'T', 'J', 'Q', 'K', 'A',
	}
	CardValues = map[rune]int{
		'2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9, 'T': 10, 'J': 11, 'Q': 12, 'K': 13, 'A': 14,
	}
	CardValuesJoker = map[rune]int{
		'2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9, 'T': 10, 'J': 1, 'Q': 12, 'K': 13, 'A': 14,
	}
	JokerIndex = slices.Index(Cards, 'J')
)

type Strength int

const (
	Unknown Strength = iota
	HighCard
	OnePair
	TwoPair
	ThreeOfAKind
	FullHouse
	FourOfAKind
	FiveOfAKind
)

func (s Strength) String() string {
	switch s {
	case Unknown:
		return "Unknown"
	case HighCard:
		return "HighCard"
	case OnePair:
		return "OnePair"
	case TwoPair:
		return "TwoPair"
	case ThreeOfAKind:
		return "ThreeOfAKind"
	case FullHouse:
		return "FullHouse"
	case FourOfAKind:
		return "FourOfAKind"
	case FiveOfAKind:
		return "FiveOfAKind"
	}
	panic("invalid strength value")
}

func ParseHands(lines []string) []Hand {
	var rv []Hand
	for _, l := range lines {
		var h Hand
		var cards string
		n, err := fmt.Sscanf(l, "%s %d", &cards, &h.Bid)
		if n != 2 || err != nil {
			panic(fmt.Errorf("invalid input, n=%d, err=%w", n, err))
		}
		h.Cards = []rune(cards)
		h.Counts = h.CountCards()
		h.Str = h.CalcStrength()
		h.JokerStr = h.CalcJokerStrength()
		rv = append(rv, h)
	}

	return rv
}

type Hand struct {
	Cards    []rune
	Counts   []int
	Bid      int
	Str      Strength
	JokerStr Strength
	Jokers   int
}

func (h Hand) CalcStrength() Strength {
	if CountCounts(h.Counts, 5) >= 1 {
		return FiveOfAKind
	}
	if CountCounts(h.Counts, 4) >= 1 {
		return FourOfAKind
	}
	if CountCounts(h.Counts, 3) >= 1 && CountCounts(h.Counts, 2) >= 1 {
		return FullHouse
	}
	if CountCounts(h.Counts, 3) >= 1 {
		return ThreeOfAKind
	}
	pairs := CountCounts(h.Counts, 2)
	if pairs >= 2 {
		return TwoPair
	}
	if pairs >= 1 {
		return OnePair
	}
	return HighCard
}

func (h Hand) CalcJokerStrength() Strength {
	jokers := h.Counts[JokerIndex]
	if jokers == 0 {
		return h.Str
	}
	if CountCountsNoJoker(h.Counts, 5-jokers) >= 1 {
		return FiveOfAKind
	}
	if CountCountsNoJoker(h.Counts, 4-jokers) >= 1 {
		return FourOfAKind
	}
	// at this point we have at most 2 jokers
	c1 := CountCountsNoJoker(h.Counts, 1)
	c2 := CountCountsNoJoker(h.Counts, 2)

	// full house: 3 + 2 cards
	switch jokers {
	case 1:
		if c2 >= 2 {
			return FullHouse
		}
	case 2:
		if c2 >= 1 && c1 >= 1 {
			return FullHouse
		}
	default:
		panic("invalid number of jokers")
	}

	// ThreeOfAKind
	if CountCountsNoJoker(h.Counts, 3-jokers) >= 1 {
		return ThreeOfAKind
	}

	// at this point we have at most 1 joker

	// TwoPair
	switch jokers {
	case 1:
		if c2 >= 1 && c1 >= 1 {
			return TwoPair
		}
	default:
		panic("invalid number of jokers")
	}

	// OnePair
	return OnePair
}

func CountCounts(counts []int, search int) int {
	var c int
	for _, count := range counts {
		if count == search {
			c++
		}
	}
	return c
}

func CountCountsNoJoker(counts []int, search int) int {
	var c int
	for i, count := range counts {
		if i == JokerIndex {
			continue
		}
		if count == search {
			c++
		}
	}
	return c
}

func (h Hand) CountCards() []int {
	counts := make([]int, len(Cards))
	for i, c := range Cards {
		for _, x := range h.Cards {
			if c == x {
				counts[i]++
			}
		}
	}
	return counts
}

func (h Hand) CardLess(b Hand) bool {
	return h.cardLess(b, CardValues)
}
func (h Hand) CardLessJoker(b Hand) bool {
	return h.cardLess(b, CardValuesJoker)
}

func (h Hand) cardLess(b Hand, values map[rune]int) bool {
	for i, ca := range h.Cards {
		cb := b.Cards[i]
		if ca == cb {
			continue
		}
		va := values[ca]
		vb := values[cb]
		return va < vb
	}
	panic("invalid state: equal strength cards found")
}

////////////////////////////////////////////////////////////

type HandByStrength []Hand

var _ sort.Interface = HandByStrength(nil)

func (h HandByStrength) Len() int      { return len(h) }
func (h HandByStrength) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h HandByStrength) Less(i, j int) bool {
	si := h[i].Str
	sj := h[j].Str
	if si != sj {
		return si < sj
	}
	// equal strength, compare cards in order
	return h[i].CardLess(h[j])
}

////////////////////////////////////////////////////////////

type HandByJokerStrength []Hand

var _ sort.Interface = HandByJokerStrength(nil)

func (h HandByJokerStrength) Len() int      { return len(h) }
func (h HandByJokerStrength) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h HandByJokerStrength) Less(i, j int) bool {
	si := h[i].JokerStr
	sj := h[j].JokerStr
	if si != sj {
		return si < sj
	}
	// equal strength, compare cards in order
	return h[i].CardLessJoker(h[j])
}
