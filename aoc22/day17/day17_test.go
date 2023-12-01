package main

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLineFrom(t *testing.T) {
	line := LineS("..#..")
	assert.Equal(t, Line(0b_00100), line)
	line = LineS("..##......")
	assert.Equal(t, Line(0b_11_000_000), line)
}

func TestRock_Intersects(t *testing.T) {
	lines := []Line{
		LineS("......."),
		LineS("......."),
		LineS("......."),
		LineS("......."),
	}
	slices.Reverse(lines)
	for _, rock := range Rocks {
		assert.False(t, rock.Intersects(lines))
	}

	lines2 := []Line{
		LineS("#######"),
		LineS("......."),
		LineS("......."),
		LineS("......."),
	}
	slices.Reverse(lines2)
	for _, rock := range Rocks {
		assert.True(t, rock.Intersects(lines2))
	}
}
