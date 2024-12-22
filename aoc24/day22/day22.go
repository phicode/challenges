package main

// https://adventofcode.com/2024/day/22

import (
	"flag"
	"fmt"

	"github.com/phicode/challenges/lib"
	"github.com/phicode/challenges/lib/assert"
)

func main() {
	flag.Parse()
	lib.Timed("Part 1", ProcessPart1, "aoc24/day22/example.txt")
	lib.Timed("Part 1", ProcessPart1, "aoc24/day22/input.txt")

	lib.Timed("Part 2", ProcessPart2, "aoc24/day22/example2.txt")
	lib.Timed("Part 2", ProcessPart2, "aoc24/day22/input.txt")

	//lib.Profile(1, "day22.pprof", "Part 2", ProcessPart2, "aoc24/day22/input.txt")
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

type Input struct {
	Numbers []int
}

func ReadAndParseInput(name string) Input {
	lines := lib.ReadLines(name)
	return ParseInput(lines)
}

func ParseInput(lines []string) Input {
	return Input{Numbers: lib.Map(lines, lib.ToInt)}
}

////////////////////////////////////////////////////////////

func SolvePart1(input Input) int {
	sum := 0
	for _, n := range input.Numbers {
		sum += SecretNumber2000(n)
	}
	return sum
}

func NextSecretNumber(n int) int {
	r := n * 64
	n = mix(n, r)
	n = prune(n)
	r = n / 32
	n = mix(n, r)
	n = prune(n)
	r = n * 2048
	n = mix(n, r)
	n = prune(n)
	return n
}
func mix(n, r int) int {
	return n ^ r
}
func prune(n int) int { return n % 16777216 }

func SecretNumber2000(n int) int {
	for i := 0; i < 2000; i++ {
		n = NextSecretNumber(n)
	}
	return n
}

////////////////////////////////////////////////////////////

func SolvePart2(input Input) int {
	wins := make(map[int]int)
	for _, n := range input.Numbers {
		FindWins(n, wins)
	}
	maxbananas := 0
	maxcomb := 0
	for k, v := range wins {
		if v > maxbananas {
			maxbananas = v
			maxcomb = k
		}
	}
	if lib.LogLevel >= lib.LogDebug {
		PrintCombo(maxcomb)
	}
	return maxbananas
}

func FindWins(n int, wins map[int]int) {
	// changes
	var c1, c2, c3, c4 int
	var prevLastDigit int = n % 10
	usedCombinations := make(map[int]bool)
	for i := 1; i <= 2000; i++ {
		n = NextSecretNumber(n)
		lastDigit := n % 10
		change := lastDigit - prevLastDigit
		c1, c2, c3, c4 = c2, c3, c4, change
		prevLastDigit = lastDigit
		if i >= 4 {
			k := combinationKey(c1, c2, c3, c4)
			if !usedCombinations[k] {
				wins[k] += lastDigit
				usedCombinations[k] = true
			}
		}
	}
}

func combinationKey(c1, c2, c3, c4 int) int {
	assert.True(c1 >= -9 && c1 <= 9)
	assert.True(c2 >= -9 && c2 <= 9)
	assert.True(c3 >= -9 && c3 <= 9)
	assert.True(c4 >= -9 && c4 <= 9)

	c1 += 9 // -9 - +9 => 0-18 -> 5 bits
	c2 += 9
	c3 += 9
	c4 += 9
	// 19^4=130321 combinations
	return c1 + (c2 << 5) + (c3 << 10) + (c4 << 15)
}

func decodeKey(x int) (int, int, int, int) {
	c1 := (x & 0b11111) - 9
	c2 := ((x >> 5) & 0b11111) - 9
	c3 := ((x >> 10) & 0b11111) - 9
	c4 := ((x >> 15) & 0b11111) - 9
	return c1, c2, c3, c4
}

func PrintCombo(comb int) {
	c1, c2, c3, c4 := decodeKey(comb)
	fmt.Println("combination:", c1, c2, c3, c4)
}
