package main

// https://adventofcode.com/2022/day/5

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

func main() {
	Process("aoc22/day05/example.txt")
	Process("aoc22/day05/input.txt")
}

func Process(name string) {
	fmt.Println("input:", name)

	stacks, instructions := ReadInput(name)
	fmt.Println(strings.Repeat("*", 40))
	stacks.Print()
	normal := stacks.Copy()
	for _, instr := range instructions {
		normal.Apply(instr)
	}
	fmt.Println(strings.Repeat("*", 40))
	normal.Print()
	fmt.Println(strings.Repeat("*", 40))

	s9001 := stacks.Copy()
	for _, instr := range instructions {
		s9001.Apply9001(instr)
	}
	fmt.Println(strings.Repeat("*", 40))
	s9001.Print()
	fmt.Println(strings.Repeat("*", 40))
	fmt.Println("solution normal", normal.Solution())
	fmt.Println("solution 9001", s9001.Solution())

	fmt.Println()
}

func ReadInput(name string) (CrateStacks, []Instruction) {
	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	s := bufio.NewScanner(f)

	setupStage := true
	var stacks CrateStacks
	var instrs []Instruction
	for s.Scan() {
		line := s.Text()
		if line == "" {
			setupStage = false
			continue
		}
		if setupStage {
			if line[1] != '1' {
				stacks.PushLine(line)
			}
		} else {
			var instr Instruction
			n, err := fmt.Sscanf(line, "move %d from %d to %d", &instr.Amount, &instr.From, &instr.To)
			if err != nil || n != 3 {
				panic(fmt.Errorf("invalid instruction line: %q", line))
			}
			instrs = append(instrs, instr)
		}
	}
	for _, stack := range stacks.stacks {
		stack.Reverse()
	}

	if err := s.Err(); err != nil {
		panic(s.Err())
	}
	return stacks, instrs
}

type CrateStacks struct {
	stacks []stack
}

func (c *CrateStacks) PushLine(line string) {
	i := 0
	for len(line) >= 3 {
		if line[0] == '[' && line[2] == ']' {
			c.getStack(i).Push(rune(line[1]))
		}
		if len(line) == 3 {
			return
		}
		// consume 4 characters
		line = line[4:]
		i++
	}
}

func (c *CrateStacks) getStack(i int) *stack {
	for len(c.stacks) <= i {
		c.stacks = append(c.stacks, stack{})
	}
	return &c.stacks[i]
}

func (c *CrateStacks) Print() {
	rows := c.MaxLen()
	for irow := rows - 1; irow >= 0; irow-- {
		for _, stack := range c.stacks {
			if len(stack) <= irow {
				fmt.Print("    ")
			} else {
				fmt.Printf("[%c] ", stack[irow])
			}
		}
		fmt.Println()
	}
}

func (c *CrateStacks) MaxLen() int {
	m := 0
	for _, s := range c.stacks {
		m = max(m, len(s))
	}
	return m
}

func (c *CrateStacks) Apply(instr Instruction) {
	from := c.getStack(instr.From - 1)
	to := c.getStack(instr.To - 1)
	for i := 0; i < instr.Amount; i++ {
		to.Push(from.Pop())
	}
}
func (c *CrateStacks) Apply9001(instr Instruction) {
	from := c.getStack(instr.From - 1)
	to := c.getStack(instr.To - 1)
	var temp stack
	for i := 0; i < instr.Amount; i++ {
		temp.Push(from.Pop())
	}
	for i := 0; i < instr.Amount; i++ {
		to.Push(temp.Pop())
	}
}

func (c *CrateStacks) Solution() string {
	sol := ""
	for _, s := range c.stacks {
		x := s.Peek()
		if x != ' ' {
			sol += string(x)
		}
	}
	return sol
}

func (c *CrateStacks) Copy() *CrateStacks {
	var stacks []stack
	for _, s := range c.stacks {
		sc := make(stack, len(s))
		copy(sc, s)
		stacks = append(stacks, sc)
	}
	return &CrateStacks{stacks}
}

type stack []rune

func (s *stack) Push(x rune) {
	*s = append(*s, x)
}

func (s *stack) Reverse() {
	slices.Reverse(*s)
}

func (s *stack) Pop() rune {
	last := len(*s) - 1
	r := (*s)[last]
	*s = (*s)[:last]
	return r
}

func (s *stack) Peek() rune {
	l := len(*s)
	if l == 0 {
		return ' '
	}
	return (*s)[l-1]
}

type Instruction struct {
	Amount int
	From   int
	To     int
}
