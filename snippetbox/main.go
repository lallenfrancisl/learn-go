package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Snippetbox"))

	log.Println("GET /")
}

// Show a specific snippet by id
func snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)

		return
	}

	msg := fmt.Sprintf("Snippet of id %d", id)
	w.Write([]byte(msg))

	log.Println("GET /snippets/{id}")
}

// Create a new snippet
func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create a snippet"))

	log.Println("GET /snippets")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/{$}", home)
	mux.HandleFunc("GET /snippets/{id}", snippetView)
	mux.HandleFunc("POST /snippets", snippetCreate)

	log.Println("starting server on localhost:4000")

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
