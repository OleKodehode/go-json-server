package service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/OleKodehode/go-json-server/internal/db"
	"github.com/OleKodehode/go-json-server/internal/model"
)

type Service struct {
	DB *db.DB[model.Data]
}

// Creates a new instance of the Service struct with an attached Database
func New(db *db.DB[model.Data]) *Service {
	return &Service{DB: db}
}

// GET /:name -> Returns all entries within the collection
func (s *Service) GetAll(collection string, filters map[string]string) []map[string]any {
	collection = normalizeInput(collection)

	if !s.collectionExists(collection) {
		return []map[string]any{}
	}
	
	items := s.DB.Data[collection]
	items = applyFilters(items, filters)

	return items
}

// GET /:name/:id -> Returns the requsted entry within a collection if it exists
func (s *Service) GetByID(collection string, id string) map[string]any {
	collection = normalizeInput(collection)

	if !s.collectionExists(collection) {
		return nil
	}

	entry, i := s.findByID(collection, id)

	if i == -1 {
		return nil
	}
	return entry
}

// POST /:name -> Creates a new entry within a collection. Creates a new collection if it doesn't exist
func (s *Service) Create(collection string, item map[string]any) (map[string]any, error) {
	collection = normalizeInput(collection)

	// checks if the collection exists and creates it if it doesn't exist
	s.ensureCollectionExists(collection)
	// add the item to the collection
	if _, ok := item["id"]; !ok {
		item["id"] = generateID(s.DB.Data[collection])
	}

	s.DB.Data[collection] = append(s.DB.Data[collection], item)

	// return the item with the added ID field and whether the DB saved successfully
	return item, s.DB.Save()
}

// PUT /:name/:id -> Replaces (or creates) a specific entry within a collection.
func (s *Service) Replace(collection string, id string, item map[string]any) (map[string]any, error) {
	collection = normalizeInput(collection)
	// Check if the collection exists - Return early if it does not
	if !s.collectionExists(collection) {
		return nil, errors.New("Entry not found - Collection doesn't exist")
	}
	// Check if the entry exists (id) - Return early if it does not
	_, index := s.findByID(collection, id)

	if index == -1 {
		return nil, errors.New("Entry not found - Entry doesn't exist")
	}
	
	item["id"] = id

	s.DB.Data[collection][index] = item
	// Return the updated/created item and whether there were any issues saving the DB
	return item, s.DB.Save()
}

// PATCH /:name/:id -> Updates a specific entry in a collection if it exists
func (s *Service) Update(collection string, id string, fields map[string]any) (map[string]any, error) {
	collection = normalizeInput(collection)
	// Check if the collection exists - Return early if it does not
	if !s.collectionExists(collection) {
		return nil, errors.New("Entry not found - Collection doesn't exist")
	}
	// Check if the entry exists (id) - Return early if it does not
	item, index := s.findByID(collection, id)
	if index == -1 {
		return nil, errors.New("Entry not found - Entry doesn't exist")
	}
	// Update the item with the fields supplied to the function
	for key, value := range fields {
		if key == "id" {
			continue
		}
		item[key] = value
	}

	s.DB.Data[collection][index] = item
	// Return the updated item and whether there were any issues saving the DB
	return item, s.DB.Save()
}

// DELETE /:name/:id -> Deletes a specific entry within a collection if it exists
func (s *Service) Delete(collection string, id string) error {
	collection = normalizeInput(collection)
	// Check if the collection exists - Return early if it does not
	if !s.collectionExists(collection) {
		return errors.New("Entry not found - Collection doesn't exist")
	}
	// check if the entry exists (id) - Return early if it does not
	items := s.DB.Data[collection]

	_, index := s.findByID(collection, id)

	if index == -1 {
		return errors.New("Entry not found - Entry doesn't exist")
	}
	// Delete the entry - Just filter it out based on the ID
	s.DB.Data[collection] = append(items[:index], items[index+1:]...)

	// Return whether there were any errors saving the DB
	return s.DB.Save()
}

// Might need some helper functions?
// Need query filtering

// collectionExists takes in a name of the collection and checks whether it exists
func (s *Service) collectionExists(name string) bool {
	_, ok := s.DB.Data[name]

	return ok
}

// ensureCollectionExists takes in the name of a collection and checks whether it's in the database
// Create a new collection with that name if it doesn't exist
func (s *Service) ensureCollectionExists(name string) {
	if _, ok := s.DB.Data[name]; !ok {
		s.DB.Data[name] = []map[string]any{}
	}
}

// findByID takes in a name of a collection and the ID for the wanted entry.
// Returns that item if it exists and the index of it. Otherwise return nil and -1
func (s *Service) findByID(collection string, id string) (map[string]any, int) {
	items, ok := s.DB.Data[collection]
	if !ok {
		return nil, -1
	}

	for i, item := range items {
		if fmt.Sprint(item["id"]) == id {
			return item, i
		}
	}

	return nil, -1
}

// normalizeInput takes in an input string and removes any white space and makes it lowercase before returning it
func normalizeInput(input string) string {
	normalizedInput := strings.TrimSpace(input)
	normalizedInput = strings.ToLower(normalizedInput)

	return normalizedInput
}

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