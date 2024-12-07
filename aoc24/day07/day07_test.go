package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOr(t *testing.T) {
	assert.Equal(t, 11, Or(1, 1))
	assert.Equal(t, 110, Or(1, 10))
	assert.Equal(t, 10, Or(1, 0))
	assert.Equal(t, 12345, Or(12, 345))
	assert.Equal(t, 9999, Or(99, 99))
}

func TestPart2Examples(t *testing.T) {
	equations := []Equation{
		{
			Result: 156, // 15 || 6
			Values: []int{15, 6},
		},
		{
			Result: 7290, // 6 * 8 || 6 * 15
			Values: []int{6, 8, 6, 15},
		},
		{
			Result: 192, // 17 || 8 + 14
			Values: []int{17, 8, 14},
		},
	}
	for _, eq := range equations {
		ok := eq.CanSolve(OpsPart2)
		assert.True(t, ok)
	}
}
