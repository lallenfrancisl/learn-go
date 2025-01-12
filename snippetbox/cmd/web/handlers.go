package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/lallenfrancisl/snippetbox/internal/models"
)

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("GET /")

	w.Header().Add("Server", "Go")

	files := []string{
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/base.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, r, err)

		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, r, err)
	}
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

	app.logger.Info(fmt.Sprintf("GET /snippets/%d", id))

	msg := fmt.Sprintf("%+v", snippet)
	w.Write([]byte(msg))
}

// Create a new snippet
func (app *Application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	title := "test title"
	content := "Snippet test content"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)

		return
	}

	app.logger.Info("POST /snippets")

	http.Redirect(w, r, fmt.Sprintf("/snippets/%d", id), http.StatusSeeOther)
}
