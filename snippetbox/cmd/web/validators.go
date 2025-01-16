package main

import (
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
)

type createSnippetForm struct {
	Title       string
	Content     string
	Expires     int
	FieldErrors map[string]string
}

func validateCreateSnippetForm(r *http.Request) (data createSnippetForm) {
	err := r.ParseForm()
	validated := createSnippetForm{
		FieldErrors: make(map[string]string),
	}

	if err != nil {
		validated.FieldErrors["form"] = err.Error()
		return validated
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")

	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		validated.FieldErrors["form"] = err.Error()
		return validated
	}

	if strings.TrimSpace(title) == "" {
		validated.FieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		validated.FieldErrors["title"] = "This field cannot be more than 100 characters long"
	}

	if strings.TrimSpace(content) == "" {
		validated.FieldErrors["content"] = "This field cannot be blank"
	}

	if expires != 1 && expires != 7 && expires != 365 {
		validated.FieldErrors["expires"] = "This field must equal 1, 7 or 365"
	}

	if len(validated.FieldErrors) == 0 {
        validated.Title = title
        validated.Content = content
        validated.Expires = expires
	}

	return validated
}
