package main

import (
	"os"
	"sync"

	"github.com/redis/go-redis/v9"
	"goauthbackend.bighead.dev/internal/jsonlog"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	redisURL string
}

type application struct {
	config      config
	logger      *jsonlog.Logger
	wg          sync.WaitGroup
	redisClient *redis.Client
}

func main() {
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	cfg, err := updateConfigWithEnvVariables()
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	db, err := openDB(*cfg)

	if err != nil {
		logger.PrintFatal(err, nil)
	}

	defer db.Close()

	opt, err := redis.ParseURL(cfg.redisURL)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	client := redis.NewClient(opt)

	logger.PrintInfo("redis connection poll estabilished", nil)

	app := application{
		config:      *cfg,
		logger:      logger,
		redisClient: client,
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}

}
