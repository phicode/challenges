package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJBox_Distance(t *testing.T) {
	var tests = []struct {
		a, b JBox
		dist float64
	}{
		{
			a:    JBox{0, 0, 0, nil},
			b:    JBox{1, 0, 0, nil},
			dist: 1,
		},
		{
			a:    JBox{0, 0, 0, nil},
			b:    JBox{1, 1, 0, nil},
			dist: 1.41421356237309504880, // sqrt(2)
		},
		{
			a:    JBox{162, 817, 812, nil},
			b:    JBox{425, 690, 689, nil},
			dist: 316.90219311326957113731,
		},
		{
			a:    JBox{162, 817, 812, nil},
			b:    JBox{984, 92, 344, nil},
			dist: 1191.77724428686756160185,
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("test-%d", i), func(t *testing.T) {
			dist := test.a.Distance(test.b)
			assert.InDelta(t, test.dist, dist, 0.001)
		})
	}
}

func TestJBox_DistanceDiff(t *testing.T) {
	gota := JBox{906, 360, 560, nil}
	gotb := JBox{819, 987, 18, nil}
	wanta := JBox{431, 825, 988, nil}
	wantb := JBox{425, 690, 689, nil}
	dx := gota.Distance(gotb)
	dy := wanta.Distance(wantb)
	if dx <= dy {
		t.Errorf("want: %v, got: %v", dx, dy)
	}
}
