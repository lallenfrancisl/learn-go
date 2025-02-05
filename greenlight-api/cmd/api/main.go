package main

import (
	"context"
	"database/sql"
	"flag"
	"os"
	"sync"
	"time"

	jsonlog "github.com/lallenfrancisl/greenlight-api/internal"
	"github.com/lallenfrancisl/greenlight-api/internal/data"
	"github.com/lallenfrancisl/greenlight-api/internal/mailer"
	_ "github.com/lib/pq"
)

const version = "1.0.0"

type config struct {
	Port int
	Env  string
	DB   struct {
		DSN          string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  string
	}
	Limiter struct {
		RPS     float64
		Burst   int
		Enabled bool
	}
	Mailer struct {
		Host     string
		Port     int
		Username string
		Password string
		Sender   string
	}
}

type application struct {
	config config
	logger *jsonlog.Logger
	repo   data.Repo
	mailer mailer.Mailer
	wg     sync.WaitGroup
}

func main() {
	var cfg config

	parseFlags(&cfg)

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err, nil)
	}
	defer db.Close()

	logger.Info("database connection pool established", nil)

	mailerInst := mailer.New(
		cfg.Mailer.Host,
		cfg.Mailer.Port,
		cfg.Mailer.Username,
		cfg.Mailer.Password,
		cfg.Mailer.Sender,
		logger,
	)

	app := &application{
		config: cfg,
		logger: logger,
		repo:   data.NewRepo(db),
		mailer: mailerInst,
	}

	err = app.serve()
	if err != nil {
		logger.Fatal(err, nil)
	}
}

func parseFlags(cfg *config) {
	// App
	flag.IntVar(&cfg.Port, "port", 4000, "API server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|staging|production)")

	// DB
	flag.StringVar(
		&cfg.DB.DSN,
		"db-dsn", "postgres://greenlight:password@localhost/greenlight?sslmode=disable",
		"PostgreSQL DSN",
	)
	flag.IntVar(&cfg.DB.MaxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.DB.MaxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.DB.MaxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	// Rate limiter
	flag.Float64Var(&cfg.Limiter.RPS, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.Limiter.Burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.Limiter.Enabled, "limiter-enabled", true, "Enable rate limiter")

	// Mailer
	flag.StringVar(&cfg.Mailer.Host, "mailer-host", "", "Mailer hostname")
	flag.IntVar(&cfg.Mailer.Port, "mailer-port", 587, "Mailer port")
	flag.StringVar(&cfg.Mailer.Username, "mailer-username", "", "Mailer username")
	flag.StringVar(&cfg.Mailer.Password, "mailer-password", "", "Mailer password")
	flag.StringVar(&cfg.Mailer.Sender, "mailer-sender", "", "Mailer sender email id")

	flag.Parse()
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DB.DSN)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	db.SetMaxIdleConns(cfg.DB.MaxIdleConns)

	duration, err := time.ParseDuration(cfg.DB.MaxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
