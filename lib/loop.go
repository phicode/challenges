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
		max_testlen := (l - i) / 2
		for testlen := 1; testlen <= max_testlen; testlen++ {
			if IsLoop(xs, i, i+testlen) {
				return i, testlen
			}
		}
	}
	return -1, 0
}
