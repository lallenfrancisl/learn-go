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

type Application struct {
	logger *log.Logger
}

var cfg Config

func main() {
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.assets, "assets", "./ui/static/", "Static assets folder")

	flag.Parse()

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(cfg.assets))

	app := Application{
		logger: log.New(),
	}

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/{$}", app.home)
	mux.HandleFunc("GET /snippets/{id}", app.snippetView)
	mux.HandleFunc("POST /snippets", app.snippetCreate)

	app.logger.Info(fmt.Sprintf("server started at %s", cfg.addr))

	err := http.ListenAndServe(cfg.addr, mux)

	if err != nil {
		app.logger.Fatal(err.Error())
	}
}
