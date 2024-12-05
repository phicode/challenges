package lib

import (
	"strconv"
	"strings"

	"git.bind.ch/phil/challenges/lib/assert"
)

func ToInt(s string) int {
	n, err := strconv.Atoi(s)
	assert.NoErr(err)
	return n
}

func ExtractInts(s string) []int {
	return Map(strings.Fields(s), ToInt)
}
