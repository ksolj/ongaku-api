package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ksolj/ongaku-api/internal/data"
	"github.com/ksolj/ongaku-api/internal/jsonlog"
)

const version = "1.0.0"

type config struct {
	port int
	env  string // Name of the current operating environment for the application (development, staging, production, etc.)
	db   struct {
		dsn string
	}
}

type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	// Read the DSN value from the db-dsn command-line flag into the config struct. We
	// default to using our development DSN if no flag is provided.
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("ONGAKU_DB_DSN"), "PostgreSQL DSN")

	flag.Parse()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo) // maybe use zerolog in the future???

	pool, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	defer pool.Close()
	logger.PrintInfo("database connection pool established", nil)

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(pool),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		ErrorLog:     log.New(logger, "", 0),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Again, we use the PrintInfo() method to write a "starting server" message at the
	// INFO level. But this time we pass a map containing additional properties (the
	// operating environment and server address) as the final parameter.
	logger.PrintInfo("starting server", map[string]string{
		"addr": srv.Addr,
		"env":  cfg.env,
	})

	err = srv.ListenAndServe()
	// Use the PrintFatal() method to log the error and exit.
	logger.PrintFatal(err, nil)
}

func openDB(cfg config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use Ping() to establish a new connection to the database, passing in the
	// context we created above as a parameter. If the connection couldn't be
	// established successfully within the 5 second deadline, then this will return an
	// error.
	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
