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

func (db *DB[T]) Save() error {
	db.Lock()
	defer db.Unlock()
	// Marshal db.Data
	jsonData, err := json.MarshalIndent(db.Data, "", "  ")
	if err != nil {
		return err
	}
	// write to db.Path
	// Return error or nil
	return os.WriteFile(db.Path, jsonData, 0644)
}

// Lock locks the database during writes using the Mutex's Lock method
func (db *DB[T]) Lock() { db.mu.Lock() }
// unlock unlocks the database after writing using the Mutex's Unlock method
func (db *DB[T]) Unlock() { db.mu.Unlock() }

// Readlock locks the database during reads using the Mutex's RLock method
func (db *DB[T]) ReadLock() { db.mu.RLock() }
// ReadUnlock unlocks the database after reading using the Mutex's RUnlock method
func (db *DB[T]) ReadUnlock() { db.mu.RUnlock() }