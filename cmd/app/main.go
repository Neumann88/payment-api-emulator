package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Neumann88/payment-api-emulator/config"
	"github.com/Neumann88/payment-api-emulator/internal/payment"
	"github.com/Neumann88/payment-api-emulator/pkg/db/postgres"
	"github.com/Neumann88/payment-api-emulator/pkg/http/server"
	"github.com/Neumann88/payment-api-emulator/pkg/loggin"
	"github.com/gorilla/mux"
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

	pg, err := postgres.NewPostgres(dbOptions).Connect()

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
	rep := payment.NewPaymentRepository(pg)
	usc := payment.NewPaymentUseCase(rep)
	con := payment.NewPaymentController(
		logger,
		usc,
	)

	// HTTP-Server
	router := mux.NewRouter()

	httpServer := server.NewHttpServer(
		con.Register(router),
		cfg.HTTP.Port,
		time.Duration(cfg.HTTP.ReadTimeout),
		time.Duration(cfg.HTTP.WriteTimeout),
		time.Duration(cfg.HTTP.ShutdownTimeout),
	)

	logger.Infof("http server created and started at http://localhost:%s", cfg.HTTP.Port)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Infof("app - run - signal: %s", s.String())
	case err := <-httpServer.Notify():
		logger.Errorf("app - run - httpServer.Notify: %s", err.Error())
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		logger.Errorf("app - run - httpServer.Shutdown: %s", err.Error())
	}
}
