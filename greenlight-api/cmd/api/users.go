package main

import (
	"errors"
	"net/http"
	"time"

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

	err = app.repo.Permissions.AddForUser(
		user.ID, data.PermissionMoviesRead,
		data.PermissionMoviesWrite, data.PermissionMoviesDelete,
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)

		return
	}

	token, err := app.repo.Tokens.New(
		user.ID, 3*24*time.Hour, data.ScopeActivation,
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)

		return
	}

	app.background(func() {
		data := map[string]interface{}{
			"activationToken": token.Plaintext,
			"userID":          user.ID,
		}

		err = app.mailer.Send(user.Email, "user_welcome.tmpl.html", data)
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

type activateUserPayload struct {
	TokenPlaintext string `json:"token"`
}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input activateUserPayload

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)

		return
	}

	v := validator.New()

	data.ValidateTokenPlaintext(v, input.TokenPlaintext)
	if !v.Valid() {
		app.validationFailedResponse(w, r, v.Errors)

		return
	}

	user, err := app.repo.Users.GetByToken(data.ScopeActivation, input.TokenPlaintext)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			v.AddError("token", "invalid or expired activation token")
			app.validationFailedResponse(w, r, v.Errors)
		} else {
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	if user.Activated {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid operation")

		return
	}

	user.Activated = true

	err = app.repo.Users.Update(user)
	if err != nil {
		if errors.Is(err, data.ErrEditConflict) {
			app.editConflictResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	err = app.repo.Tokens.DeleteAllOfUser(data.ScopeActivation, user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)

		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)

		return
	}
}
