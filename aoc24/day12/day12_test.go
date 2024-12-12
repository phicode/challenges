package main

import (
	"strings"
	"testing"

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
