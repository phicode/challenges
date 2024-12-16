package main

import (
	"strings"
	"testing"

	"github.com/phicode/challenges/lib/rowcol"
	"github.com/stretchr/testify/assert"
)

const Example1 = `
AAAA
BBCD
BBCC
EEEC`
const Example2 = `
EEEEE
EXXXX
EEEEE
EXXXX
EEEEE`
const Example3 = `
AAAAAA
AAABBA
AAABBA
ABBAAA
ABBAAA
AAAAAA`

func TestPart2(t *testing.T) {
	testPart2(t, Example1, 80)
	testPart2(t, Example2, 236)
	testPart2(t, Example3, 368)
}

func testPart2(t *testing.T, strInput string, expected int) {
	lines := strings.Split(strInput, "\n")
	input := ParseInput(lines[1:])
	result := SolvePart2(input)
	assert.Equal(t, expected, result)
}

func TestNumEdges(t *testing.T) {
	lines := strings.Split(Example1, "\n")
	grid := rowcol.NewByteGridFromStrings(lines[1:])
	assert.Equal(t, 2, NumEdges(grid, rowcol.Pos{0, 0}))
	assert.Equal(t, 0, NumEdges(grid, rowcol.Pos{0, 1}))
	assert.Equal(t, 0, NumEdges(grid, rowcol.Pos{0, 2}))
	assert.Equal(t, 2, NumEdges(grid, rowcol.Pos{0, 3}))
}
