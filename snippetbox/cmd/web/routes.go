package main

import (
	"net/http"

	"github.com/justinas/alice"
	"github.com/lallenfrancisl/snippetbox/ui"
)

func (app *Application) routes() http.Handler {
	fileServer := http.FileServerFS(ui.Files)
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("GET /ping", ping)

	dynamic := alice.New(app.sessionManager.LoadAndSave, app.noSurf, app.authenticate)
	protected := dynamic.Append(app.requireAuthentication)

	mux.Handle("/{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippets/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /users/signup", dynamic.ThenFunc(app.userSignupPage))
	mux.Handle("POST /users/signup", dynamic.ThenFunc(app.createUser))
	mux.Handle("GET /users/login", dynamic.ThenFunc(app.userLoginPage))
	mux.Handle("POST /users/login", dynamic.ThenFunc(app.loginUser))

	mux.Handle("POST /users/logout", protected.ThenFunc(app.logoutUser))
	mux.Handle("POST /snippets", protected.ThenFunc(app.snippetCreate))
	mux.Handle("GET /snippets/create", protected.ThenFunc(app.createSnippetPage))

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
}
