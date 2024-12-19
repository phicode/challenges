package main

import (
	"strings"
	"testing"

	"github.com/phicode/challenges/lib"
	"github.com/stretchr/testify/assert"
)

func TestIsPossible(t *testing.T) {
	design := []byte("bwurrg")
	patterns := lib.Map(strings.Split("r, wr, b, g, bwu, rb, gb, br", ", "), ToBytes)
	isPossible := IsPossible(design, patterns)
	assert.True(t, isPossible)
}
