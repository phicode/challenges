package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNextSecretNumber(t *testing.T) {
	sequence := []int{
		15887950,
		16495136,
		527345,
		704524,
		1553684,
		12683156,
		11100544,
		12249484,
		7753432,
		5908254,
	}
	last := 123
	for _, expected := range sequence {
		next := NextSecretNumber(last)
		assert.Equal(t, expected, next)
		last = next
	}
}

func TestSecretNumber2000(t *testing.T) {
	assert.Equal(t, 8685429, SecretNumber2000(1))
	assert.Equal(t, 4700978, SecretNumber2000(10))
	assert.Equal(t, 15273692, SecretNumber2000(100))
	assert.Equal(t, 8667524, SecretNumber2000(2024))
}

func TestMix(t *testing.T) {
	assert.Equal(t, 37, mix(42, 15))
}
func TestPrune(t *testing.T) {
	assert.Equal(t, 16113920, prune(100000000))
}
