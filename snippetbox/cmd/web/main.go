package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/lallenfrancisl/snippetbox/internal"
)

type Config struct {
	addr   string
	assets string
	dsn    string
}

type Application struct {
	logger *log.Logger
}

var cfg Config

func main() {
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.assets, "assets", "./ui/static/", "Static assets folder")
	flag.StringVar(&cfg.dsn, "dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")

	flag.Parse()

	app := Application{
		logger: log.New(),
	}

	db, err := openDB(cfg.dsn)
	if err != nil {
		app.logger.Fatal(err.Error())

		return
	}

	defer db.Close()

	app.logger.Info(fmt.Sprintf("server started at %s", cfg.addr))

	err = http.ListenAndServe(cfg.addr, app.routes())

	if err != nil {
		app.logger.Fatal(err.Error())
	}
}

// The openDB() function wraps sql.Open() and returns a sql.DB connection pool
// for a given DSN.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
