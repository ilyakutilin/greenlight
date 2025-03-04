package data

import (
	"database/sql"
	"errors"
)

// A custom ErrRecordNotFound error.
var (
	ErrRecordNotFound = errors.New("record not found")
)

// A Models struct which wraps the MovieModel.
type Models struct {
	Movies MovieModel
}

// A New() method which returns a Models struct containing the initialized MovieModel
// (for ease of use).
func NewModels(db *sql.DB) Models {
	return Models{
		Movies: MovieModel{DB: db},
	}
}
