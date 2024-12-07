package main

// https://adventofcode.com/2022/day/2

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/phicode/challenges/lib/assets"
)

// https://adventofcode.com/2022/day/2

func main() {
	Process("aoc22/day02/example.txt")
	Process("aoc22/day02/input.txt")
}

func Process(name string) {
	fmt.Println("input:", name)
	rounds := ReadInput(name)
	fmt.Println("rounds:", len(rounds))
	total := 0
	for _, round := range rounds {
		total += round.Score()
	}
	fmt.Println("score:", total)

	fixed := fix(rounds)
	totalFixed := 0
	for _, round := range fixed {
		totalFixed += round.Score()
	}
	fmt.Println("score fixed:", totalFixed)

	fmt.Println()
}

type Move int
type Outcome int

const (
	ROCK Move = iota
	PAPER
	SCISSORS
)

const (
	LOSE Outcome = iota
	DRAW
	WIN
)

func ParseMove(x uint8) Move {
	switch x {
	case 'A', 'X':
		return ROCK
	case 'B', 'Y':
		return PAPER
	case 'C', 'Z':
		return SCISSORS
	}
	panic(fmt.Errorf("invalid move: %v", x))
}

func ParseOutcome(x uint8) Outcome {
	switch x {
	case 'X':
		return LOSE
	case 'Y':
		return DRAW
	case 'Z':
		return WIN
	}
	panic(fmt.Errorf("invalid outcome: %v", x))
}

func (m Move) MoveScore() int {
	switch m {
	case ROCK:
		return 1
	case PAPER:
		return 2
	case SCISSORS:
		return 3
	}
	panic("invalid move")
}

func (m Move) OutcomeAgainst(o Move) int {
	if m == o {
		return 3 // tie
	}
	if m.Wins(o) {
		return 6
	}
	return 0
}

// Wins check rules:
// rock beats scissors
// scissors beats paper
// paper beats rock
func (m Move) Wins(o Move) bool {
	return (m == ROCK && o == SCISSORS) ||
		(m == SCISSORS && o == PAPER) ||
		(m == PAPER && o == ROCK)
}

func (m Move) WinsAgainst() Move {
	switch m {
	case ROCK:
		return SCISSORS
	case SCISSORS:
		return PAPER
	case PAPER:
		return ROCK
	}
	panic("invalid move")
}
func (m Move) LoosesAgainst() Move {
	switch m {
	case ROCK:
		return PAPER
	case PAPER:
		return SCISSORS
	case SCISSORS:
		return ROCK
	}
	panic("invalid move")
}

type Round struct {
	A Move
	B Move
	O Outcome
}

func (r Round) Score() int {
	return r.B.MoveScore() + r.B.OutcomeAgainst(r.A)
}

func ReadInput(name string) []Round {
	f, err := os.Open(assets.MustFind(name))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	var rounds []Round
	for s.Scan() {
		line := s.Text()
		inputs := strings.Split(line, " ")
		if len(inputs) != 2 || len(inputs[0]) != 1 || len(inputs[1]) != 1 {
			panic(fmt.Errorf("invalid input line: %q", line))
		}
		rounds = append(rounds, Round{
			A: ParseMove(inputs[0][0]),
			B: ParseMove(inputs[1][0]),
			O: ParseOutcome(inputs[1][0]),
		})
	}
	if err := s.Err(); err != nil {
		panic(s.Err())
	}
	return rounds
}

func fix(rounds []Round) []Round {
	fixed := make([]Round, len(rounds))
	for i, r := range rounds {
		fixed[i] = r.Fix()
	}
	return fixed
}

func (r Round) Fix() Round {
	switch r.O {
	case LOSE:
		return Round{A: r.A, B: r.A.WinsAgainst(), O: r.O}
	case DRAW:
		return Round{A: r.A, B: r.A, O: r.O}
	case WIN:
		return Round{A: r.A, B: r.A.LoosesAgainst(), O: r.O}
	}
	panic("invalid outcome")
}
