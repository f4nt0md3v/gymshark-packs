package utilx

func MergeMaps(m1 map[int]int, m2 map[int]int) map[int]int {
	merged := make(map[int]int)

	for k, v := range m1 {
		merged[k] = v
	}

	for key, value := range m2 {
		// Rather than replacing the existing value,
		// add on to any value we already stored.
		merged[key] = merged[key] + value
	}

	return merged
}
