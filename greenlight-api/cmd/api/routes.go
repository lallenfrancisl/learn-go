package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.RedirectTrailingSlash = true
	router.RedirectFixedPath = true

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// Healthcheck routes
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	// Movies routes
	router.HandlerFunc(
		http.MethodPost,
		"/v1/movies",
		app.requireActivatedUser(app.createMovieHandler),
	)
	router.HandlerFunc(
		http.MethodGet,
		"/v1/movies",
		app.requireActivatedUser(app.listMoviesHandler),
	)
	router.HandlerFunc(
		http.MethodGet,
		"/v1/movies/:id",
		app.requireActivatedUser(app.showMovieHandler),
	)
	router.HandlerFunc(
		http.MethodPatch,
		"/v1/movies/:id",
		app.requireActivatedUser(app.updateMovieHandler),
	)
	router.HandlerFunc(
		http.MethodDelete,
		"/v1/movies/:id",
		app.requireActivatedUser(app.deleteMovieHandler),
	)

	// Users routes
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/:id/activate", app.activateUserHandler)
	router.HandlerFunc(
		http.MethodPost, "/v1/users/login",
		app.loginHandler,
	)

	return app.recoverPanic(
		app.rateLimit(app.authenticate(router)),
	)
}
