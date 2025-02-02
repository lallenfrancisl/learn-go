package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/lallenfrancisl/greenlight-api/internal/data"
	"github.com/lallenfrancisl/greenlight-api/internal/validator"
)

type createMoviePayload struct {
	Title   string       `json:"title"`
	Year    int32        `json:"year"`
	Runtime data.Runtime `json:"runtime"`
	Genres  []string     `json:"genres"`
}

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input createMoviePayload

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)

		return
	}

	movie := data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	v := validator.New()
	data.ValidateMovie(v, movie)

	if !v.Valid() {
		app.validationFailedResponse(w, r, v.Errors)

		return
	}

	err = app.repo.Movies.Insert(&movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)

		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParams(r)
	if err != nil {
		app.notFoundResponse(w, r)

		return
	}

	movie, err := app.repo.Movies.Get(id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)

			return
		}

		app.serverErrorResponse(w, r, err)

		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)

		return
	}
}

type updateMoviePayload struct {
	Title   *string       `json:"title"`
	Year    *int32        `json:"year"`
	Runtime *data.Runtime `json:"runtime"`
	Genres  []string      `json:"genres"`
}

func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParams(r)
	if err != nil {
		app.notFoundResponse(w, r)

		return
	}

	movie, err := app.repo.Movies.Get(id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		}

		app.serverErrorResponse(w, r, err)

		return
	}

	var payload updateMoviePayload

	err = app.readJSON(w, r, &payload)
	if err != nil {
		app.serverErrorResponse(w, r, err)

		return
	}

	if payload.Title != nil {
		movie.Title = *payload.Title
	}

	if payload.Year != nil {
		movie.Year = *payload.Year
	}

	if payload.Runtime != nil {
		movie.Runtime = *payload.Runtime
	}

	if payload.Genres != nil {
		movie.Genres = payload.Genres
	}

	validator := validator.New()
	data.ValidateMovie(validator, *movie)

	if !validator.Valid() {
		app.validationFailedResponse(w, r, validator.Errors)

		return
	}

	err = app.repo.Movies.Update(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)

		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)

		return
	}
}

func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParams(r)
	if err != nil {
		app.notFoundResponse(w, r)

		return
	}

	err = app.repo.Movies.Delete(id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	err = app.writeJSON(
		w, http.StatusOK, envelope{"message": "movie deleted sucessfully"}, nil,
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)

		return
	}
}
