package db

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type DB[T any] struct {
	Path string
	mu sync.RWMutex
	Data T
}

func Load[T any](file string) (*DB[T], error){
	// Check if the file exists
	path := filepath.Join("data", file + ".json")

	// making sure the directory exists
	// 0755 for owner r/w/e, group r/e, others r/e
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) || len(data) == 0 {
		// If not - Create it with a default value
		data = []byte("{}")
		if writeErr := os.WriteFile(path, data, 0644); writeErr != nil {
			return nil, writeErr
		}
	} else if err != nil {
		return nil, err
	}

	// Unmarshal JSON into db.Data
	var dbData  T
	if err := json.Unmarshal(data, &dbData); err != nil {
		return nil, err
	}

	// Return &DB{Path: path, Data: data}
	return &DB[T]{Path: path, Data: dbData}, nil
}

func (db *DB[T]) save() error {
	// Marshal db.Data
	jsonData, err := json.MarshalIndent(db.Data, "", "  ")
	if err != nil {
		return err
	}
	// write to db.Path
	// Return error or nil
	return os.WriteFile(db.Path, jsonData, 0644)
}

// GetCollection returns a copy of the data avilable in the DB
func (db *DB[T]) GetCollection(name string) ([]map[string]any, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	data, ok := any(db.Data).(map[string][]map[string]any)
	if !ok {
		 	return nil, false
	}

	original, exists := data[name]
	if !exists {
		return nil, false
	}

	copyItems := make([]map[string]any, len(original))
	copy(copyItems, original)

	return copyItems, true
}

func (db *DB[T]) UpdateCollection(name string, items []map[string]any) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	data, ok := any(db.Data).(map[string][]map[string]any)
	if !ok {
		return os.ErrInvalid
	}

	data[name] = items

	return db.save()
}