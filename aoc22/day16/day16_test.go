package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValve_ValueAtTime(t *testing.T) {
	v := Valve{Name: "AA", Flow: 5}
	assert.Equal(t, 0, v.ValueAtTime(0))
	assert.Equal(t, 0, v.ValueAtTime(1))
	assert.Equal(t, 5, v.ValueAtTime(2))
	assert.Equal(t, 10, v.ValueAtTime(3))
}

func TestActivatedValve(t *testing.T) {
	var av OpenValves
	av.SetOpen(5)

	for i := 0; i < 64; i++ {
		assert.Equal(t, i == 5, av.IsOpen(i))
	}
}
