package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/lallenfrancisl/snippetbox/internal"
	"github.com/lallenfrancisl/snippetbox/internal/models"
)

type Config struct {
	addr   string
	assets string
	dsn    string
}

type Application struct {
	logger        *log.Logger
	snippets      *models.SnippetRepo
	templateCache map[string]*template.Template
	sessionManager *scs.SessionManager
}

var cfg Config

func main() {
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.assets, "assets", "./ui/static/", "Static assets folder")
	flag.StringVar(&cfg.dsn, "dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")

	flag.Parse()

	logger := log.New()

	db, err := openDB(cfg.dsn)
	if err != nil {
		logger.Fatal(err.Error())

		return
	}

	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Fatal(err.Error())

		return
	}

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	app := &Application{
		logger:        logger,
		snippets:      &models.SnippetRepo{DB: db},
		templateCache: templateCache,
		sessionManager: sessionManager,
	}

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
