package slice

func SliceIntContains(list []int, target int) bool {
	for _, b := range list {
		if b == target {
			return true
		}
	}
	return false
}

func SliceIntEq(a []int, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, av := range a {
		if av != b[i] {
			return false
		}
	}
	return true
}

func SliceIntLt(a []int, b []int) bool {
	for _, i := range a {
		if !SliceIntContains(b, i) {
			return false
		}
	}
	return true
}
