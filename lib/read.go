package lib

import (
	"bufio"
	"os"

	"github.com/phicode/challenges/lib/assets"
)

func ReadLines(name string) []string {
	path := assets.MustFind(name)
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	s := bufio.NewScanner(f)

	var lines []string
	for s.Scan() {
		line := s.Text()
		lines = append(lines, line)
	}

	if err := s.Err(); err != nil {
		panic(s.Err())
	}
	return lines
}
