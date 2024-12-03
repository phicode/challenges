package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTestmul(t *testing.T) {
	var tests = []struct {
		input string
		x, y  int
		next  int
		ok    bool
	}{
		{"mul(4)", 0, 0, 0, false},
		{"mul(4,)", 0, 0, 0, false},
		{"mul(,444)", 0, 0, 0, false},
		{"mul(1,2)", 1, 2, 8, true},
		{"mul(1,22)", 1, 22, 9, true},
		{"mul(1,2)mul", 1, 2, 8, true},
		{"mul(1,2]asdf)", 0, 0, 0, false},
	}

	for _, test := range tests {
		test := test
		t.Run(test.input, func(t *testing.T) {
			x, y, next, ok := testmul(test.input)
			assert.Equal(t, test.x, x)
			assert.Equal(t, test.y, y)
			assert.Equal(t, test.next, next)
			assert.Equal(t, test.ok, ok)
		})
	}
}
