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
func TestResolve2(t *testing.T) {
	tests := []struct {
		pos, move, result, zeros int
	}{
		{1, 1, 2, 0},
		{1, -1, 0, 1},
		{1, -101, 0, 2},
		{1, 99, 0, 1},
		{1, 100, 1, 1},
		{1, 499, 0, 5},
		{1, 500, 1, 5},
		{90, 9, 99, 0},
		{90, 10, 0, 1},
		{90, -90, 0, 1},
		{90, -91, 99, 1},
		{90, -591, 99, 6},
		{0, -100, 0, 1},
		{0, -5, 95, 0},
		{50, 50, 0, 1},
		{50, 51, 1, 1},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test-%d", i), func(t *testing.T) {
			gotResult, gotZeros := resolve2(test.pos, test.move)
			if gotResult != test.result || gotZeros != test.zeros {
				t.Errorf("pos=%d, move=%d, want=(%d, %d), got=(%d, %d)",
					test.pos, test.move,
					test.result, test.zeros,
					gotResult, gotZeros,
				)
			}
		})
	}
}
