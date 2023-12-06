package lib

import "strings"

func SplitStringTokens(s string) []string {
	s = strings.TrimSpace(s)
	for strings.Index(s, "  ") != -1 {
		s = strings.ReplaceAll(s, "  ", " ")
	}
	return strings.Split(s, " ")
}
