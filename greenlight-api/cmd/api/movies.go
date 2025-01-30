package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/lallenfrancisl/greenlight-api/internal/data"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "create a new movie")
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParams(r)
	if err != nil {
		http.NotFound(w, r)

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
		app.serverError(w, err)

		return
	}
}
