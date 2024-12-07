package main

// https://adventofcode.com/2023/day/2

import (
	"fmt"
	"io"

	"github.com/phicode/challenges/lib"
	"github.com/phicode/challenges/lib/parser"
)

func main() {
	ProcessPart1("aoc23/day02v2/example.txt")
	ProcessPart1("aoc23/day02v2/input.txt")

	ProcessPart2("aoc23/day02v2/example.txt")
	ProcessPart2("aoc23/day02v2/input.txt")
}

func ProcessPart1(name string) {
	fmt.Println("input:", name)
	lines := lib.ReadLines(name)

	games := LoadGames(lines)
	possible := games.FindPossible(12, 13, 14)
	var sum int
	for _, p := range possible {
		sum += p.Nr
	}
	fmt.Printf("possible games found: %d\n", len(possible))
	fmt.Printf("sum of IDs: %d\n", sum)

	fmt.Println()
}

func ProcessPart2(name string) {
	fmt.Println("input:", name)
	lines := lib.ReadLines(name)

	games := LoadGames(lines)
	minimums := games.MinColors()
	var sum int
	for i, m := range minimums {
		p := m.Power()
		_ = i
		//fmt.Printf("power of game %d: %d\n", (i + 1), p)
		sum += p
	}
	fmt.Printf("sum of power: %d\n", sum)

	fmt.Println()
}

////////////////////////////////////////////////////////////

type Games []Game

type Game struct {
	Nr   int
	Sets []Set
}

func (g Game) IsPossible(red, green, blue int) bool {
	for _, s := range g.Sets {
		if !s.IsPossible(red, green, blue) {
			return false
		}
	}
	return true
}

func (g Game) MinSet() Set {
	var s Set
	for _, set := range g.Sets {
		s.Red = max(s.Red, set.Red)
		s.Blue = max(s.Blue, set.Blue)
		s.Green = max(s.Green, set.Green)
	}
	return s
}

type Set struct {
	Red   int
	Green int
	Blue  int
}

func (s Set) IsPossible(red, green, blue int) bool {
	return red >= s.Red && green >= s.Green && blue >= s.Blue
}

func (s Set) Power() int {
	return s.Red * s.Green * s.Blue
}

func LoadGames(lines []string) Games {
	games := Games{}
	for _, l := range lines {
		games = append(games, ParseGame(l))
	}
	return games
}

// Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green
func ParseGame(l string) Game {
	p := parser.New()
	p.SetString(l)
	if !p.AcceptString("Game") {
		panic(p.Error())
	}
	nr, ok := p.Int()
	if !ok {
		panic(p.Error())
	}
	if !p.Accept(':') {
		panic(p.Error())
	}
	return Game{
		Nr:   nr,
		Sets: ParseSets(p),
	}
}

func ParseSets(p *parser.Parser) []Set {
	var rv []Set
	for p.Error() != io.EOF {
		rv = append(rv, ParsetSet(p))
	}
	return rv
}

func ParsetSet(p *parser.Parser) Set {
	var rv Set
	for {
		num, ok := p.Int()
		if !ok {
			panic(p.Error())
		}
		color, ok := p.String()
		if !ok {
			panic(p.Error())
		}
		switch color {
		case "blue":
			rv.Blue = num
		case "red":
			rv.Red = num
		case "green":
			rv.Green = num
		default:
			panic(fmt.Errorf("invalid color: %q", color))
		}
		if next, ok := p.String(); !ok || next == ";" {
			return rv
		}
	}
}

func (x Games) FindPossible(red, green, blue int) []Game {
	var ss []Game
	for _, g := range x {
		if g.IsPossible(red, green, blue) {
			ss = append(ss, g)
		}
	}
	return ss
}

func (x Games) MinColors() []Set {
	var ss []Set
	for _, g := range x {
		ss = append(ss, g.MinSet())
	}
	return ss
}
