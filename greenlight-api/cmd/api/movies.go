package main

import (
	"net/http"
	"time"

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

	v := validator.New()
	app.validateCreateMoviePayload(v, input)

	if !v.Valid() {
		app.validationFailedResponse(w, r, v.Errors)

		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"movie": input}, nil)
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParams(r)
	if err != nil {
		app.notFoundResponse(w, r)

		return
	}

	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Maheshinte Prathikaram",
		Runtime:   102,
		Genres:    []string{"drama", "romance", "naturalism"},
		Version:   1,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)

		return
	}
}
