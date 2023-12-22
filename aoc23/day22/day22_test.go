package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBrick_OnTopOf(t *testing.T) {
	a := ParseBrick("1,0,1~1,2,1")
	b := ParseBrick("0,0,2~2,0,2")
	assert.True(t, b.OnTopOf(a))
}

func TestBrick_OnTopOf_F_and_G(t *testing.T) {
	f := ParseBrick("0,1,6~2,1,6")
	g := ParseBrick("1,1,8~1,1,9")
	assert.False(t, g.OnTopOf(f))

	g.S.Z--
	g.E.Z--
	assert.True(t, g.OnTopOf(f))
}
