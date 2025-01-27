package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/lallenfrancisl/snippetbox/internal/models"
	"github.com/lallenfrancisl/snippetbox/internal/validator"
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

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *Application) userSignupPage(
	w http.ResponseWriter, r *http.Request,
) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, r, http.StatusOK, "signup.tmpl.html", data)
}

func (app *Application) createUser(
	w http.ResponseWriter, r *http.Request,
) {
	data := app.newTemplateData(r)
	form := app.validateUserSignupForm(r)

	if !form.Valid() {
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl.html", data)

		return
	}

	err := app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		} else {
			app.serverError(w, r, err)
		}

		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")
	http.Redirect(w, r, "/users/login", http.StatusSeeOther)
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *Application) userLoginPage(
	w http.ResponseWriter, r *http.Request,
) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, r, http.StatusOK, "login.tmpl.html", data)
}

func (app *Application) loginUser(
	w http.ResponseWriter, r *http.Request,
) {
	data := app.newTemplateData(r)

	form := app.validateLoginUserForm(r)
	if !form.Valid() {
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl.html", data)

		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddFormError("Email or password is incorrect")
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		} else {
			app.serverError(w, r, err)
		}

		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)

		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	http.Redirect(w, r, "/snippets", http.StatusSeeOther)
}

func (app *Application) logoutUser(
	w http.ResponseWriter, r *http.Request,
) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)

		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
