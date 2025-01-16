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

	data := app.newTemplateData(r)
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

	data := app.newTemplateData(r)
	data.Snippet = snippet 

	app.render(w, r, http.StatusOK, "view.tmpl.html", data)
}

// Create a new snippet
func (app *Application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 4096)

	validated := validateCreateSnippetForm(r)

	if !validated.Valid() {
		data := app.newTemplateData(r)
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

	app.sessionManager.Put(
		r.Context(), "flash", "Snippet successfully created!",
	)

	http.Redirect(
		w, r, fmt.Sprintf("/snippets/%d", id), http.StatusSeeOther,
	)
}

func (app *Application) createSnippetPage(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = createSnippetForm{Expires: 1}

	app.render(w, r, http.StatusOK, "create.tmpl.html", data)
}

func (app *Application) userSignupPage(
	w http.ResponseWriter, r *http.Request,
) {
	fmt.Fprintln(w, "Display a form for signing up users")
}

func (app *Application) createUser(
	w http.ResponseWriter, r *http.Request,
) {
	fmt.Fprintln(w, "Create a user in the database")
}

func (app *Application) userLoginPage(
	w http.ResponseWriter, r *http.Request,
) {
	fmt.Fprintln(w, "Show a login page")
}

func (app *Application) loginUser(
	w http.ResponseWriter, r *http.Request,
) {
	fmt.Fprintln(w, "Login an user")
}

func (app *Application) logoutUser(
	w http.ResponseWriter, r *http.Request,
) {
	fmt.Fprintln(w, "Logout the user")
}
