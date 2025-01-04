package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /")

	w.Header().Add("Server", "Go")

	files := []string{
	    "./ui/html/partials/nav.tmpl.html",
	    "./ui/html/base.tmpl.html",
	    "./ui/html/pages/home.tmpl.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
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
	w.WriteHeader(http.StatusCreated)

	w.Write([]byte("Create a snippet"))

	log.Println("GET /snippets")
}
