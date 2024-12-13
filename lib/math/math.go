package math

// Lcm computes the Least Common Multiply of two values.
func Lcm(a, b int) int {
	return a * (b / Gcd(a, b))
}

// LcmN computes the Least Common Multiply of multiple values.
func LcmN(xs ...int) int {
	l := xs[0]
	for i := 1; i < len(xs); i++ {
		l = (l * xs[i]) / Gcd(l, xs[i])
	}
	return l
}

// Gcd computes the Greatest Common Divisor of two values.
func Gcd(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}
