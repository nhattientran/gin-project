package server

import (
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
	port int
	env  string
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
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	cfg := config{
		port: port,
		env:  os.Getenv("APP_ENV"),
	}

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
