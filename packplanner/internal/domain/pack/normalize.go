package pack

import "sort"

// NormalizePackSizes validates the input, removes duplicates, and sorts sizes ascending.
func NormalizePackSizes(packSizes []int) ([]int, error) {
	if len(packSizes) == 0 {
		return nil, ErrInvalidPackSizes
	}

	unique := make(map[int]struct{}, len(packSizes))
	normalized := make([]int, 0, len(packSizes))

	for _, size := range packSizes {
		if size <= 0 {
			return nil, ErrInvalidPackSizes
		}

		if _, exists := unique[size]; exists {
			continue
		}

		unique[size] = struct{}{}
		normalized = append(normalized, size)
	}

	if len(normalized) == 0 {
		return nil, ErrInvalidPackSizes
	}

	sort.Ints(normalized)
	return normalized, nil
}
