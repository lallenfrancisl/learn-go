package main

import (
	"errors"
	"net/http"

	"github.com/lallenfrancisl/greenlight-api/internal/data"
	"github.com/lallenfrancisl/greenlight-api/internal/validator"
)

type RegisterUserPayload struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input RegisterUserPayload

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)

		return
	}

	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)

		return
	}

	v := validator.New()

	data.ValidateUser(v, *user)
	if !v.Valid() {
		app.validationFailedResponse(w, r, v.Errors)

		return
	}

	err = app.repo.Users.Insert(user)
	if err != nil {
		if errors.Is(err, data.ErrDuplicateEmail) {
			v.AddError("email", "a user with this email address already exists")
			app.validationFailedResponse(w, r, v.Errors)
		} else {
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	app.background(func() {
		err = app.mailer.Send(user.Email, "user_welcome.tmpl.html", user)
		if err != nil {
			app.logError(r, err)
		}
	})

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)

		return
	}
}
