package service

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
)

type Comparator func(itemVal any, filterVal string) bool

var Comparisons = map[string]Comparator { 
	"eq": func(iv any, fv string) bool {	// equal
		return fmt.Sprint(iv) == fv
	},
	"ne": func(iv any, fv string) bool { // Not equal
		return fmt.Sprint(iv) != fv
	},
	"contains": func(iv any, fv string) bool {
		s := fmt.Sprint(iv)
		return strings.Contains(strings.ToLower(s), strings.ToLower(fv))
	},
	"gte": func(iv any, fv string) bool {
		inputNum, filterNum, invalid := numbConvert(iv,fv)
		if invalid {
			return false
		}

		return inputNum >= filterNum
	},
	"lte": func(iv any, fv string) bool {
		inputNum, filterNum, invalid := numbConvert(iv,fv)
		if invalid {
			return false
		}

		return inputNum <= filterNum
	},
	"gt": func(iv any, fv string) bool {
		inputNum, filterNum, invalid := numbConvert(iv,fv)
		if invalid {
			return false
		}

		return inputNum > filterNum
	},
	"lt": func(iv any, fv string) bool {
		inputNum, filterNum, invalid := numbConvert(iv,fv)
		if invalid {
			return false
		}

		return inputNum < filterNum
	},
}

// Function to extract the comparison operation from the Comparisons map. 
func GetComparator(op string) Comparator {
	if op == "" || op == "eq" {
		return Comparisons["eq"]
	}

	if op == "like" {
		op = "contains"
		}

	if comparison, ok := Comparisons[op]; ok {
		return comparison
	}

	// fallback - Invalid input
	return func(_ any, _ string) bool {
		slog.Error("GetComparator failed. ", "Operator:", op)
		return false
	}
}

// numbConvert is a helper function for the Comparison map.
// Returns 2 float64 numbers, and a bool whether the conversion failed.
func numbConvert(iv any, fv string) (float64, float64, bool) {
	inputValue, err := toFloat64(iv)
	if err != nil {
		return 0, 0, true
	}
	filterValue, err := toFloat64(fv)
	if err != nil {
		return 0, 0, true
	}

	return inputValue, filterValue, false
}

// Helper function for consistency instead of using float(64) and strconv
func toFloat64(v any)(float64, error) {
	switch val := v.(type) {
	case float64:
		return val, nil
	case float32:
		return float64(val), nil
	case int:
		return float64(val), nil
	case int64:
		return float64(val), nil
	case string:
		return strconv.ParseFloat(val, 64)
	default:
		return 0, fmt.Errorf("Not a numeric value: %T", v)
	}
}