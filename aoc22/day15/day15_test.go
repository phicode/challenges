package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitLR(t *testing.T) {
	space := Space{1, 1, 3, 3}
	quads := space.Split()
	assert.Equal(t, Space{1, 1, 1, 3}, quads[0]) // left
	assert.Equal(t, Space{2, 1, 2, 3}, quads[1]) // right
}

func TestSplitTB(t *testing.T) {
	space := Space{1, 1, 2, 3}
	quads := space.Split()
	assert.Equal(t, Space{1, 1, 2, 1}, quads[0]) // top
	assert.Equal(t, Space{1, 2, 2, 2}, quads[1]) // bottom
}
