package lib

import "strings"

func ConcatStrings(lines []string) string {
	b := strings.Builder{}
	for _, l := range lines {
		b.WriteString(l)
	}
	return b.String()
}

func AllStringIndexes(s, substr string) []int {
	var idxs []int
	offset := 0
	for {
		index := strings.Index(s, substr)
		if index == -1 {
			return idxs
		}
		idxs = append(idxs, offset+index)
		s = s[index+1:]
		offset += index + 1
	}
}

func RemoveEmptyLines(lines []string) []string {
  rv := lines[:0]
  for _, l := range lines {
    if len(l) > 0 {
	  rv = append(rv, l)
	}
  }
  return rv
}

