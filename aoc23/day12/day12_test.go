package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCombinations(t *testing.T) {
	// all examples
	assert.Equal(t, 1, run(t, "???.### 1,1,3"))
	assert.Equal(t, 4, run(t, ".??..??...?##. 1,1,3"))
	assert.Equal(t, 1, run(t, "?#?#?#?#?#?#?#? 1,3,1,6"))
	assert.Equal(t, 1, run(t, "????.#...#... 4,1,1"))
	assert.Equal(t, 4, run(t, "????.######..#####. 1,6,5"))
	assert.Equal(t, 10, run(t, "?###???????? 3,2,1"))

	// other trouble makers
	assert.Equal(t, 5, run(t, "??????.??#. 2,3"))
	assert.Equal(t, 15, run(t, ".?.??????#????.?? 1,6"))
}

func run(t *testing.T, s string) int {
	seq := ParseSequence(s)
	//return MatchRec(seq, 0, 0, 0)
	a, b := SolveCombinations(seq)
	assert.Equal(t, a, b)
	return a
}
