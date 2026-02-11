package service

import (
	"fmt"
)

// applyFilters takes in a collection of items and filters to apply.
// Returns a collection of items with filters applied.
func applyFilters(items []map[string]any, filters map[string]string) []map[string]any {
	// early return if there are no filters to apply
	if len(filters) == 0 {
		return items
	}

	result := make([]map[string]any, 0, len(items))

	for _, item := range items {
		match := true

		for key, value := range filters {
			if fmt.Sprint(item[key]) != value {
				match = false
				break
			}
		}
		if match {
			result = append(result, item)
		}
	}

	return result
}
