package lib

import "strings"

func SplitStringTokens(s string) []string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "  ", " ")
	return strings.Split(s, " ")
}
