package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRange_Translate(t *testing.T) {
	var tests = []struct {
		R Range
		T Translation
		O []Range
	}{
		{
			R: Range{1, 2},
			T: Translation{Range{3, 2}, 1},
			O: []Range{{1, 2}},
		},

		{
			R: Range{5, 2},
			T: Translation{Range{3, 2}, 1},
			O: []Range{{5, 2}},
		},

		{
			R: Range{1, 5},
			T: Translation{Range{3, 2}, 3},
			O: []Range{
				{1, 2},
				{5, 1},
				{3, 2},
			},
		},

		{
			R: Range{1, 5},
			T: Translation{Range{3, 2}, 4},
			O: []Range{
				{1, 2},
				{5, 1},
				{4, 2},
			},
		},

		{
			R: Range{1, 7},
			T: Translation{Range{6, 5}, 10},
			O: []Range{
				{1, 5},
				{10, 2},
			},
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("test-%d", i), func(t *testing.T) {
			o := test.R.Translate(test.T)
			assert.Equal(t, test.O, o)
		})
	}

}
