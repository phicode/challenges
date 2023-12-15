package main

// https://adventofcode.com/2023/day/XX

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"git.bind.ch/phil/challenges/lib"
)

var VERBOSE = 1

func main() {
	fmt.Println(`Hash of "HASH":`, Hash("HASH"))

	ProcessPart1("aoc23/day15/example.txt")
	ProcessPart1("aoc23/day15/input.txt")

	ProcessPart2("aoc23/day15/example.txt")
	VERBOSE = 0
	ProcessPart2("aoc23/day15/input.txt")
}

func ProcessPart1(name string) {
	fmt.Println("Part 1 input:", name)
	lines := lib.ReadLines(name)

	parts := strings.Split(lines[0], ",")
	var sum int
	for _, p := range parts {
		h := Hash(p)
		fmt.Printf("%q=%d\n", p, h)
		sum += h
	}
	fmt.Println("Sum:", sum)

	fmt.Println()
}

func ProcessPart2(name string) {
	fmt.Println("Part 2 input:", name)
	lines := lib.ReadLines(name)
	SolvePart2(lines[0])

	fmt.Println()
}

func log(v int, msg string) {
	if v <= VERBOSE {
		fmt.Println(msg)
	}
}

////////////////////////////////////////////////////////////

// - Determine the ASCII code for the current character of the string.
// - Increase the current value by the ASCII code you just determined.
// - Set the current value to itself multiplied by 17.
// - Set the current value to the remainder of dividing itself by 256.
func Hash(s string) int {
	var h int
	for _, b := range s {
		h += int(b)
		h *= 17
		h = h & 0xff
	}
	return h
}

type Instruction struct {
	Label     string
	Operation byte
	Lens      int
}

func (i Instruction) String() string {
	if i.Operation == '-' {
		return fmt.Sprintf("%s-", i.Label)
	}
	return fmt.Sprintf("%s=%d", i.Label, i.Lens)
}

func ParseInstructions(s string) []Instruction {
	var rv []Instruction
	parts := strings.Split(s, ",")
	for _, i := range parts {
		rv = append(rv, ParseInstruction(i))
	}
	return rv
}

func ParseInstruction(instr string) Instruction {
	i := 0
	for instr[i] != '-' && instr[i] != '=' {
		i++
	}
	label := instr[:i]
	op := instr[i]
	var lens int
	if op == '=' {
		var err error
		lens, err = strconv.Atoi(instr[i+1:])
		if err != nil {
			panic(err)
		}
	}
	return Instruction{
		Label:     label,
		Operation: op,
		Lens:      lens,
	}
}

type Boxes [256]Box

func (b Boxes) Print() {
	for i, lenses := range b {
		if len(lenses) == 0 {
			continue
		}
		fmt.Printf("Box %d: ", i)
		for i, lens := range lenses {
			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Printf("[%s %d]", lens.Label, lens.Lens)
		}

		fmt.Println()
	}
}

type Box []Lens

type Lens struct {
	Label string
	Lens  int
}

func SolvePart2(s string) {
	instrs := ParseInstructions(s)
	var boxes Boxes
	for _, instr := range instrs {
		RunInstr(&boxes, instr)

		if VERBOSE >= 1 {
			fmt.Printf("After %q\n", instr)
			boxes.Print()
			fmt.Println()
		}
	}
	fmt.Println("Focus Power:", boxes.FocusPower())
}

func RunInstr(boxes *Boxes, instr Instruction) {
	box := Hash(instr.Label)
	if instr.Operation == '-' {
		// If the operation character is a dash (-),
		// go to the relevant box and remove the lens with the given label if it is present in the box.
		// Then, move any remaining lenses as far forward in the box as they can go without changing their order,
		// filling any space made by removing the indicated lens. (If no lens in that box has the given label, nothing happens.)
		for i, lens := range boxes[box] {
			if lens.Label == instr.Label {
				boxes[box] = slices.Delete(boxes[box], i, i+1)
				return
			}
		}
	} else { // operation == '='
		// If there is already a lens in the box with the same label,
		// replace the old lens with the new lens:
		// remove the old lens and put the new lens in its place,
		//not moving any other lenses in the box.
		for i, lens := range boxes[box] {
			if lens.Label == instr.Label {
				boxes[box][i].Lens = instr.Lens
				return
			}
		}

		// If there is not already a lens in the box with the same label,
		// add the lens to the box immediately behind any lenses already in the box.
		// Don't move any of the other lenses when you do this.
		// If there aren't any lenses in the box, the new lens goes all the way to the front of the box.
		boxes[box] = append(boxes[box], Lens{instr.Label, instr.Lens})
	}
}

func (b Boxes) FocusPower() int {
	var sum int
	for i, box := range b {
		a := i + 1
		for slot, lens := range box {
			b := slot + 1
			c := lens.Lens
			focus := a * b * c
			sum += focus
		}
	}
	return sum
}
