package main

import (
	"testing"

	"git.bind.ch/phil/challenges/lib/rowcol"
	"github.com/stretchr/testify/assert"
)

func TestDirSet_Empty(t *testing.T) {
	var set DirSet
	for _, dir := range rowcol.Directions {
		assert.False(t, set.IsSet(dir))
	}

}
func TestDirSet(t *testing.T) {
	var set DirSet
	set = set.Add(rowcol.Up)
	assert.True(t, set.IsSet(rowcol.Up))
	assert.False(t, set.IsSet(rowcol.Down))
	assert.False(t, set.IsSet(rowcol.Left))
	assert.False(t, set.IsSet(rowcol.Right))
	set = set.Add(rowcol.Down)
	assert.True(t, set.IsSet(rowcol.Up))
	assert.True(t, set.IsSet(rowcol.Down))
	assert.False(t, set.IsSet(rowcol.Left))
	assert.False(t, set.IsSet(rowcol.Right))
	set = set.Add(rowcol.Left)
	assert.True(t, set.IsSet(rowcol.Up))
	assert.True(t, set.IsSet(rowcol.Down))
	assert.True(t, set.IsSet(rowcol.Left))
	assert.False(t, set.IsSet(rowcol.Right))
	set = set.Add(rowcol.Right)
	assert.True(t, set.IsSet(rowcol.Up))
	assert.True(t, set.IsSet(rowcol.Down))
	assert.True(t, set.IsSet(rowcol.Left))
	assert.True(t, set.IsSet(rowcol.Right))
}
