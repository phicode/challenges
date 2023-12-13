package main

import (
	"testing"

	"git.bind.ch/phil/challenges/lib"
	"github.com/stretchr/testify/assert"
)

func TestPart2(t *testing.T) {
	//lines := lib.ReadLines("aoc23/day13/test_idx_1.txt")
	lines := lib.ReadLines("test_idx_1.txt")
	grids := ParseGrids(lines)
	assert.Equal(t, 1, len(grids))
	g := grids[0]
	ts := g.Transpose()
	// toggling row 11 (index 10) and col 7 (index 6) should find the mirror at line 8 (index 7)
	assert.Equal(t, 7, FindMirrorPart2v2(ts))
}
