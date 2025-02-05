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
	Users  UserRepo
	Tokens TokenRepo
}

func NewRepo(db *sql.DB) Repo {
	return Repo{
		Movies: MovieRepo{DB: db},
		Users:  UserRepo{DB: db},
		Tokens: TokenRepo{DB: db},
	}
}
