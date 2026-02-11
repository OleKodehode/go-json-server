package service

import (
	"fmt"
	"sort"
	"strings"
)

// sortItems sorts the slice by a single field (asc by defautl, desc if _order=desc)
func sortItems(items []map[string]any, field string, order string) []map[string]any {
	desc := strings.ToLower(order) == "desc" || strings.HasPrefix(field, "-")

	// Remove leading "-" if present
	if strings.HasPrefix(field, "-") {
		field = field[1:]
		desc = true
	}

	sort.SliceStable(items, func(i, j int) bool {
		a := items[i][field]
		b := items[j][field]

		// Number comparison first
		aNumb, aErr := toFloat64(a)
		bNumb, bErr := toFloat64(b)

		if aErr != nil && bErr != nil {
			if desc { return aNumb > bNumb}
			return aNumb < bNumb
		}
		// Fallback to string

		aStr := fmt.Sprint(a)
		bStr := fmt.Sprint(b)
		if desc { return aStr > bStr}
		return aStr < bStr
	})

	return items
}