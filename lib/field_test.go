package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestField(t *testing.T) {
	f := NewField(2000)
	// 2000/64 = 31.25
	assert.Equal(t, 32, len(f.data))
	assert.Equal(t, 2000, f.n)

	for x := 0; x <= 543; x++ {
		f.Set(x)
	}
	x, ok := f.FindMissingField()
	assert.Equal(t, 544, x)
	assert.True(t, ok)
}

func TestFieldLarge(t *testing.T) {
	f := NewField(4_000_001)

	for x := 0; x < 1_000_000; x++ {
		f.Set(x)
	}
	for x := 1_000_001; x <= 4_000_000; x++ {
		f.Set(x)
	}
	x, ok := f.FindMissingField()
	assert.Equal(t, 1_000_000, x)
	assert.True(t, ok)
}
func TestFieldLarge2(t *testing.T) {
	f := NewField(4_000_001)

	for x := 0; x <= 4_000_000; x++ {
		f.Set(x)
	}
	x, ok := f.FindMissingField()
	assert.Equal(t, 0, x)
	assert.False(t, ok)
}
