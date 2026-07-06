package utils

// MergeMaps returns a new map that contains all key-value pairs from the provided maps.
// If a key exists in multiple maps, the value from the last map will be used.
func MergeMaps[K comparable, V any](mapsToMerge ...map[K]V) map[K]V {
	m := make(map[K]V)

	for _, src := range mapsToMerge {
		for k, v := range src {
			m[k] = v
		}
	}

	return m
}
