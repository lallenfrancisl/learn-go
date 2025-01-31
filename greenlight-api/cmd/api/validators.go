package main

import (
	"time"

	"github.com/lallenfrancisl/greenlight-api/internal/validator"
)

func (app *application) validateCreateMoviePayload(v *validator.Validator, input createMoviePayload) {
	v.Check(validator.NotBlank(
		input.Title), "title", "must be provided",
	)
	v.Check(
		validator.MaxLen(input.Title, 500),
		"title", "must not be more than 500 characters long",
	)

	v.Check(!validator.Equal(string(input.Year), "0"), "year", "must be provided")
	v.Check(validator.Min(int(input.Year), 1888), "year", "must be greater than 1888")
	v.Check(validator.Max(int(input.Year), time.Now().Year()), "year", "must not be in the future")

	v.Check(!validator.Equal(string(input.Runtime), "0"), "runtime", "must be provided")
	v.Check(
		validator.GreaterThan(int(input.Runtime), 0),
		"runtime", "must be a positive integer",
	)

	v.Check(
		validator.Min(len(input.Genres), 1), "genres", "must contain atleast 1 genre",
	)
	v.Check(validator.Unique(input.Genres), "genres", "must not contain duplicate values")
}
