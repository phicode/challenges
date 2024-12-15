package main

import (
	"strings"
	"testing"

	"github.com/phicode/challenges/lib/rowcol"
	"github.com/stretchr/testify/assert"
)

func TestPart2Expand(t *testing.T) {
	const (
		expectedP2 = `
####################
##....[]....[]..[]##
##............[]..##
##..[][]....[]..[]##
##....[]@.....[]..##
##[]##....[]......##
##[]....[]....[]..##
##..[][]..[]..[][]##
##........[]......##
####################`
	)
	lines := strings.Split(expectedP2, "\n")[1:]
	exp := rowcol.NewByteGridFromStrings(lines)
	input := ReadAndParseInput("aoc24/day15/example.txt")
	widened := Part2Expand(input.grid)
	assert.Equal(t, exp, widened)
}
