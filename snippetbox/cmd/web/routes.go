package main

import "net/http"

func (app *Application) routes() *http.ServeMux {
	fileServer := http.FileServer(http.Dir(cfg.assets))
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/{$}", app.home)
	mux.HandleFunc("GET /snippets/{id}", app.snippetView)
	mux.HandleFunc("POST /snippets", app.snippetCreate)

	return mux
}
