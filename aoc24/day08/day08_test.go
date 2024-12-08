package main

import (
	"fmt"
	"testing"

	"github.com/phicode/challenges/lib/rowcol"
	"github.com/stretchr/testify/assert"
)

func TestAntinodes(t *testing.T) {
	var tests = []struct {
		A, B rowcol.Pos
		P, Q rowcol.Pos
	}{
		{
			rowcol.Pos{4, 4}, rowcol.Pos{6, 5},
			rowcol.Pos{2, 3}, rowcol.Pos{8, 6},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v,%v", test.A, test.B), func(t *testing.T) {
			p, q := Antinodes(test.A, test.B)
			assert.Equal(t, test.P, p)
			assert.Equal(t, test.Q, q)
		})
	}
}
