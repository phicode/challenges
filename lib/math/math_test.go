package math

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGcd(t *testing.T) {
	assert.Equal(t, 2, Gcd(4, 6))
	assert.Equal(t, 1, Gcd(3, 5))
	assert.Equal(t, 1, Gcd(1, 1))
}

func TestLcm(t *testing.T) {
	assert.Equal(t, 12, Lcm(4, 6))
	assert.Equal(t, 15, Lcm(3, 5))
	assert.Equal(t, 1, Lcm(1, 1))
}

func TestLcmN(t *testing.T) {
	assert.Equal(t, 105, LcmN(3, 5, 7))
	assert.Equal(t, 24, LcmN(2, 6, 8))
}
