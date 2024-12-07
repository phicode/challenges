package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCondition_Split(t *testing.T) {
	gt := Condition{
		Index:  0,
		Op:     '>',
		Value:  5,
		Target: "target",
	}
	lt := Condition{
		Index:  1,
		Op:     '<',
		Value:  5,
		Target: "target",
	}
	r := Range{0, 1000}
	// Range: 0-1000
	// Cond: >5
	// Split: 0-6 ; 6-1000
	ok, nok := gt.Split(r)
	assert.Equal(t, Range{6, 1000}, ok)
	assert.Equal(t, Range{0, 6}, nok)

	// Range: 0-1000
	// Cond: <5
	// Split: 0-5 ; 5-1000
	ok, nok = lt.Split(r)
	assert.Equal(t, Range{0, 5}, ok)
	assert.Equal(t, Range{5, 1000}, nok)
}
