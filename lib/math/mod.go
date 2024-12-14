package math

// ModUnsigned return the result of number % divisor, as if the integer were unsigned, ie:
// the result will be zero or positive.
func ModUnsigned(number, divisor int) int {
	return (number%divisor + divisor) % divisor
}
