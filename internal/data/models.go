package data

import (
	"database/sql"
	"errors"
)

// A custom ErrRecordNotFound error.
var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Movies MovieModel
	Tokens TokenModel
	Users  UserModel
}

// A New() method which returns a Models struct containing the initialized MovieModel
// (for ease of use).
func NewModels(db *sql.DB) Models {
	return Models{
		Movies: MovieModel{DB: db},
		Tokens: TokenModel{DB: db},
		Users:  UserModel{DB: db},
	}
}
