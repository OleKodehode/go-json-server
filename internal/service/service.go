package service

import (
	"errors"
	"maps"
	"strconv"

	"github.com/OleKodehode/go-json-server/internal/db"
	"github.com/OleKodehode/go-json-server/internal/model"
)

type Service struct {
	DB *db.DB[model.Data]
}

var (
	ErrCollectionNotFound = errors.New("Collection not found")
	ErrEntryNotFound = errors.New("Entry not found")
)

// Creates a new instance of the Service struct with an attached Database
func New(db *db.DB[model.Data]) *Service {
	return &Service{DB: db}
}

// GET /:name -> Returns all entries within the collection
func (s *Service) GetAll(collection string, filters map[string]string, control map[string]string) ([]map[string]any, int) {
	collection = normalizeInput(collection)

	if !s.collectionExists(collection) {
		return []map[string]any{}, 0
	}
	
	items := s.DB.Data[collection]
	items = applyFilters(items, filters)
	if sortField, ok := control["_sort"]; ok && sortField != "" {
		items = sortItems(items, sortField, control["_order"])
	}

	total := len(items)

	// Pagination
	page := 1
	perPage := 10 // 10 by default if not supplied

	if reqPage, ok := control["_page"]; ok && reqPage != "" {
		// Check if the request's page is larger than 1
		if n, err := strconv.Atoi(reqPage); err == nil && n >= 1 {
			page = n
		}
	}

	// check if request has a per page and if it's greater than 0
	if reqPerPage, ok := control["_per_page"]; ok && reqPerPage != "" {
		if n, err := strconv.Atoi(reqPerPage); err == nil && n > 0 {
			perPage = n
		}
	} else if limit, ok := control["_limit"]; ok && limit != "" {
		// Legacy Fallback
		if n, err := strconv.Atoi(limit); err == nil && n > 0 {
			perPage = n
		}
	}

	start := (page - 1) * perPage
	if start >= total {
		// return nothing, as the start can't be above the total either way.
		return []map[string]any{}, total
	}

	end := start + perPage
	if end > total {
		end = total
	}

	return items[start:end], total
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

		if err := s.DB.Save(); err != nil {
		return nil, err
	}

	// return the item with the added ID field and whether the DB saved successfully
	return item, nil
}

// PUT /:name/:id -> Replaces (or creates) a specific entry within a collection.
func (s *Service) Replace(collection string, id string, item map[string]any) (map[string]any, error) {
	collection = normalizeInput(collection)
	// Check if the collection exists - Return early if it does not
	if !s.collectionExists(collection) {
		return nil, ErrCollectionNotFound
	}
	
	// Make a copy instead of the original input
	itemCopy := maps.Clone(item)
	// add the ID to the item itself
	itemCopy["id"] = id
	
	// Check if the entry exists (id) - Return early if it does not
	_, index := s.findByID(collection, id)

	if index != -1 {
		s.DB.Data[collection][index] = itemCopy
	} else {
		// entry didn't exist - Create it
		s.DB.Data[collection] = append(s.DB.Data[collection], itemCopy)
	}
	
	if err := s.DB.Save(); err != nil {
		return nil, err
	}
	// Return the updated/created item and whether there were any issues saving the DB
	return itemCopy, nil
}

// PATCH /:name/:id -> Updates a specific entry in a collection if it exists
func (s *Service) Update(collection string, id string, fields map[string]any) (map[string]any, error) {
	collection = normalizeInput(collection)
	// Check if the collection exists - Return early if it does not
	if !s.collectionExists(collection) {
		return nil, ErrCollectionNotFound
	}
	// Check if the entry exists (id) - Return early if it does not
	item, index := s.findByID(collection, id)
	if index == -1 {
		return nil, ErrEntryNotFound
	}
	itemCopy := maps.Clone(item)
	// Update the item with the fields supplied to the function
	for key, value := range fields {
		if key == "id" {
			continue
		}
		itemCopy[key] = value
	}

	s.DB.Data[collection][index] = itemCopy

	if err := s.DB.Save(); err != nil {
		return nil, err
	}

	// Return the updated item and whether there were any issues saving the DB
	return itemCopy, nil
}

// DELETE /:name/:id -> Deletes a specific entry within a collection if it exists
func (s *Service) Delete(collection string, id string) error {
	collection = normalizeInput(collection)
	// Check if the collection exists - Return early if it does not
	if !s.collectionExists(collection) {
		return ErrCollectionNotFound
	}
	// check if the entry exists (id) - Return early if it does not
	items := s.DB.Data[collection]

	_, index := s.findByID(collection, id)

	if index == -1 {
		return ErrEntryNotFound
	}
	// Delete the entry - Just filter it out based on the ID
	s.DB.Data[collection] = append(items[:index], items[index+1:]...)

	if err := s.DB.Save(); err != nil {
		return err
	}

	return nil
}
