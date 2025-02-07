package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/lallenfrancisl/greenlight-api/internal/data"
	"github.com/lallenfrancisl/greenlight-api/internal/validator"
)

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)

		return
	}

	v := validator.New()

	data.ValidateEmail(v, input.Email)
	data.ValidatePassword(v, input.Password)

	if !v.Valid() {
		app.validationFailedResponse(w, r, v.Errors)

		return
	}

	user, err := app.repo.Users.GetByEmail(input.Email)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	match, err := user.Password.Matches(input.Password)
	if !match {
		app.invalidCredentialsResponse(w, r)

		return
	}

	token, err := app.repo.Tokens.New(
		user.ID, 24*time.Hour, data.ScopeAuthentication,
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)

		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"credentials": token}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)

		return
	}
}
