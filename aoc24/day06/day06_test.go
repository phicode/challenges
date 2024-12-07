package main

import (
	"testing"

	"github.com/phicode/challenges/lib/rowcol"
	"github.com/stretchr/testify/assert"
)

func TestDirSet(t *testing.T) {
	var set DirSet
	assert.False(t, set.IsSet(rowcol.Up))
	assert.False(t, set.IsSet(rowcol.Down))
	assert.False(t, set.IsSet(rowcol.Left))
	assert.False(t, set.IsSet(rowcol.Right))
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
