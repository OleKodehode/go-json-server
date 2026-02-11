package service

import (
	"fmt"
	"sort"
	"strings"
)

// sortItems sorts the slice provided (asc by default, "-" prefix for desc)
func sortItems(items []map[string]any, sortStr string) []map[string]any {

	if sortStr == "" {
		return items
	}

	fields := strings.Split(sortStr, ",")

	sort.SliceStable(items, func(i, j int) bool {
		for _, field := range fields {
			field := strings.TrimSpace(field)
			desc := strings.HasPrefix(field, "-")
			field = strings.TrimPrefix(field, "-")

			var a, b any

			// checks for potential missing values
			if value, ok := items[i][field]; ok{ a = value }
			if value, ok := items[j][field]; ok{ b = value }

			// Number comparison first
			aNumb, aErr := toFloat64(a)
			bNumb, bErr := toFloat64(b)

			if aErr == nil && bErr == nil {
				// check if the numbers are the same
				if aNumb == bNumb { continue }

				if desc {
					return aNumb > bNumb
				}
				return aNumb < bNumb
			}
			// Fallback to string

			aStr := fmt.Sprint(a)
			bStr := fmt.Sprint(b)
			// check if the strings are the same
			if aStr == bStr { continue }

			if desc {
				return aStr > bStr
			} 

			return aStr < bStr
		}
		return false // stable if all equal
})

	return items
}