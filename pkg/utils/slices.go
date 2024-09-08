package utils

import "sort"

// Int64SliceEquals returns whether specific two int64 arrays equal
func Int64SliceEquals(s1, s2 []int64) bool {
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

// Int64SliceMinus returns a int64 array which contains items in s1 but not in s2
func Int64SliceMinus(s1, s2 []int64) []int64 {
	if s1 == nil {
		return nil
	}

	s2ItemsMap := make(map[int64]bool)
	ret := make([]int64, 0, len(s1))

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

// ToUniqueInt64Slice returns a int64 array which does not have duplicated items
func ToUniqueInt64Slice(items []int64) []int64 {
	uniqueItems := make([]int64, 0, len(items))
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

// Int64Sort sorts the specified array with int64 items
func Int64Sort(items []int64) {
	sort.Slice(items, func(i, j int) bool {
		return items[i] < items[j]
	})
}

// ToSet returns a map where the keys are the items in the specified array
func ToSet(items []int64) map[int64]bool {
	itemExistMap := make(map[int64]bool)

	for i := 0; i < len(items); i++ {
		item := items[i]
		itemExistMap[item] = true
	}

	return itemExistMap
}
