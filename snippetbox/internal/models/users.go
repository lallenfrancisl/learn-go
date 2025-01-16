package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserRepo struct {
	DB *sql.DB
}

func (r *UserRepo) Insert(name, email, password string) error {
	return nil
}

func (r *UserRepo) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (r *UserRepo) Exists(id int) (bool, error) {
	return false, nil
}

