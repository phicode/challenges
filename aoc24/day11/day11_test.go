package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNext(t *testing.T) {
	testNext(t, 0, 1, -1)
	testNext(t, 1, 2024, -1)
	testNext(t, 2024, 20, 24)
	testNext(t, 125, 253000, -1)
	testNext(t, 253000, 253, 0)
}

func testNext(t *testing.T, num, a, b int) {
	gotA, gotB := next(num)
	assert.Equal(t, a, gotA)
	assert.Equal(t, b, gotB)
}

func TestSolve(t *testing.T) {
	cache := make(map[Key]int)
	stones := solve(cache, 125, 1, 1) // 1 blink
	// 125 => 253000
	// = 1 stone
	assert.Equal(t, 1, stones)

	stones = solve(cache, 125, 1, 2) // 2 blinks
	// 253000 => 253 , 0
	// = 2 stones
	assert.Equal(t, 2, stones)
	stones = solve(cache, 125, 1, 2)
	assert.Equal(t, 2, stones)

	stones = solve(cache, 125, 1, 3) // 3 blinks
	// 125 => 253000
	// 253000 => 253 , 0
	// 253, 0 => 512072, 1
	// = 2 stones
	assert.Equal(t, 2, stones)

	stones = solve(cache, 125, 1, 4) // 4 blinks
	// 125 => 253000
	// 253000 => 253 , 0
	// 253, 0 => 512072, 1
	// 512072, 1 => 512, 72, 2024
	// = 3 stones
	assert.Equal(t, 3, stones)
}

func TestExample(t *testing.T) {
	in := ParseInput("125 17")
	assert.Equal(t, 3, Solve(in, 1))
	assert.Equal(t, 55312, Solve(in, 25))
}
