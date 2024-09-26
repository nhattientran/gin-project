package server

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"gin-project/internal/data"
	logger "gin-project/internal/log"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
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
	config     config
	models     data.Models
	infoLog    *logger.Logger
	errorLog   *logger.Logger
	warningLog *logger.Logger
	fatalLog   *logger.Logger
	db         database.Service
}

var (
	InfoLog    = logger.New(os.Stdout, logger.LevelInfo)
	WarningLog = logger.New(os.Stdout, logger.LevelWarn)
	ErrorLog   = logger.New(os.Stderr, logger.LevelError)
	FatalLog   = logger.New(os.Stderr, logger.LevelFatal)
)

func NewServer() error {

	var cfg config
	port, _ := strconv.ParseInt(os.Getenv("PORT"), 10, 64)
	flag.IntVar(&cfg.port, "port", int(port), "API server port")
	flag.StringVar(&cfg.env, "env", os.Getenv("APP_ENV"), "Environment (development|staging|production)")

	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")
	flag.Parse()

	db := database.New()
	defer db.Close()
	NewServer := &Server{
		config:     cfg,
		infoLog:    InfoLog,
		errorLog:   ErrorLog,
		warningLog: WarningLog,
		fatalLog:   FatalLog,
		models:     data.NewModels(&db),
		db:         db,
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

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		InfoLog.PrintInfo("shutting down server", map[string]string{
			"signal": s.String(),
		})
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		shutdownError <- server.Shutdown(ctx)

		os.Exit(0)
	}()

	InfoLog.PrintInfo("starting server", map[string]string{
		"addr": server.Addr,
		"env":  cfg.env,
	})

	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	InfoLog.PrintInfo("stopped server", map[string]string{
		"addr": server.Addr,
	})

	return nil
}
