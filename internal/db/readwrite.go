package db

type DB[T any] struct {
	Path string
	Data T
}

func Load(path string) {
	// Check if the file exists
	// If not - Create it with a default value
	// if exists - Read it
	// If emtpy write a default value of {}
	// Unmarshal JSOn into db.Data
	// Return &DB{Path: path, Data: data}
}

func (db *DB[T]) Save() {
	// Marshal db.Data
	// write to db.Path
	// Return error or nil
}