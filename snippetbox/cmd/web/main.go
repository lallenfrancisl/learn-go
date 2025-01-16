package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
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
	logger         *log.Logger
	snippets       *models.SnippetRepo
	templateCache  map[string]*template.Template
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
		logger:         logger,
		snippets:       &models.SnippetRepo{DB: db},
		templateCache:  templateCache,
		sessionManager: sessionManager,
	}

	server := &http.Server{
		Addr:     cfg.addr,
		Handler:  app.routes(),
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
		TLSConfig: &tls.Config{
			CurvePreferences: []tls.CurveID{
				tls.X25519,
				tls.CurveP256,
			},
		},
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app.logger.Info(fmt.Sprintf("server started at %s", cfg.addr))

	err = server.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")

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
