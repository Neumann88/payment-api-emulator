package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Neumann88/payment-api-emulator/config"
	"github.com/Neumann88/payment-api-emulator/pkg/db/migrate"
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
	logger := loggin.NewLogger().Init(cfg.Logger.Debug)

	// Database
	pg, err := postgres.NewPostgres().Connect(cfg.Postgres.Dsn)
	if err != nil {
		logger.Fatalf("postgres connection failed, %s", err.Error())
	}
	defer pg.Close()

	// Scheme migrate
	migrate.InitMigrate(logger, cfg.Postgres.Dsn)

	// Entities
	// repos := repository.NewRepository(pg)
	// services := service.NewService(repos)
	// handlers := handler.NewHandler(services)

	h := mux.NewRouter()

	h.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	// HTTP-Server
	httpServer := server.NewHttpServer(
		h,
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
