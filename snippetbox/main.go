package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {	
	w.Write([]byte("Hello from Snippetbox"))
}

// Show a specific snippet by id
func snippetView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a snippet"))
}

// Create a new snippet
func snippetCreate(w http.ResponseWriter, r *http.Request)  {
	w.Write([]byte("Create a snippet"))	
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippets/view", snippetView)
	mux.HandleFunc("/snippets/create", snippetCreate)

	log.Println("starting server on localhost:4000")

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
