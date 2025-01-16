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
