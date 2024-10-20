package utils

func StringSliceContains(source []string, target string) bool {
	for _, s := range source {
		if s == target {
			return true
		}
	}

	return false
}
