package main

// https://adventofcode.com/2025/day/2

import (
	"flag"
	"fmt"
	"strings"

	"github.com/phicode/challenges/lib"
)

func main() {
	flag.Parse()
	lib.Timed("Part 1", ProcessPart1, "aoc25/day02/example.txt")
	lib.Timed("Part 1", ProcessPart1, "aoc25/day02/input.txt")

	lib.Timed("Part 2", ProcessPart2, "aoc25/day02/example.txt")
	lib.Timed("Part 2", ProcessPart2, "aoc25/day02/input.txt")

	//lib.Profile(1, "day02.pprof", "Part 2", ProcessPart2, "aoc25/day02/input.txt")
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

type Range struct {
	a, b int
}
type Input struct {
	ranges []Range
}

func ReadAndParseInput(name string) Input {
	lines := lib.ReadLines(name)
	return ParseInput(lines)
}

func ParseInput(lines []string) Input {
	var rv Input
	for _, l := range lines {
		parts := strings.Split(l, ",")
		for _, p := range parts {
			if p == "" {
				continue
			}
			var a, b int
			n, err := fmt.Sscanf(p, "%d-%d", &a, &b)
			if n != 2 || err != nil {
				panic("invalid range: " + p + " err: " + err.Error())
			}
			rv.ranges = append(rv.ranges, Range{a: a, b: b})
		}
	}
	return rv
}

////////////////////////////////////////////////////////////

func SolvePart1(input Input) int {
	//for i, r := range input.ranges {
	//	fmt.Printf("range %d: %d - %d\n", i, r.a, r.b)
	//}
	n := 0
	for _, r := range input.ranges {
		n += findPatterns(r)
	}
	return n
}

func findPatterns(r Range) int {
	alen := log10(r.a)
	blen := log10(r.b)
	//fmt.Println("lengths: ", alen, blen)
	idminlen := alen / 2
	if idminlen*2 < idminlen {
		idminlen++
	}
	idmaxlen := blen / 2
	idmin := 1
	for i := 1; i < idminlen; i++ {
		idmin *= 10
	}
	idmax := 9
	for i := 1; i < idmaxlen; i++ {
		idmax *= 10
		idmax += 9
	}
	//fmt.Println("idmin:", idmin, "idmax:", idmax)

	var n int
	for i := idmin; i <= idmax; i++ {
		num := i
		log := log10(num)
		for l := 0; l < log; l++ {
			num *= 10
		}
		num += i
		//fmt.Println("testing", num)
		if num < r.a {
			continue
		}
		if num > r.b {
			break
		}
		//fmt.Println("INVALID:", num)
		n += num
	}

	return n
}

func log10(x int) int {
	l := 0
	for x > 0 {
		x = x / 10
		l++
	}
	return l
}

////////////////////////////////////////////////////////////

func SolvePart2(input Input) int {
	var n int
	for _, r := range input.ranges {
		n += findRanges2(r)
	}
	return n
}

func findRanges2(r Range) int {
	alen := log10(r.a)
	blen := log10(r.b)
	//fmt.Println("lengths: ", alen, blen)
	idminlen := alen / 2
	if idminlen*2 < idminlen {
		idminlen++
	}
	idmaxlen := blen / 2
	idmin := 1
	for i := 1; i < idminlen; i++ {
		idmin *= 10
	}
	idmax := 9
	for i := 1; i < idmaxlen; i++ {
		idmax *= 10
		idmax += 9
	}
	n := 0
	//fmt.Printf("Range %d - %d; testing %d - %d\n", r.a, r.b, 1, idmax)
	found := make(map[int]bool)
	for i := 1; i <= idmax; i++ {
		ilen := log10(i)
		for rep := 2; rep <= blen; rep++ {
			if ilen*rep > blen {
				//fmt.Printf("abort repeat=%d\n",rep)
				break
			}
			num := repeat(i, ilen, rep)
			//fmt.Printf("\ti=%d*%d=%d\n", i, rep,num)

			if num >= r.a && num <= r.b {
				//			fmt.Println("\tINVALID:", num)
				if !found[num] {
					n += num
				}
				found[num] = true
			}
		}
	}
	return n
}

func repeat(i, ilen, rep int) int {
	num := i
	for j := 1; j < rep; j++ {

		for z := 0; z < ilen; z++ {
			num *= 10
		}
		num += i
	}
	return num
}
