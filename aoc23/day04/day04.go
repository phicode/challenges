package main

// https://adventofcode.com/2023/day/XX

import (
	"fmt"
	"sort"
	"strings"

	"git.bind.ch/phil/challenges/lib"
)

func main() {
	ProcessStep1("aoc23/day04/example.txt")
	ProcessStep1("aoc23/day04/input.txt")

	ProcessStep2("aoc23/day04/example.txt")
	ProcessStep2("aoc23/day04/input.txt")
}

func ProcessStep1(name string) {
	fmt.Println("input:", name)
	lines := lib.ReadLines(name)

	cards := lib.Map(lines, ParseCard)
	var sum int
	for _, c := range cards {
		p := c.Points()
		sum += p
		fmt.Println(c)
		fmt.Println("Points:", p)
	}
	fmt.Println("Total Points:", sum)

	fmt.Println()
}

func ProcessStep2(name string) {
	fmt.Println("input:", name)
	lines := lib.ReadLines(name)

	cards := lib.Map(lines, ParseCard)
	counts := CountCards(cards)
	var sum int
	for _, v := range counts {
		sum += v
	}
	fmt.Println("Counts:", counts)
	fmt.Println("Sum:", sum)

	fmt.Println()
}

////////////////////////////////////////////////////////////

type Card struct {
	Num     int
	Winning []int
	Numbers []int
}

func (c Card) Matches() int {
	var matches int
	l := len(c.Numbers)
	for _, w := range c.Winning {
		if idx := sort.SearchInts(c.Numbers, w); idx < l && c.Numbers[idx] == w {
			matches++
		}
	}
	return matches
}

func (c Card) Points() int {
	m := c.Matches()
	switch m {
	case 0:
		return 0
	case 1:
		return 1
	default:
		return 1 << (m - 1)
	}
}

func ParseCard(l string) Card {
	var c Card
	if n, err := fmt.Sscanf(l, "Card %d:", &c.Num); n != 1 || err != nil {
		panic(fmt.Errorf("parse error, n=%d, err=%v", n, err))
	}
	byColon := strings.Split(l, ":")
	if len(byColon) != 2 {
		panic("invalid input")
	}
	numbers := strings.Split(byColon[1], "|")
	if len(numbers) != 2 {
		panic("invalid input")
	}
	c.Winning = lib.ExtractInts(numbers[0])
	c.Numbers = lib.ExtractInts(numbers[1])
	sort.Ints(c.Numbers)
	return c
}

func (c Card) String() string {
	return fmt.Sprintf("Card %d: %v | %v", c.Num, c.Winning, c.Numbers)
}

func CountCards(cs []Card) []int {
	ns := make([]int, len(cs)+1)
	for i := 1; i < len(ns); i++ {
		ns[i] = 1
	}

	for _, c := range cs {
		m := c.Matches()
		n := ns[c.Num] // how many of these cards we have
		for i := 0; i < m; i++ {
			incrNumber := c.Num + i + 1
			if incrNumber >= len(ns) {
				panic("out of range")
			}
			ns[incrNumber] += n
		}
	}

	return ns
}
