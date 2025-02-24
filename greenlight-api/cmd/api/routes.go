package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/lallenfrancisl/greenlight-api/internal/data"
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
		app.requirePermission(data.PermissionMoviesWrite, app.createMovieHandler),
	)
	router.HandlerFunc(
		http.MethodGet,
		"/v1/movies",
		app.listMoviesHandler,
	)
	router.HandlerFunc(
		http.MethodGet,
		"/v1/movies/:id",
		app.showMovieHandler,
	)
	router.HandlerFunc(
		http.MethodPatch,
		"/v1/movies/:id",
		app.requirePermission(data.PermissionMoviesWrite, app.updateMovieHandler),
	)
	router.HandlerFunc(
		http.MethodDelete,
		"/v1/movies/:id",
		app.requirePermission(data.PermissionMoviesDelete, app.deleteMovieHandler),
	)

	// Users routes
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/:id/activate", app.activateUserHandler)
	router.HandlerFunc(
		http.MethodPost, "/v1/users/login",
		app.loginHandler,
	)

	// Docs routes
	router.HandlerFunc(http.MethodGet, "/v1/docs", app.GetDocs)

	return app.recoverPanic(
		app.enableCORS(app.rateLimit(
			app.logRequest(app.authenticate(router)),
		)),
	)
}
