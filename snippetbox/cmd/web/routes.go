package main

import (
	"github.com/justinas/alice"
	"net/http"
)

func (app *Application) routes() http.Handler {
	fileServer := http.FileServer(http.Dir(cfg.assets))
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/{$}", app.home)
	mux.HandleFunc("GET /snippets/{id}", app.snippetView)
	mux.HandleFunc("POST /snippets", app.snippetCreate)

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
}
