package main

// https://adventofcode.com/2022/day/XX

import (
	"bufio"
	"fmt"
	"os"

	"github.com/phicode/challenges/lib/assets"
)

func main() {
	Process("aoc22/dayXX/example.txt")
	Process("aoc22/dayXX/input.txt")
}

func Process(name string) {
	fmt.Println("input:", name)
	lines := ReadInput(name)
	_ = lines

	fmt.Println()
}

func ReadInput(name string) []string {
	f, err := os.Open(assets.MustFind(name))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	s := bufio.NewScanner(f)

	var lines []string
	for s.Scan() {
		line := s.Text()
		lines = append(lines, line)
	}

	if err := s.Err(); err != nil {
		panic(s.Err())
	}
	return lines
}

////////////////////////////////////////////////////////////
