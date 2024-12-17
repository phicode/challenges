package main

// https://adventofcode.com/2024/day/17

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/phicode/challenges/lib"
	"github.com/phicode/challenges/lib/assert"
)

func main() {
	flag.Parse()
	lib.Timed("Part 1", ProcessPart1, "aoc24/day17/example.txt")
	lib.Timed("Part 1", ProcessPart1, "aoc24/day17/input.txt")

	lib.Timed("Part 2", ProcessPart2, "aoc24/day17/example.txt")
	lib.Timed("Part 2", ProcessPart2, "aoc24/day17/input.txt")

	//lib.Profile(1, "day17.pprof", "Part 2", ProcessPart2, "aoc24/day17/input.txt")
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

func ReadAndParseInput(name string) *CPU {
	lines := lib.ReadLines(name)
	return ParseInput(lines)
}

type CPU struct {
	Reg     [3]int
	IP      int
	Program []int
	Output  []int
}

const (
	InstrAdv int = iota
	InstrBxl
	InstrBst
	InstrJnz
	InstrBxc
	InstrOut
	InstrBdv
	InstrCdv
)

const (
	RegA = 0
	RegB = 1
	RegC = 2
)

func ParseInput(lines []string) *CPU {
	cpu := &CPU{}
	for _, line := range lines {
		if strings.HasPrefix(line, "Register ") {
			var reg rune
			var val int
			n, err := fmt.Sscanf(line, "Register %c: %d", &reg, &val)
			assert.True(n == 2 && err == nil)
			reg -= 'A'
			cpu.Reg[reg] = val
		}
		if after, found := strings.CutPrefix(line, "Program: "); found {
			split := strings.Split(after, ",")
			cpu.Program = lib.Map(split, lib.ToInt)
		}
	}
	return cpu
}

////////////////////////////////////////////////////////////

func SolvePart1(cpu *CPU) string {
	cpu.Run()
	return strings.Join(lib.Map(cpu.Output, strconv.Itoa), ",")
}

func (cpu *CPU) Run() {
	for cpu.Tick() {
	}
}

func (cpu *CPU) Tick() bool {
	instr := cpu.Program[cpu.IP]
	switch instr {
	case InstrAdv:
		cpu.instrAdv()
	case InstrBxl:
		cpu.instrBxl()
	case InstrBst:
		cpu.instrBst()
	case InstrJnz:
		cpu.instrJnz()
	case InstrBxc:
		cpu.instrBxc()
	case InstrOut:
		cpu.instrOut()
	case InstrBdv:
		cpu.instrBdv()
	case InstrCdv:
		cpu.instrCdv()
	}
	return cpu.IP < len(cpu.Program)
}

func (cpu *CPU) loadCombo(value int) int {
	if value <= 3 {
		return value
	}
	reg := value - 4
	return cpu.Reg[reg]
}

// Division, store in A
func (cpu *CPU) instrAdv() {
	operand := cpu.Program[cpu.IP+1]
	operand = cpu.loadCombo(operand)
	denom := pow2(operand)
	cpu.Reg[RegA] = cpu.Reg[RegA] / denom
	cpu.IP += 2
}

// Division, store in B
func (cpu *CPU) instrBdv() {
	operand := cpu.Program[cpu.IP+1]
	operand = cpu.loadCombo(operand)
	denom := pow2(operand)
	cpu.Reg[RegB] = cpu.Reg[RegA] / denom
	cpu.IP += 2
}

// Division, store in C
func (cpu *CPU) instrCdv() {
	operand := cpu.Program[cpu.IP+1]
	operand = cpu.loadCombo(operand)
	denom := pow2(operand)
	cpu.Reg[RegC] = cpu.Reg[RegA] / denom
	cpu.IP += 2
}

// Bitwise XOR
func (cpu *CPU) instrBxl() {
	literal := cpu.Program[cpu.IP+1]
	cpu.Reg[RegB] = cpu.Reg[RegB] ^ literal
	cpu.IP += 2
}

// Modulo 8
func (cpu *CPU) instrBst() {
	operand := cpu.Program[cpu.IP+1]
	operand = cpu.loadCombo(operand)
	result := operand % 8
	cpu.Reg[RegB] = result
	cpu.IP += 2
}

// Jump if not zero
func (cpu *CPU) instrJnz() {
	if cpu.Reg[RegA] == 0 {
		cpu.IP += 2
	} else {
		literal := cpu.Program[cpu.IP+1]
		cpu.IP = literal
	}
}

// Bitwise XOR of register B, C
func (cpu *CPU) instrBxc() {
	cpu.Reg[RegB] = cpu.Reg[RegB] ^ cpu.Reg[RegC]
	cpu.IP += 2
}

// Output
func (cpu *CPU) instrOut() {
	operand := cpu.Program[cpu.IP+1]
	operand = cpu.loadCombo(operand)
	operand %= 8
	cpu.Output = append(cpu.Output, operand)
	cpu.IP += 2
}

////////////////////////////////////////////////////////////

// pow2 calculates 2^x
func pow2(x int) int {
	// TODO: fast power-of-two
	if x == 0 {
		return 1
	}
	return 2 * pow2(x-1)
}

////////////////////////////////////////////////////////////

func SolvePart2(input *CPU) int {
	return 0
}
