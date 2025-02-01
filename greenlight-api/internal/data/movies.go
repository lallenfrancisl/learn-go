package data

import (
	"database/sql"
	"time"
)

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Runtime   Runtime   `json:"runtime,omitempty"`
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}

type MovieRepo struct {
	DB *sql.DB
}

func (r *MovieRepo) Insert(movie *Movie) error {
	return nil
}

func (r *MovieRepo) Get(id int64) (*Movie, error) {
	return nil, nil
}

func (r *MovieRepo) Update(movie *Movie) error {
	return nil
}

func (r *MovieRepo) Delete(id int64) error {
	return nil
}
