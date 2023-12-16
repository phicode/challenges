package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCombinations(t *testing.T) {
	// all examples
	assert.Equal(t, 1, run("???.### 1,1,3"))
	assert.Equal(t, 4, run(".??..??...?##. 1,1,3"))
	assert.Equal(t, 1, run("?#?#?#?#?#?#?#? 1,3,1,6"))
	assert.Equal(t, 1, run("????.#...#... 4,1,1"))
	assert.Equal(t, 4, run("????.######..#####. 1,6,5"))
	assert.Equal(t, 10, run("?###???????? 3,2,1"))

	// other trouble makers
	assert.Equal(t, 5, run("??????.??#. 2,3"))
	assert.Equal(t, 15, run(".?.??????#????.?? 1,6"))
}

func run(s string) int {
	seq := ParseSequence(s)
	return SolveCombinations(seq)
}
