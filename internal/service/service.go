package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/OleKodehode/go-json-server/internal/db"
	"github.com/OleKodehode/go-json-server/internal/model"
)

type Service struct {
	DB *db.DB[model.Data]
}

func New(db *db.DB[model.Data]) *Service {
	return &Service{DB: db}
}

// GET /:name
func (s *Service) GetAll(collection string, filters map[string]string) []map[string]any {
	collection = normalizeInput(collection)

	if !s.collectionExists(collection) {
		return []map[string]any{}
	}
	// Check if there are any filters - Apply

	return s.DB.Data[collection]
}

// GET /:name/:id
func (s *Service) GetByID(collection string, id string) map[string]any {
	collection = normalizeInput(collection)

	if !s.collectionExists(collection) {
		return map[string]any{}
	}

	entry, i := s.findByID(collection, id)

	if i == -1 {
		return map[string]any{}
	}
	return entry
}

// POST /:name
func (s *Service) Create(collection string, item map[string]any) (map[string]any, error) {
	collection = normalizeInput(collection)

	// checks if the collection exists and creates it if it doesn't exist
	s.ensureCollectionExists(collection)
	// add the item to the collection
	return map[string]any{}, s.DB.Save()
}

// PUT /:name/:id
func (s *Service) Replace(collection string, id string, item map[string]any) (map[string]any, error) {
	collection = normalizeInput(collection)
	// Check if the collection exists - Return early if it does not
	if !s.collectionExists(collection) {
		return nil, errors.New("Entry not found")
	}
	// Check if the entry exists (id) - Return early if it does not
	// Replace the entry with the new item
	// Return nil (successful PUT)
	return map[string]any{}, s.DB.Save()
}

// PATCH /:name/:id
func (s *Service) Update(collection string, id string, fields map[string]any) (map[string]any, error) {
	collection = normalizeInput(collection)
	// Check if the collection exists - Return early if it does not
	if !s.collectionExists(collection) {
		return nil, errors.New("Entry not found")
	}
	// Check if the entry exists (id) - Return early if it does not
	// Update the item with the fields supplied to the function
	// return nil (successfull PATCH)
	return map[string]any{}, s.DB.Save()
}

// DELETE /:name/:id
func (s *Service) Delete(collection string, id string) error {
	collection = normalizeInput(collection)
	// Check if the collection exists - Return early if it does not
	if !s.collectionExists(collection) {
		return errors.New("Entry not found")
	}
	// check if the entry exists (id) - Return early if it does not
	// Delete the entry - Just filter it out based on the ID
	// Return nil (Successful DELETE)
	return s.DB.Save()
}

// Might need some helper functions?
// Need query filtering

// collectionExists takes in a name of the collection and checks whether it exists
func (s *Service) collectionExists(name string) bool {
	_, ok := s.DB.Data[name]

	return ok
}

func (s *Service) ensureCollectionExists(name string) {
	if _, ok := s.DB.Data[name]; !ok {
		s.DB.Data[name] = []map[string]any{}
	}
}

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
