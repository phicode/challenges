package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindLoop(t *testing.T) {
	{
		test := []int{1, 2, 3, 3, 4, 5}
		start, l := FindLoop(test)
		assert.Equal(t, 2, start)
		assert.Equal(t, 1, l)
	}
	{
		test := []int{0, 1, 2, 3, 1, 2, 3}
		start, l := FindLoop(test)
		assert.Equal(t, 1, start)
		assert.Equal(t, 3, l)
	}
}
