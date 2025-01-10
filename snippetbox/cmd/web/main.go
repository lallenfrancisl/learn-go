package main

import (
	"flag"
	"fmt"
	"net/http"

	log "github.com/lallenfrancisl/snippetbox/internal"
)

type Config struct {
	addr   string
	assets string
}

var cfg Config
var logger = log.New()

func main() {
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.assets, "assets", "./ui/static/", "Static assets folder")

	flag.Parse()

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(cfg.assets))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/{$}", home)
	mux.HandleFunc("GET /snippets/{id}", snippetView)
	mux.HandleFunc("POST /snippets", snippetCreate)

	logger.Info(fmt.Sprintf("starting server on %s", cfg.addr))

	err := http.ListenAndServe(cfg.addr, mux)
	logger.Fatal(err.Error())
}
