package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/lallenfrancisl/greenlight-api/internal/validator"
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

func ValidateMovie(v *validator.Validator, movie Movie) {
	v.Check(validator.NotBlank(
		movie.Title), "title", "must be provided",
	)
	v.Check(
		validator.MaxLen(movie.Title, 500),
		"title", "must not be more than 500 characters long",
	)

	v.Check(!validator.Equal(string(movie.Year), "0"), "year", "must be provided")
	v.Check(validator.Min(int(movie.Year), 1888), "year", "must be greater than 1888")
	v.Check(validator.Max(int(movie.Year), time.Now().Year()), "year", "must not be in the future")

	v.Check(!validator.Equal(string(movie.Runtime), "0"), "runtime", "must be provided")
	v.Check(
		validator.GreaterThan(int(movie.Runtime), 0),
		"runtime", "must be a positive integer",
	)

	v.Check(
		validator.Min(len(movie.Genres), 1), "genres", "must contain atleast 1 genre",
	)
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")
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
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
		SELECT id, created_at, title, year, runtime, genres, version
		FROM movies
		WHERE id = $1
	`

	var movie Movie

	err := r.DB.QueryRow(query, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.Version,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}

		return nil, err
	}

	return &movie, nil
}

func (r *MovieRepo) Update(movie *Movie) error {
	query := `
		UPDATE movies
		SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
		WHERE id = $5 AND version = $6
		RETURNING version
	`

	args := []interface{}{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
		movie.ID,
		movie.Version,
	}

	err := r.DB.QueryRow(query, args...).Scan(&movie.Version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrEditConflict
		}

		return err
	}

	return nil
}

func (r *MovieRepo) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM movies WHERE id = $1`

	result, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
