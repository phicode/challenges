package main

import "testing"

func TestRepeat(t *testing.T) {
	if x := repeat(1, 1, 5); x != 11111 {
		t.Error("1", x)
	}
	if x := repeat(123, 3, 4); x != 123123123123 {
		t.Error("1", x)
	}

}
