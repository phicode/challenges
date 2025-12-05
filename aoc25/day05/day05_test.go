package main

import "testing"

func TestMergeAll(t *testing.T) {
	ranges := []Range{
		{99300079175333, 99300079175333},
		{91024067485337, 99300079175333},
	}
	out := mergeAll(ranges)
	if len(out) != 1 {
		t.Errorf("len(out)=%d; want 1", len(out))
	}
}

func TestMergeAll2(t *testing.T) {
	ranges := []Range{
		{2, 10},
		{4, 5},
		{10, 10},
		{10, 11},
		{1, 5},
	}
	out := mergeAll(ranges)
	if len(out) != 1 {
		t.Errorf("len(out)=%d; want 1", len(out))
	}
}
