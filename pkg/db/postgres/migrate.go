package postgres

import (
	"github.com/golang-migrate/migrate/v4"
	// need for migrate
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	// need for migrate
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/Neumann88/payment-api-emulator/pkg/loggin"
)

func InitMigrate(logger loggin.ILogger, options DBOptions) {
	dsn := getDSN(options)

	if dsn == "" {
		logger.Fatal("migrate: environment variable not declared")
	}

	m, err := migrate.New("file://migrations", dsn)

	if err != nil {
		logger.Fatal(err)
	}

	if err := m.Up(); err != nil {
		logger.Debug(err)
	}

	m.Close()
}
