package main

// https://adventofcode.com/2023/day/1

import (
	"fmt"
	"strings"

	"git.bind.ch/phil/challenges/lib"
)

func main() {
	ProcessPart1("aoc23/day01/example.txt")
	ProcessPart1("aoc23/day01/input.txt")

	ProcessPart2("aoc23/day01/example.txt")
	ProcessPart2("aoc23/day01/example2.txt")
	ProcessPart2("aoc23/day01/input.txt")
}

func ProcessPart1(name string) {
	fmt.Println("input:", name)
	lines := lib.ReadLines(name)

	var cal int
	for _, l := range lines {
		cal += parse(l)
	}
	fmt.Println("Calibration-1:", cal)

	fmt.Println()
}

func ProcessPart2(name string) {
	fmt.Println("input:", name)
	lines := lib.ReadLines(name)

	fmt.Println("lines:", len(lines))
	var cal int
	for _, l := range lines {
		cal += parse2(l)
	}
	fmt.Println("Calibration-2:", cal)

	fmt.Println()
}

func parse(s string) int {
	firstIdx := strings.IndexFunc(s, IsDigit)
	lastIdx := strings.LastIndexFunc(s, IsDigit)
	if firstIdx == -1 || lastIdx == -1 {
		panic(fmt.Sprintf("less than two digits found: %q", s))
	}
	first, last := s[firstIdx], s[lastIdx]
	first, last = first-'0', last-'0'
	return (int)(first*10 + last)
}

func IsDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func IsDigitI8(r uint8) bool {
	return r >= '0' && r <= '9'
}

func ToDigit(r uint8) int {
	return (int)(r - '0')
}

var digitStrings = map[string]int{
	"one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9,
}

func parse2(s string) int {
	//input := s
	var digits []int
outer:
	for len(s) > 0 {
		if IsDigitI8(s[0]) {
			digits = append(digits, ToDigit(s[0]))
			s = s[1:]
			continue
		}
		for word, value := range digitStrings {
			if strings.HasPrefix(s, word) {
				digits = append(digits, value)
				s = s[1:]
				continue outer
			}
		}
		if len(s) < 1 {
			panic("wtf")
		}
		s = s[1:]
	}
	if len(digits) < 1 {
		panic("not enough digits")
	}
	return digits[0]*10 + digits[len(digits)-1]
	//fmt.Printf("%q = %v = %d\n", input, digits, out)
}
