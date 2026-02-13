package collections

import "slices"

func MapKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func SortedMapKeys[K ~string, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	slices.Sort(keys)

	return keys
}

// TruncateSlice truncates a slice to the specified maximum length.
// If the slice is shorter than the maximum length, it is returned unchanged.
func TruncateSlice[T any](s []T, max int) []T {
	if len(s) > max {
		return s[:max]
	}
	return s
}
