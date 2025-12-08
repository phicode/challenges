package main

// https://adventofcode.com/2025/day/8

import (
	"flag"
	"fmt"
	"math"
	"slices"
	"sort"

	"github.com/phicode/challenges/lib"
)

func main() {
	flag.Parse()
	lib.TimedFunc("Part 1", func() { ProcessPart1("aoc25/day08/example.txt", 10) })
	lib.TimedFunc("Part 1", func() { ProcessPart1("aoc25/day08/input.txt", 1000) })
	//lib.TimedFunc("Part 1", func() { ProcessPart1("aoc25/day08/example.txt", 1000) })
	//
	lib.Timed("Part 2", ProcessPart2, "aoc25/day08/example.txt")
	lib.Timed("Part 2", ProcessPart2, "aoc25/day08/input.txt")

	//lib.Profile(1, "day08.pprof", "Part 2", ProcessPart2, "aoc25/day08/input.txt")
}

func ProcessPart1(name string, loops int) {
	fmt.Println("Part 1 input:", name)
	input := ReadAndParseInput(name)
	result := SolvePart1(input, loops)
	fmt.Println("Result:", result)
}

func ProcessPart2(name string) {
	fmt.Println("Part 2 input:", name)
	input := ReadAndParseInput(name)
	result := SolvePart2(input)
	fmt.Println("Result:", result)
}

////////////////////////////////////////////////////////////

// junction box
type JBox struct {
	X, Y, Z int
	C       *Circuit
}

func (a JBox) Distance(b JBox) float64 {
	x, y, z := a.X-b.X, a.Y-b.Y, a.Z-b.Z
	return math.Sqrt(float64(x*x + y*y + z*z))
}

func (a JBox) String() string { return fmt.Sprintf("(%d,%d,%d)", a.X, a.Y, a.Z) }

type Circuit struct {
	boxes []*JBox
}

type Distance struct {
	dist float64
	a, b *JBox
}

type Input struct {
	boxes    []*JBox
	circuits []*Circuit
}

func (i *Input) Connect(a, b *JBox) {
	ca, cb := a.C, b.C
	founda, foundb := ca != nil, cb != nil
	switch {
	case founda && foundb:
		if ca != cb {
			//fmt.Println("different circuits, merging sizes: ", len(ca.boxes), len(cb.boxes))
			for _, bbox := range b.C.boxes {
				bbox.C = ca
				ca.boxes = append(ca.boxes, bbox)
			}
			cb.boxes = nil
			idx := slices.Index(i.circuits, cb)
			if idx == -1 {
				panic("invalid state")
			}
			i.circuits = RemoveIdx(i.circuits, idx)
		}
	case founda:
		//fmt.Println("adding b to existing circuit of a")
		ca.boxes = append(ca.boxes, b)
		b.C = a.C
	case foundb:
		//fmt.Println("adding a to existing circuit of b")
		cb.boxes = append(cb.boxes, a)
		a.C = b.C
	default:
		//fmt.Println("creating new circuit")
		c := new(Circuit)
		i.circuits = append(i.circuits, c)
		c.boxes = append(c.boxes, a)
		c.boxes = append(c.boxes, b)
		a.C, b.C = c, c
	}
}

func (i *Input) validate() {
	numBoxes := len(i.boxes)
	boxesInCircuits := 0
	for _, c := range i.circuits {
		boxesInCircuits += len(c.boxes)

		for _, b := range c.boxes {
			if b.C != c {
				panic("box in invalid circuit")
			}
		}
	}
	boxesWithoutCircuit := lib.Reduce(i.boxes, 0, func(b *JBox, acc int) int {
		if b.C == nil {
			return acc + 1
		}
		return acc
	})
	if boxesWithoutCircuit+boxesInCircuits != numBoxes {
		panic("box-count missmatch")
	}
}

func ReadAndParseInput(name string) Input {
	lines := lib.ReadLines(name)
	return ParseInput(lines)
}

func ParseInput(lines []string) Input {
	var boxes []*JBox
	for _, line := range lines {
		v := new(JBox)
		n, err := fmt.Sscanf(line, "%d,%d,%d", &v.X, &v.Y, &v.Z)
		if n != 3 || err != nil {
			panic(fmt.Errorf("vector parsing error: n=%d, err=%v", n, err))
		}
		boxes = append(boxes, v)
	}

	return Input{
		boxes: boxes,
	}
}

////////////////////////////////////////////////////////////

func SolvePart1(input Input, loops int) int {
	distances := BuildDistances(input.boxes)
	connected := 0
	for _, dist := range distances {
		//fmt.Printf("connecting: %v, %v d=%f\n", dist.a, dist.b, dist.dist)
		input.Connect(dist.a, dist.b)
		connected++
		//input.validate()
		if connected == loops {
			break
		}
	}
	lens := lib.Map(input.circuits, func(c *Circuit) int { return len(c.boxes) })
	sort.Ints(lens)
	slices.Reverse(lens)

	//fmt.Println("Circuit lenghts:", lens)
	return lens[0] * lens[1] * lens[2]
}

func BuildDistances(boxes []*JBox) []Distance {
	l := len(boxes)
	var rv []Distance
	for i := 0; i < l-1; i++ {
		for j := i + 1; j < l; j++ {
			a, b := boxes[i], boxes[j]
			d := a.Distance(*b)
			rv = append(rv, Distance{d, a, b})
		}
	}
	sort.Slice(rv, func(i, j int) bool {
		return rv[i].dist < rv[j].dist
	})
	return rv
}

func RemoveIdx[T any](ts []T, i int) []T {
	l := len(ts)
	if i < 0 || i >= l {
		panic("index out of range")
	}
	ts[i] = ts[l-1]
	return ts[:l-1]
}

////////////////////////////////////////////////////////////

func SolvePart2(input Input) int {
	distances := BuildDistances(input.boxes)
	connected := 0
	for _, dist := range distances {
		//fmt.Printf("connecting: %v, %v d=%f\n", dist.a, dist.b, dist.dist)
		input.Connect(dist.a, dist.b)
		connected++
		input.validate()
		if len(input.circuits) == 1 && len(input.circuits[0].boxes) == len(input.boxes) {
			//fmt.Println("all boxes connected after", connected, "lines")
			//fmt.Printf("last connected: %v, %v d=%f\n", dist.a, dist.b, dist.dist)
			return dist.a.X * dist.b.X
		}
	}
	return 0
}
