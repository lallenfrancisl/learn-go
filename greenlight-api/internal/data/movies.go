package data

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
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
	query := `
		INSERT INTO movies (title, year, runtime, genres)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version
	`

	args := []interface{}{
		movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres),
	}

	return r.DB.QueryRow(query, args...).Scan(
		&movie.ID, &movie.CreatedAt, &movie.Version,
	)
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
