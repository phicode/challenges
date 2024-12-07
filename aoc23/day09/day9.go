package main

// https://adventofcode.com/2023/day/9

import (
	"fmt"

	"github.com/phicode/challenges/lib"
)

var VERBOSE = 0

func main() {
	// 18 + 28 + 68 = 114
	ProcessPart1("aoc23/day09/example.txt")
	ProcessPart1("aoc23/day09/input.txt")

	// -3 + 0 + 5 = 2
	ProcessPart2("aoc23/day09/example.txt")
	ProcessPart2("aoc23/day09/input.txt")
}

func ProcessPart1(name string) {
	fmt.Println("Part 1 input:", name)
	lines := lib.ReadLines(name)
	inputs := ParseLines(lines)
	var sum int
	for i, input := range inputs {
		next := FindNextNumber(input)
		if VERBOSE >= 1 {
			fmt.Println("line", i+1, "next number:", next)
		}
		sum += next
	}
	fmt.Println("Sum:", sum)
	fmt.Println()
}

func ProcessPart2(name string) {
	fmt.Println("Part 2 input:", name)
	lines := lib.ReadLines(name)
	inputs := ParseLines(lines)
	var sum int
	for i, input := range inputs {
		next := FindPreviousNumber(input)
		if VERBOSE >= 1 {
			fmt.Println("line", i+1, "previous number:", next)
		}
		sum += next
	}
	fmt.Println("Sum:", sum)

	fmt.Println()
}

////////////////////////////////////////////////////////////

func ParseLines(lines []string) [][]int {
	return lib.Map(lines, lib.ExtractInts)
}

func FindNextNumber(xs []int) int {
	derivatives := Derivatives(xs)
	d := 0
	for i := len(derivatives) - 1; i >= 0; i-- {
		d += derivatives[i][len(derivatives[i])-1]
	}
	return xs[len(xs)-1] + d
}
func FindPreviousNumber(xs []int) int {
	derivatives := Derivatives(xs)
	d := 0
	for i := len(derivatives) - 1; i >= 0; i-- {
		d = derivatives[i][0] - d
	}
	return xs[0] - d
}

func Derivatives(xs []int) [][]int {
	var derivatives [][]int
	current := xs
	for {
		n := len(current)
		d := make([]int, n-1)
		nZero := 0
		for i := 0; i < n-1; i++ {
			d[i] = current[i+1] - current[i]
			if d[i] == 0 {
				nZero++
			}
		}
		if nZero == n-1 {
			break
		}
		current = d
		derivatives = append(derivatives, d)
	}
	return derivatives
}
