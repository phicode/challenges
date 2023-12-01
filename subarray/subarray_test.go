package subarray

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubArray(t *testing.T) {
	data := []int{-2, 1, -3, 4, -1, 2, 1, -5, 4}
	expected := []int{4, -1, 2, 1}
	sum, array := SubArray(data)
	assert.Equal(t, expected, array)
	assert.Equal(t, 6, sum)
}

func BenchmarkSubArray10(b *testing.B) {
	benchmarkSubArray(b, 10)
}
func BenchmarkSubArray1000(b *testing.B) {
	benchmarkSubArray(b, 1000)
}
func benchmarkSubArray(b *testing.B, size int) {
	data := make([]int, size)
	for i := 0; i < len(data); i++ {
		data[i] = rand.Int()
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sum, array := SubArray(data)
		assert.NotNil(b, array)
		assert.True(b, sum > 0)
	}
}
