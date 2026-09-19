package support

import (
	"cmp"
	"slices"
)

func SortedMapKeys[TKey cmp.Ordered, TValue any](m map[TKey]TValue) []TKey {
	keys := MapKeys(m)

	slices.Sort(keys)

	return keys
}

func MapKeys[TKey comparable, TValue any](m map[TKey]TValue) []TKey {
	keys := make([]TKey, 0, len(m))

	for key := range m {
		keys = append(keys, key)
	}

	return keys
}
