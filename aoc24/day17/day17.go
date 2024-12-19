package main

// https://adventofcode.com/2024/day/17

import (
	"flag"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/phicode/challenges/lib"
	"github.com/phicode/challenges/lib/assert"
)

func main() {
	flag.Parse()
	lib.Timed("Part 1", ProcessPart1, "aoc24/day17/example.txt")
	lib.Timed("Part 1", ProcessPart1, "aoc24/day17/input.txt")

	//lib.Timed("Part 2", ProcessPart2, "aoc24/day17/example2.txt")
	//lib.Timed("Part 2", ProcessPart2, "aoc24/day17/input.txt")

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
	//ExpectEqualProgram bool
	PC int
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
	cpu.Run(false)
	return strings.Join(lib.Map(cpu.Output, strconv.Itoa), ",")
}

// return value: success
func (cpu *CPU) Run(quine bool) bool {
	for {
		cpu.PC++
		//if cpu.IP >= len(cpu.Program) {
		//	panic(fmt.Errorf("pc: %d", cpu.PC))
		//}
		//fmt.Println("pc:", cpu.PC, "instr:", cpu.Program[cpu.IP])
		assert.True(cpu.IP < len(cpu.Program))
		cont := cpu.Tick()
		if !cont {
			return !quine || slices.Equal(cpu.Program, cpu.Output)
		}
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
		//ok := cpu.instrOut()
		//if !ok {
		//	return false, false // halt, without success
		//}
	case InstrBdv:
		cpu.instrBdv()
	case InstrCdv:
		cpu.instrCdv()
	}
	return cpu.IP < len(cpu.Program)
}

// 0-3 = 0-3
// 4 = A
// 5 = B
// 6 = C
func (cpu *CPU) loadCombo(value int) int {
	if value <= 3 {
		return value
	}
	reg := value - 4
	return cpu.Reg[reg]
}

// Division (Opcode 0), store in A
func (cpu *CPU) instrAdv() {
	operand := cpu.Program[cpu.IP+1]
	operand = cpu.loadCombo(operand)
	denom := pow2(operand)
	cpu.Reg[RegA] = cpu.Reg[RegA] / denom
	cpu.IP += 2
}

// Division (Opcode 6), store in B
func (cpu *CPU) instrBdv() {
	operand := cpu.Program[cpu.IP+1]
	operand = cpu.loadCombo(operand)
	denom := pow2(operand)
	cpu.Reg[RegB] = cpu.Reg[RegA] / denom
	cpu.IP += 2
}

// Division (Opcode 7), store in C
func (cpu *CPU) instrCdv() {
	operand := cpu.Program[cpu.IP+1]
	operand = cpu.loadCombo(operand)
	denom := pow2(operand)
	cpu.Reg[RegC] = cpu.Reg[RegA] / denom
	cpu.IP += 2
}

// Bitwise XOR (Opcode 1)
func (cpu *CPU) instrBxl() {
	literal := cpu.Program[cpu.IP+1]
	cpu.Reg[RegB] = cpu.Reg[RegB] ^ literal
	cpu.IP += 2
}

// Modulo 8 (Opcode 2)
// B = combo() % 8
func (cpu *CPU) instrBst() {
	operand := cpu.Program[cpu.IP+1]
	operand = cpu.loadCombo(operand)
	result := operand % 8
	cpu.Reg[RegB] = result
	cpu.IP += 2
}

// Jump if not zero (Opcode 3)
func (cpu *CPU) instrJnz() {
	if cpu.Reg[RegA] == 0 {
		cpu.IP += 2
	} else {
		literal := cpu.Program[cpu.IP+1]
		cpu.IP = literal
	}
}

// Bitwise XOR (Opcode 4) of register B, C
func (cpu *CPU) instrBxc() {
	cpu.Reg[RegB] = cpu.Reg[RegB] ^ cpu.Reg[RegC]
	cpu.IP += 2
}

// Output (Opcode 5): Combo % 8
func (cpu *CPU) instrOut() {
	operand := cpu.Program[cpu.IP+1]
	operand = cpu.loadCombo(operand)
	operand %= 8
	cpu.Output = append(cpu.Output, operand)
	cpu.IP += 2
	//if cpu.ExpectEqualProgram {
	//	i := len(cpu.Output) - 1
	//	return i < len(cpu.Program) && cpu.Program[i] == cpu.Output[i]
	//}
}

////////////////////////////////////////////////////////////

// pow2 calculates 2^x
func pow2(x int) int {
	return 1 << x
}

////////////////////////////////////////////////////////////

func (cpu *CPU) Copy() *CPU {
	cpy := &CPU{
		IP:      cpu.IP,
		Program: make([]int, len(cpu.Program)),
		//ExpectEqualProgram: cpu.ExpectEqualProgram,
		Output: make([]int, len(cpu.Output)),
	}
	copy(cpy.Reg[:], cpu.Reg[:])
	copy(cpy.Program, cpu.Program)
	copy(cpy.Output, cpu.Output)
	return cpy
}

func (cpu *CPU) Reset(cpy *CPU) {
	copy(cpu.Reg[:], cpy.Reg[:])
	cpu.Output = cpu.Output[:0]
	cpu.IP = 0
	cpu.PC = 0
}

func (cpu *CPU) ResetSimple(a int) {
	cpu.Reg[RegA] = a
	cpu.Reg[RegB] = 0
	cpu.Reg[RegC] = 0
	cpu.Output = cpu.Output[:0]
	cpu.IP = 0
	cpu.PC = 0
}

func (cpu *CPU) numcorrect() int {
	l := min(len(cpu.Program), len(cpu.Output))
	for i := 0; i < l; i++ {
		if cpu.Program[i] != cpu.Output[i] {
			return i
		}
	}
	return l
}

func (cpu *CPU) out() string {
	return strings.Join(lib.Map(cpu.Output, strconv.Itoa), ",")
}

func SolvePart2(cpu *CPU) int {
	//incr1 := AnalyseLoops(cpu.Copy(), true)
	//_ = incr1
	//incr2 := AnalyseLoops(cpu.Copy(), false)
	//_ = incr2
	//fmt.Println("incr:", incr1)
	//fmt.Println("incr:", incr2)
	//incr := max(incr1, incr2)

	PrintAll(cpu)
	//PrintAllCorrectTransitions(cpu)
	//PrintAll2(cpu)
	return 0
	//
	//Analyse(cpu.Copy())
	//
	//incr := 1
	//
	//instructions := len(cpu.Program) - 1
	//minA := 1
	//for instructions > 0 {
	//	minA *= 8
	//	instructions--
	//}
	////fmt.Println("min-A:", minA)
	//
	////cpu.ExpectEqualProgram = true
	//cpy := cpu.Copy()
	//a := 3 * minA
	//maxoutlen := 0
	//_ = maxoutlen
	//for {
	//	cpu.Reg[RegA] = a
	//	ok := cpu.Run(true)
	//	//outlen := len(cpu.Output)
	//	//if outlen > maxoutlen {
	//	//	fmt.Println("input-a", a, "registers:", cpu.Reg[RegA], cpu.Reg[RegB], cpu.Reg[RegC], outlen)
	//	//	maxoutlen = outlen
	//	//	//if outlen > 1 && a%incr == 0 && a < 1_000_000 {
	//	//	//	fmt.Println("incr", a)
	//	//	//	incr = a
	//	//	//}
	//	//}
	//	//if outlen > maxoutlen {
	//	//	fmt.Println("a", a, "outlen", outlen)
	//	//	maxoutlen = outlen // expect: 15
	//	//
	//	//}
	//	if ok {
	//		return a
	//	}
	//	cpu.Reset(cpy)
	//	a += incr
	//}
}

func PrintAll(cpu *CPU) {
	minA, _ := MinMax(cpu)
	a := minA - 1
	for {
		cpu.ResetSimple(a)
		if cpu.Run(true) {
			fmt.Println("result:", a)
			return
		}
		fmt.Printf("a: %d, A=%d, B=%d, C=%d, correct=%d, out=%v\n",
			a, cpu.Reg[RegA], cpu.Reg[RegB], cpu.Reg[RegC], cpu.numcorrect(), cpu.out())
		a++
	}
}

func PrintAllCorrectTransitions(cpu *CPU) {
	a := 0
	lastCorrect := -1
	// key: num correct
	// value: A
	lastCorrectA := make(map[int]int)
	for {
		cpu.ResetSimple(a)
		if cpu.Run(true) {
			fmt.Println("result:", a)
			return
		}
		correct := cpu.numcorrect()
		if correct > lastCorrect {
			lca, found := lastCorrectA[correct]
			interval := 0
			if found {
				interval = a - lca
			}
			lastCorrectA[correct] = a
			fmt.Printf("a: %d, A=%d, B=%d, C=%d, correct=%d, interval=%d\n",
				a, cpu.Reg[RegA], cpu.Reg[RegB], cpu.Reg[RegC], correct, interval)
		}
		lastCorrect = correct
		a++
	}
}
func PrintAll2(cpu *CPU) {
	minA, maxA := MinMax(cpu)
	a := minA
	// key: num correct
	// value: A
	lastCorrectA := make(map[int]int)
	bestCorrect := -1
	incr := 1
	searchCorrect := 1
	for {
		cpu.ResetSimple(a)
		if cpu.Run(true) {
			fmt.Println("result:", a)
			return
		}
		correct := cpu.numcorrect()
		//if incr > 1 && correct < lastCorrect {
		//	fmt.Println("asddf")
		//	os.Exit(1)
		//}

		if correct == searchCorrect {
			lca, found := lastCorrectA[correct]
			interval := 0
			if found {
				interval = a - lca
				if correct > bestCorrect {
					bestCorrect = correct
					incr = interval
					searchCorrect++

					fmt.Printf("a: %d, correct=%d, interval=%d, incr=%d\n", a, correct, interval, incr)
					//printNext(cpu, a, 128, incr)
				}
			}
			lastCorrectA[correct] = a
		}
		a += incr
		if a > maxA {
			panic("too many")
		}
	}
}

func MinMax(cpu *CPU) (int, int) {
	instructions := len(cpu.Program) - 1
	minA := 1
	for instructions > 0 {
		minA *= 8
		instructions--
	}
	fmt.Println("min-A:", minA)
	maxA := minA * 8
	fmt.Println("max-A:", maxA)
	combinations := maxA - minA
	fmt.Printf("combinations: 2^%f\n", math.Log2(float64(combinations)))
	return minA, maxA
}

func Analyse(cpu *CPU) {
	minA, maxA := MinMax(cpu)
	_ = maxA
	a := minA
	incr := 1
	pos := 0
	printNext(cpu, a, 32, incr)
	for a < maxA {
		fmt.Println(strings.Repeat("#", 60))
		start, loop, end := FindLoop(cpu, a, pos, incr)
		fmt.Printf("start=%d loop=%d end=%t\n", start, loop, end)

		if end {
			fmt.Println("a:", start)
			return
		}

		if pos == 0 {
			incr *= 8
		}
		incr *= (loop) // * 8)
		a = start
		pos++

		printNext(cpu, start, 32, incr)

		//cpu.Reg[RegA] = a
		//ok := cpu.Run(true)
		//correct := cpu.numcorrect()
		//_ = correct
		//want := cpu.Program[0]
		//got := cpu.Output[0]
		//if got == want {
		//	incr = 64
		//}
		//
		//fmt.Println("a: ", a, "first-result:", got, want)
		//
		////if correct > maxCorrect {
		////	maxCorrect = correct
		////	diff := a - last
		////	last = a
		////	fmt.Println("input-a", a, "registers:", cpu.Reg[RegA], cpu.Reg[RegB], cpu.Reg[RegC], len(cpu.Output), correct, diff)
		////}
		//if ok {
		//	return
		//}
		//cpu.Reset(cpy)
		//a += incr
	}
	panic("no solution found")
}

func printNext(cpu *CPU, startA int, n int, incr int) {
	cpu = cpu.Copy()
	for i := 0; i < n; i++ {
		a := startA + (i * incr)
		cpu.ResetSimple(a)
		cpu.Run(true)
		fmt.Printf("a=%d, A=%d, B=%d, C=%d, output=%d, correct=%d\n",
			a, cpu.Reg[RegA], cpu.Reg[RegB], cpu.Reg[RegC], len(cpu.Output), cpu.numcorrect())
	}

}

func FindStart(cpu *CPU, a int, pos int, incr int) (int, bool) {
	for {
		cpu.ResetSimple(a)
		ok := cpu.Run(true)
		if ok {
			return a, true
		}
		assert.True(len(cpu.Program) == len(cpu.Output))
		if cpu.numcorrect() == pos+1 {
			return a, false
		}
		a += incr
	}
}

// increment a until pos is a match
// increment a until pos does not match any longer
func FindLoop(cpu *CPU, a int, pos, incr int) (int, int, bool) {
	start, end := FindStart(cpu, a, pos, incr)
	if end {
		return start, 0, true
	}
	a = start + incr
	for {
		cpu.ResetSimple(a)
		ok := cpu.Run(true)
		assert.True(len(cpu.Program) == len(cpu.Output))
		want := cpu.Program[pos]
		got := cpu.Output[pos]

		if ok {
			return a, 0, true
		}
		if got != want {
			loop := (a - start) / incr
			return start, loop, false
		}
		a += incr
	}
}
