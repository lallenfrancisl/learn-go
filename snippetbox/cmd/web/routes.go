package main

import (
	"github.com/justinas/alice"
	"net/http"
)

func (app *Application) routes() http.Handler {
	fileServer := http.FileServer(http.Dir(cfg.assets))
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	mux.Handle("/{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippets/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("POST /snippets", dynamic.ThenFunc(app.snippetCreate))
	mux.Handle("GET /snippets", dynamic.ThenFunc(app.createSnippetPage))

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
}
