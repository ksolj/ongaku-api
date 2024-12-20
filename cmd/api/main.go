package main

import (
	"context"
	"expvar"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ksolj/ongaku-api/internal/data"
	"github.com/ksolj/ongaku-api/internal/jsonlog"
	"github.com/ksolj/ongaku-api/internal/mailer"
	"github.com/ksolj/ongaku-api/internal/vcs"
	"github.com/redis/go-redis/v9"
)

var (
	version = vcs.Version()
)

type config struct {
	port int
	env  string // Name of the current operating environment for the application (development, staging, production, etc.)
	db   struct {
		sql   string
		redis string
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
	cors struct {
		trustedOrigins []string
	}
}

type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
	mailer mailer.Mailer
	wg     sync.WaitGroup
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.StringVar(&cfg.db.sql, "sql-dsn", "", "PostgreSQL DSN")
	flag.StringVar(&cfg.db.redis, "redis-dsn", "", "Redis DSN")

	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.StringVar(&cfg.smtp.host, "smtp-host", "localhost", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 25, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "not in use for now", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "not in use for now", "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "Ongaku API <no-reply@ongaku.ksolj.net>", "SMTP sender")

	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.cors.trustedOrigins = strings.Fields(val)
		return nil
	})

	displayVersion := flag.Bool("version", false, "Display version and exit")

	flag.Parse()

	if *displayVersion {
		fmt.Printf("Version:\t%s\n", version)
		os.Exit(0)
	}

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo) // maybe use zerolog in the future???

	pool, err := openSQL(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	defer pool.Close()
	logger.PrintInfo("SQL database connection pool established", nil)

	rdb, err := openInMemoryDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	defer rdb.Close()
	logger.PrintInfo("Redis connection established", nil)

	mailer, err := mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	expvar.NewString("version").Set(version)

	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	// Currently pool.Stat() returns something that can't be JSON encoded
	// TODO: Deal with this problem cuz expvar can't display it properly due to that
	// expvar.Publish("database", expvar.Func(func() any {
	// 	return pool.Stat()
	// }))

	expvar.Publish("timestamp", expvar.Func(func() any {
		return time.Now().Unix()
	}))

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(pool, rdb),
		mailer: mailer,
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}

func openSQL(cfg config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), cfg.db.sql)
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

func openInMemoryDB(cfg config) (*redis.Client, error) {
	opts, err := redis.ParseURL(cfg.db.redis)
	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}
