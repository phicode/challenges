package lib

import (
	"strconv"
	"strings"
)

func ToInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return n
}

func ExtractInts(s string) []int {
	return Map(strings.Fields(s), ToInt)
}
