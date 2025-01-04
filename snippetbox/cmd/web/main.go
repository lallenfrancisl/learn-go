package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/{$}", home)
	mux.HandleFunc("GET /snippets/{id}", snippetView)
	mux.HandleFunc("POST /snippets", snippetCreate)

	log.Println("starting server on localhost:4000")

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
