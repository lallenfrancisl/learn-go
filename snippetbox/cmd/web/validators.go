package main

import (
	"net/http"
	"strconv"

	"github.com/lallenfrancisl/snippetbox/internal/validator"
)

type createSnippetForm struct {
	Title   string
	Content string
	Expires int
	validator.Validator
}

func validateCreateSnippetForm(r *http.Request) (data createSnippetForm) {
	err := r.ParseForm()
	form := createSnippetForm{}

	if err != nil {
		form.AddFieldError("form", err.Error())

		return form
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")

	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		form.AddFieldError("form", err.Error())

		return form
	}

	form.Title = title
	form.Content = content
	form.Expires = expires

	form.CheckField(
		validator.NotBlank(form.Title),
		"title", "This field cannot be blank",
	)
	form.CheckField(
		validator.MaxChars(form.Title, 100),
		"title", "This field cannot be more than 100 characters long",
	)
	form.CheckField(
		validator.NotBlank(form.Content),
		"content", "This field cannot be blank",
	)
	form.CheckField(
		validator.PermittedValue(form.Expires, 1, 7, 365),
		"expires", "This field must equal 1, 7 or 365",
	)

	return form
}

func (app *Application) validateUserSignupForm(r *http.Request) (data userSignupForm) {
	form := userSignupForm{}
	err := app.decodePostForm(r, &form)

	if err != nil {
		form.AddFieldError("form", "Cannot parse form")

		return form
	}

	form.CheckField(
		validator.NotBlank(form.Name),
		"name", "This field cannot be blank",
	)
	form.CheckField(
		validator.NotBlank(form.Email),
		"email", "This field cannot be blank",
	)
	form.CheckField(
		validator.Matches(form.Email, validator.EmailRX),
		"email", "This field must be a valid email address",
	)
	form.CheckField(
		validator.NotBlank(form.Password),
		"password", "This field cannot be blank",
	)
	form.CheckField(
		validator.MinChars(form.Password, 8),
		"password", "This field must be at least 8 characters long",
	)

	return form
}
