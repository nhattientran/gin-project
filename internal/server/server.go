package server

import (
	"flag"
	"fmt"
	"gin-project/internal/data"
	logger "gin-project/internal/log"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"gin-project/internal/database"
)

const version = "1.0.0"

type config struct {
	port    int
	env     string
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
}

type Server struct {
	config config

	models   data.Models
	infoLog  *logger.Logger
	errorLog *logger.Logger
	db       database.Service
}

var (
	InfoLog    = logger.New(os.Stdout, logger.LevelInfo)
	WarningLog = logger.New(os.Stdout, logger.LevelWarn)
	ErrorLog   = logger.New(os.Stderr, logger.LevelError)
	FatalLog   = logger.New(os.Stderr, logger.LevelFatal)
)

func NewServer() *http.Server {

	var cfg config
	port, _ := strconv.ParseInt(os.Getenv("PORT"), 10, 64)
	flag.IntVar(&cfg.port, "port", int(port), "API server port")
	flag.StringVar(&cfg.env, "env", os.Getenv("APP_ENV"), "Environment (development|staging|production)")

	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")
	flag.Parse()

	db := database.New()
	NewServer := &Server{
		config:   cfg,
		infoLog:  InfoLog,
		errorLog: ErrorLog,
		models:   data.NewModels(db),
		db:       db,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      NewServer.RegisterRoutes(),
		ErrorLog:     log.New(InfoLog, "", 0),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	InfoLog.PrintInfo("starting server", map[string]string{
		"addr": server.Addr,
		"env":  cfg.env,
	})

	return server
}
