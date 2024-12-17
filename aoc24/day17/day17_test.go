package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCPU_loadCombo(t *testing.T) {
	cpu := CPU{
		Reg: [3]int{11, 22, 33},
	}
	assert.Equal(t, 0, cpu.loadCombo(0))
	assert.Equal(t, 1, cpu.loadCombo(1))
	assert.Equal(t, 2, cpu.loadCombo(2))
	assert.Equal(t, 3, cpu.loadCombo(3))
	assert.Equal(t, 11, cpu.loadCombo(4))
	assert.Equal(t, 22, cpu.loadCombo(5))
	assert.Equal(t, 33, cpu.loadCombo(6))

	assert.Panics(t, func() { cpu.loadCombo(7) }, "loadCombo should panic")
}

// If register C contains 9, the program 2,6 would set register B to 1.
func TestCPU_Example1(t *testing.T) {
	cpu := CPU{
		Reg:     [3]int{0, 0, 9},
		Program: []int{2, 6},
	}
	cpu.Run()
	assert.Equal(t, 1, cpu.Reg[RegB])
}

// If register A contains 10, the program 5,0,5,1,5,4 would output 0,1,2.
func TestCPU_Example2(t *testing.T) {
	cpu := CPU{
		Reg:     [3]int{10, 0, 0},
		Program: []int{5, 0, 5, 1, 5, 4},
	}
	cpu.Run()
	assert.Equal(t, []int{0, 1, 2}, cpu.Output)
}

// If register A contains 2024, the program 0,1,5,4,3,0 would output 4,2,5,6,7,7,7,7,3,1,0 and leave 0 in register A.
func TestCPU_Example3(t *testing.T) {
	cpu := CPU{
		Reg:     [3]int{2024, 0, 0},
		Program: []int{0, 1, 5, 4, 3, 0},
	}
	cpu.Run()
	expected := []int{4, 2, 5, 6, 7, 7, 7, 7, 3, 1, 0}
	assert.Equal(t, expected, cpu.Output)
	assert.Equal(t, 0, cpu.Reg[RegA])
}

// If register B contains 29, the program 1,7 would set register B to 26.
func TestCPU_Example4(t *testing.T) {
	cpu := CPU{
		Reg:     [3]int{0, 29, 0},
		Program: []int{1, 7},
	}
	cpu.Run()
	assert.Equal(t, 26, cpu.Reg[RegB])
}

// If register B contains 2024 and register C contains 43690, the program 4,0 would set register B to 44354.
func TestCPU_Example5(t *testing.T) {
	cpu := CPU{
		Reg:     [3]int{0, 2024, 43690},
		Program: []int{4, 0},
	}
	cpu.Run()
	assert.Equal(t, 44354, cpu.Reg[RegB])
}
