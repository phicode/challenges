package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVec3_Distance(t *testing.T) {
	var tests = []struct {
		a, b Vec3
		dist float64
	}{
		{
			a:    Vec3{0, 0, 0},
			b:    Vec3{1, 0, 0},
			dist: 1,
		},
		{
			a:    Vec3{0, 0, 0},
			b:    Vec3{1, 1, 0},
			dist: 1.41421356237309504880, // sqrt(2)
		},
		{
			a:    Vec3{162, 817, 812},
			b:    Vec3{425, 690, 689},
			dist: 316.90219311326957113731,
		},
		{
			a:    Vec3{162, 817, 812},
			b:    Vec3{984, 92, 344},
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

func TestVec3_DistanceDiff(t *testing.T) {
	gota := Vec3{906, 360, 560}
	gotb := Vec3{819, 987, 18}
	wanta := Vec3{431, 825, 988}
	wantb := Vec3{425, 690, 689}
	dx := gota.Distance(gotb)
	dy := wanta.Distance(wantb)
	fmt.Println(dx, dy)
	if dx != dy {
		t.Errorf("want: %v, got: %v", dx, dy)
	}
}
