package main

// https://adventofcode.com/2022/day/10

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	Process("aoc22/day10/example.txt")
	Process("aoc22/day10/input.txt")
}

func Process(name string) {
	fmt.Println("input:", name)
	lines := ReadInput(name)

	m := NewMachine(lines)
	sumss := 0
	for i := 0; i < 240; i++ {
		x := m.Tick()
		if (m.Cycle+20)%40 == 0 {
			ss := x * m.Cycle
			sumss += ss
			fmt.Printf("cycle=%d, x=%d, signal-strength=%d\n", m.Cycle, m.X, ss)
		}
	}
	fmt.Println("sum of signal strength:", sumss)
	m.RenderScreen()

	fmt.Println()
}

func ReadInput(name string) []string {
	f, err := os.Open(name)
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

type Machine struct {
	X     int
	Cycle int

	Active   bool
	NextX    int
	EndCycle int

	PC      int
	Program []string

	Finished bool

	Screen [40 * 6]byte
}

func NewMachine(program []string) *Machine {
	return &Machine{
		X:       1,
		Program: program,
	}
}

func (m *Machine) Tick() int {
	if m.Finished {
		return 0
	}
	m.Cycle++
	sample := m.X

	m.updateScreen()

	if m.Active {
		if m.EndCycle == m.Cycle {
			m.X = m.NextX
			m.Active = false
		}
		return sample
	}

	// start next instruction
	if m.PC >= len(m.Program) {
		m.Finished = true
		return sample
	}
	instr := m.Program[m.PC]
	m.PC++
	if instr == "noop" {
		return sample
	}
	var addx int
	n, err := fmt.Sscanf(instr, "addx %d", &addx)
	if n != 1 || err != nil {
		panic("invalid instruction")
	}
	m.EndCycle = m.Cycle + 1
	m.NextX = m.X + addx
	m.Active = true
	return sample
}

func (m *Machine) updateScreen() {
	xpos := (m.Cycle - 1) % 40
	mem := (m.Cycle - 1) % 240
	if xpos-1 <= m.X && m.X <= xpos+1 {
		m.Screen[mem] = '#'
	} else {
		m.Screen[mem] = '.'
	}
}

func (m *Machine) RenderScreen() {
	for i := 0; i < 6; i++ {
		fmt.Println(string(m.Screen[i*40 : (i+1)*40]))
	}
}
