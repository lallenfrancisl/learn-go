package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/lallenfrancisl/snippetbox/internal/models"
)

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)

		return
	}

	data := newTemplateData()
	data.Snippets = snippets

	app.render(w, r, http.StatusOK, "home.tmpl.html", data)
}

// Show a specific snippet by id
func (app *Application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)

		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.clientError(w, http.StatusNotFound)
		} else {
			app.serverError(w, r, err)
		}

		return
	}

	data := newTemplateData()
	data.Snippet = snippet

	app.render(w, r, http.StatusOK, "view.tmpl.html", data)
}

// Create a new snippet
func (app *Application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 4096)

	validated := validateCreateSnippetForm(r)

	if !validated.Valid() {
		data := newTemplateData()
		data.Form = validated
		app.render(
			w, r, http.StatusUnprocessableEntity, "create.tmpl.html", data,
		)

		return
	}

	id, err := app.snippets.Insert(
		validated.Title, validated.Content, validated.Expires,
	)
	if err != nil {
		app.serverError(w, r, err)

		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippets/%d", id), http.StatusSeeOther)
}

func (app *Application) createSnippetPage(w http.ResponseWriter, r *http.Request) {
	data := newTemplateData()
	data.Form = createSnippetForm{Expires: 1}

	app.render(w, r, http.StatusOK, "create.tmpl.html", data)
}
