package main

import (
	"fmt"
	"testing"
)

func TestResolve(t *testing.T) {
	tests := []struct {
		pos, move, result int
	}{
		{1, 1, 2},
		{90, 9, 99},
		{90, 10, 0},
		{90, -90, 0},
		{90, -91, 99},
		{90, -591, 99},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test-%d", i), func(t *testing.T) {
			got := resolve(test.pos, test.move)
			if got != test.result {
				t.Errorf("pos=%d, move=%d, want-result=%d, got-result=%d", test.pos, test.move, test.result, got)
			}
		})
	}
}
