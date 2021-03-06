package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"github.com/Neumann88/payment-api-emulator/config"

	"github.com/Neumann88/payment-api-emulator/internal/controller"
	"github.com/Neumann88/payment-api-emulator/internal/repository"
	"github.com/Neumann88/payment-api-emulator/internal/usecase"

	"github.com/Neumann88/payment-api-emulator/pkg/db/postgres"
	"github.com/Neumann88/payment-api-emulator/pkg/http/server"
	"github.com/Neumann88/payment-api-emulator/pkg/loggin"
)

func main() {
	// Config
	cfg, err := config.NewConfig()

	if err != nil {
		log.Fatalf("config initialization error: %s", err.Error())
	}

	// Logger
	logger := loggin.NewLogger(cfg.Logger.Debug)

	// Database
	dbOptions := postgres.DBOptions{
		User:     cfg.Postgres.User,
		Password: cfg.Postgres.Password,
		Host:     cfg.Postgres.Host,
		Port:     cfg.Postgres.Port,
		DB:       cfg.Postgres.DB,
		SSLmode:  cfg.Postgres.SSLMode,
	}

	db := postgres.NewPostgres(
		dbOptions,
		cfg.Postgres.ConnAttempts,
		time.Duration(cfg.Postgres.ConnTimeout)*time.Second,
	)

	pg, err := db.Connect()

	if err != nil {
		logger.Fatalf("postgres connection failed, %s", err.Error())
	}

	defer pg.Close()

	// Scheme migrate
	postgres.InitMigrate(
		logger,
		dbOptions,
	)

	// Layers
	repo := repository.NewPaymentRepository(pg)
	usc := usecase.NewPaymentUseCase(repo)
	con := controller.NewPaymentController(
		logger,
		usc,
	)

	port := os.Getenv("PORT")
	if port == "" {
	  port = cfg.HTTP.Port
	}
	// HTTP-Server
	router := mux.NewRouter()

	httpServer := server.NewHTTPServer(
		con.Register(router),
		port,
		time.Duration(cfg.HTTP.ReadTimeout)*time.Second,
		time.Duration(cfg.HTTP.WriteTimeout)*time.Second,
		time.Duration(cfg.HTTP.ShutdownTimeout)*time.Second,
	)

	logger.Infof("http server created and started at http://localhost:%s", cfg.HTTP.Port)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Infof("app - run - signal: %s", s.String())
	case err = <-httpServer.Notify():
		logger.Errorf("app - run - httpServer.Notify: %s", err.Error())
	}

	// Shutdown
	err = httpServer.Shutdown()

	if err != nil {
		logger.Errorf("app - run - httpServer.Shutdown: %s", err.Error())
	}
}
