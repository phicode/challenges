package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	Example1 = `
.....0.
..4321.
..5..2.
..6543.
..7..4.
..8765.
..9....`

	Example2 = `
..90..9
...1.98
...2..7
6543456
765.987
876....
987....`

	Example3 = `
012345
123456
234567
345678
4.6789
56789.`

	Example4 = `
89010123
78121874
87430965
96549874
45678903
32019012
01329801
10456732`
)

func TestProcessPart2(t *testing.T) {
	testPart2(t, Example1, 3)
	testPart2(t, Example2, 13)
	testPart2(t, Example3, 227)
	testPart2(t, Example4, 81)
}

func testPart2(t *testing.T, trailMap string, expected int) {
	lines := strings.Split(trailMap, "\n")
	input := ParseInput(lines[1:])
	result := SolvePart2(input)
	assert.Equal(t, expected, result)
}
