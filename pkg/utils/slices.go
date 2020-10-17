package utils

func IsStringSliceEuqals(s1, s2 []string) bool {
	if (s1 == nil) != (s2 == nil) {
		return false
	}

	if len(s1) != len(s2) {
		return false
	}

	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			return false
		}
	}

	return true
}
