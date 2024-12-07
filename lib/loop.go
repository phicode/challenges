package lib

// start inclusive, end exclusive
func IsLoop[T comparable](xs []T, start, end int) bool {
	l := end - start
	if end+l > len(xs) {
		return false
	}
	for i := 0; i < l; i++ {
		if xs[start+i] != xs[end+i] {
			return false
		}
	}
	return true
}

func FindLoop[T comparable](xs []T) (int, int) {
	l := len(xs)
	for i := 0; i < len(xs); i++ {
		maxTestLen := (l - i) / 2
		for testLen := 1; testLen <= maxTestLen; testLen++ {
			if IsLoop(xs, i, i+testLen) {
				return i, testLen
			}
		}
	}
	return -1, 0
}
