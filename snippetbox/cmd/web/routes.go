package main

import "net/http"

func (app *Application) routes() http.Handler {
	fileServer := http.FileServer(http.Dir(cfg.assets))
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/{$}", app.home)
	mux.HandleFunc("GET /snippets/{id}", app.snippetView)
	mux.HandleFunc("POST /snippets", app.snippetCreate)

	return app.recoverPanic(app.logRequest(commonHeaders(mux)))
}
