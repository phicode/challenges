package main

// https://adventofcode.com/2025/day/3

import (
	"flag"
	"fmt"

	"github.com/phicode/challenges/lib"
)

func main() {
	flag.Parse()
	lib.Timed("Part 1", ProcessPart1, "aoc25/day03/example.txt")
	lib.Timed("Part 1", ProcessPart1, "aoc25/day03/input.txt")

	lib.Timed("Part 2", ProcessPart2, "aoc25/day03/example.txt")
	lib.Timed("Part 2", ProcessPart2, "aoc25/day03/input.txt")

	//lib.Profile(1, "day03.pprof", "Part 2", ProcessPart2, "aoc25/day03/input.txt")
}

func ProcessPart1(name string) {
	fmt.Println("Part 1 input:", name)
	input := ReadAndParseInput(name)
	result := SolvePart1(input)
	fmt.Println("Result:", result)
}

func ProcessPart2(name string) {
	fmt.Println("Part 2 input:", name)
	input := ReadAndParseInput(name)
	result := SolvePart2(input)
	fmt.Println("Result:", result)
}

////////////////////////////////////////////////////////////

type Bank []int

type Input struct {
	banks []Bank
}

func ReadAndParseInput(name string) Input {
	lines := lib.ReadLines(name)
	return ParseInput(lines)
}

func ParseInput(lines []string) Input {
	var banks []Bank
	for _, line := range lines {
		var b Bank
		//			fmt.Println("line:", line)
		for _, v := range line {
			battery := int(v - '0')
			if battery < 0 || battery > 9 {
				panic(fmt.Errorf("invalid battery value: %c", v))
			}
			b = append(b, battery)
		}
		banks = append(banks, b)
	}
	return Input{banks: banks}
}

////////////////////////////////////////////////////////////

func SolvePart1(input Input) int {
	//	fmt.Println(len(input.banks), "Battery Banks")
	var n int
	for _, bank := range input.banks {
		//		fmt.Printf("bank %v = %d\n", bank, findMax(bank))
		n += findMax(bank)
	}
	return n
}

func findMax(b Bank) int {
	l := len(b)
	first := -1
	firsti := -1
	for i := 0; i < l-1; i++ {
		if v := b[i]; v > first {
			first = v
			firsti = i
		}
	}
	second := b[firsti+1]
	for i := firsti + 2; i < l; i++ {
		if v := b[i]; v > second {
			second = v
		}
	}
	return first*10 + second
}

////////////////////////////////////////////////////////////

func SolvePart2(input Input) int {
	var n int
	for _, bank := range input.banks {
		num := findMax2(bank)
		//		fmt.Printf("bank %v = %d\n", bank, num)
		n += num
		//		fmt.Println("==================================================")
	}
	return n
}

func findMax2(b Bank) int {
	return findMax2Rec(b, 0, 0, 0)
}

func findMax2Rec(b Bank, acc, digit, startIdx int) int {
	keepSpace := 12 - digit - 1
	endIdx := len(b) - keepSpace // exclusive
	//	fmt.Printf("acc=%d, len=%d, digit=%d, start=%d, end=%d\n", acc, len(b),digit, startIdx, endIdx)

	//example:
	// len=16
	// digit=0
	// keepSpace=12-0-1=11
	// endIdx = 16-11 = 5
	// digit0 can be at index 0-4

	maxV := b.maxInRange(startIdx, endIdx)
	if maxV == -1 {
		panic("invalid state")
	}
	maxNum := -1
	//  fmt.Printf("digit=%d=%d\n", digit, maxV)
	for i := startIdx; i < endIdx; i++ {
		if b[i] == maxV {
			//		fmt.Printf("\tchecking idx=%d\n", i)
			num := acc*10 + b[i]
			if digit == 11 {
				maxNum = max(maxNum, num)
			} else {
				num = findMax2Rec(b, num, digit+1, i+1)
				maxNum = max(maxNum, num)
			}
		}
	}
	if maxNum == -1 {
		panic("maxnum=-1")
	}
	return maxNum
}

func (b Bank) maxInRange(start, end int) int {
	m := -1
	for i := start; i < end; i++ {
		m = max(m, b[i])
	}
	return m
}
