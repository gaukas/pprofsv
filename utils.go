package pprofsv

func contains(s []int, e int) bool {
	if len(s) == 0 {
		return false
	}

	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
