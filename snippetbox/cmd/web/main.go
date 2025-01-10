package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	assets := flag.String("assets", "./ui/static/", "Static assets folder")

	flag.Parse()

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(*assets))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/{$}", home)
	mux.HandleFunc("GET /snippets/{id}", snippetView)
	mux.HandleFunc("POST /snippets", snippetCreate)

	log.Println(fmt.Sprintf("starting server on %s", *addr))

	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}
