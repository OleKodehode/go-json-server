package db

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type DB[T any] struct {
	Path string
	Data T
}

func Load[T any](file string) (*DB[T], error){
	// Check if the file exists
	path := filepath.Join("data", file + ".json")

	data, err := os.ReadFile(path)
	if err != nil { 
		if os.IsNotExist(err) {
			// If not - Create it with a default value
			emptyDB := []byte("{}")
			if writeErr := os.WriteFile(path, emptyDB, 0644); writeErr != nil {
				return nil, writeErr
			}
			data = emptyDB
		} else {
			return nil, err
		}
	}
	// check if data is empty
	if len(data) <= 0 {
		// If emtpy write a default value of {}
		data = []byte("{}")
		if writeErr := os.WriteFile(path, data, 0644); writeErr != nil {
			return nil, writeErr
		}
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
	// Marshal db.Data
	jsonData, err := json.MarshalIndent(db.Data, "", "  ")
	if err != nil {
		return err
	}
	// write to db.Path
	// Return error or nil
	return os.WriteFile(db.Path, jsonData, 0644)
}