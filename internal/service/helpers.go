package service

import (
	"fmt"
	"strconv"
	"strings"
)

// generateID takes in a collection of Items, checks the IDs already present and gets highest number
// Returns the highest number + 1. Simple incrementing ID.
func generateID(items []map[string]any) string {
	max := 0

	// Items could have been deleted, leaving a potential void -> Can't utilize just len()
	for _, item := range items {
		idString := fmt.Sprint(item["id"])
		id, err := strconv.Atoi(idString)
		if err != nil {
			continue
		}
		
		if id > max {
			max = id
		}
	}
	return strconv.Itoa(max + 1)
}

// normalizeInput takes in an input string and removes any white space and makes it lowercase before returning it
func normalizeInput(input string) string {
	normalizedInput := strings.TrimSpace(input)
	normalizedInput = strings.ToLower(normalizedInput)

	return normalizedInput
}

// ensureCollectionExists takes in the name of a collection and checks whether it's in the database
// Create a new collection with that name if it doesn't exist
func (s *Service) ensureCollectionExists(name string) {
	if items, ok := s.DB.GetCollection(name); !ok {
		items = []map[string]any{}
		s.DB.UpdateCollection(name, items)
	}
}

// findByID takes in a slice of items and the ID for the wanted entry.
// Returns that item if it exists and the index of it. Otherwise return nil and -1
func (s *Service) findByID(items []map[string]any, id string) (map[string]any, int) {

	for i, item := range items {
		if fmt.Sprint(item["id"]) == id {
			return item, i
		}
	}

	return nil, -1
}
