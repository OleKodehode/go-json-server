package service

import (
	"strings"
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

		for rawKey, filterValue := range filters {
			// Parse field + operator (title_contains -> field="title", op="contains")
			field, op := parseFilterKey(rawKey)

			itemValue := item[field]
			if itemValue == nil {
				match = false
				break
			}

			comparator := GetComparator(op)
			if !comparator(itemValue, filterValue) {
				match = false
				break
			}
		}

		if match { result = append(result, item) }
	}

	return result
}

// Expand as needed
var operatorSuffixes = map[string]string{
	"_gte": "gte",
	"_lte": "lte",
	"_gt": "gt",
	"_lt": "lt",
	"_ne": "ne",
	"_contains": "contains",
	"_like": "contains", // alias
}

// parseFilterKey takes in a key (I.E: author_contains/author_like) and returns the field and the operations (author, contains)
func parseFilterKey(key string) (field, op string) {
	for suffix, operator := range operatorSuffixes {
		if trimmed, ok := strings.CutSuffix(key, suffix); ok {
			return trimmed, operator
		}
	}
	// If the loop didn't find any matching operators, return the input and "eq" for equals.
	// Could do lots of guard rails and edge cases, but that feels like a deep rabbit-hole to go down.
	return key, "eq"
}

