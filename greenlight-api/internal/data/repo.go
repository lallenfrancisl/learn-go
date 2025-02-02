package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Repo struct {
	Movies MovieRepo
}

func NewRepo(db *sql.DB) Repo {
	return Repo{
		Movies: MovieRepo{DB: db},
	}
}
