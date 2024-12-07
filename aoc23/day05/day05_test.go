package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRange_Translate(t *testing.T) {
	var tests = []struct {
		R   Range
		T   Translation
		Out Range
		Rem []Range
	}{

		{
			R:   Range{1, 5},
			T:   Translation{Range{3, 2}, 3},
			Out: Range{3, 2},
			Rem: []Range{
				{1, 2},
				{5, 1},
			},
		},

		{
			R:   Range{1, 5},
			T:   Translation{Range{3, 2}, 4},
			Out: Range{4, 2},
			Rem: []Range{
				{1, 2},
				{5, 1},
			},
		},

		{
			R:   Range{1, 7},
			T:   Translation{Range{6, 5}, 10},
			Out: Range{10, 2},
			Rem: []Range{
				{1, 5},
			},
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("test-%d", i), func(t *testing.T) {
			out, rem := test.R.Translate(test.T)
			assert.Equal(t, test.Out, out)
			assert.Equal(t, test.Rem, rem)
		})
	}

}
