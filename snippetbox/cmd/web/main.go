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

	app := Application{
		logger: log.New(),
	}

	app.logger.Info(fmt.Sprintf("server started at %s", cfg.addr))

	err := http.ListenAndServe(cfg.addr, app.routes())

	if err != nil {
		app.logger.Fatal(err.Error())
	}
}
