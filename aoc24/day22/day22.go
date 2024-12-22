package main

// https://adventofcode.com/2024/day/22

import (
	"flag"
	"fmt"

	"github.com/phicode/challenges/lib"
)

func main() {
	flag.Parse()
	lib.Timed("Part 1", ProcessPart1, "aoc24/day22/example.txt")
	lib.Timed("Part 1", ProcessPart1, "aoc24/day22/input.txt")

	lib.Timed("Part 2", ProcessPart2, "aoc24/day22/example.txt")
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
	for _, n := range input.Numbers {
		if i := HasLoop(n, 2000); i != -1 {
			fmt.Println("Found loop:", n, i)
		}
	}
	return 0
}

func HasLoop(n int, maxIter int) int {
	mem := make(map[int]bool)
	mem[n] = true
	for i := 0; i < maxIter; i++ {
		n = NextSecretNumber(n)
		if mem[n] {
			return i
		}
		mem[n] = true
	}
	return -1
}
