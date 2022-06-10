package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Neumann88/payment-api-emulator/config"
	"github.com/Neumann88/payment-api-emulator/internal/controller"
	"github.com/Neumann88/payment-api-emulator/internal/repository"
	"github.com/Neumann88/payment-api-emulator/internal/usecase"
	"github.com/Neumann88/payment-api-emulator/pkg/db/migrate"
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
	logger := loggin.NewLogger().Init(cfg.Logger.Debug)

	// Database
	pg, err := postgres.NewPostgres().Connect(cfg.Postgres.Dsn)
	if err != nil {
		logger.Fatalf("postgres connection failed, %s", err.Error())
	}
	defer pg.Close()

	// Scheme migrate
	migrate.InitMigrate(logger, cfg.Postgres.Dsn)

	// Layers
	r := repository.NewRepository(pg)
	u := usecase.NewUsecase(r)
	h := controller.NewController(logger, u)

	// HTTP-Server
	httpServer := server.NewHttpServer(
		h.InitRoutes(),
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
