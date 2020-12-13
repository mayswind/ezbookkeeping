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

func Int64SliceMinus(s1, s2 []int64) []int64 {
	if s1 == nil {
		return nil
	}

	s2ItemsMap := make(map[int64]bool)
	var ret []int64

	for i := 0; i < len(s2); i++ {
		s2ItemsMap[s2[i]] = true
	}

	for i := 0; i < len(s1); i++ {
		if _, exists := s2ItemsMap[s1[i]]; !exists {
			ret = append(ret, s1[i])
		}
	}

	return ret
}

func ToUniqueInt64Slice(items []int64) []int64 {
	var uniqueItems []int64
	itemExistMap := make(map[int64]bool)

	for i := 0; i < len(items); i++ {
		item := items[i]

		if _, exists := itemExistMap[item]; !exists {
			uniqueItems = append(uniqueItems, item)
			itemExistMap[item] = true
		}
	}

	return uniqueItems
}
