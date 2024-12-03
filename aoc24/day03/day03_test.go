package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseMul(t *testing.T) {
	var tests = []struct {
		input string
		x, y  int
		rem   string
		ok    bool
	}{
		{"mul(4)", 0, 0, "", false},
		{"mul(4,)", 0, 0, "", false},
		{"mul(,444)", 0, 0, "", false},
		{"mul(1,2)", 1, 2, "", true},
		{"mul(1,2)mul", 1, 2, "mul", true},
		{"mul(1,2]asdf)", 0, 0, "", false},
	}

	for _, test := range tests {
		test := test
		t.Run(test.input, func(t *testing.T) {
			x, y, rem, ok := parseMul(test.input)
			assert.Equal(t, test.x, x)
			assert.Equal(t, test.y, y)
			assert.Equal(t, test.rem, rem)
			assert.Equal(t, test.ok, ok)
		})
	}
}
