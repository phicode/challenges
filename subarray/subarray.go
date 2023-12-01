package subarray

func SubArray(ts []int) (int, []int) {
	current, best := 0, 0
	endIdx := -1
	for i, t := range ts {
		current = max(0, current+t)
		if current > best {
			best = current
			endIdx = i
		}
	}
	startIdx := endIdx
	rem := best - ts[endIdx]
	for rem > 0 {
		startIdx--
		rem -= ts[startIdx]
	}
	return best, ts[startIdx : endIdx+1]
}
