package lib

import (
	"strconv"
)

func ToInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return n
}

func ExtractInts(s string) []int {
	return Map(SplitStringTokens(s), ToInt)
}
