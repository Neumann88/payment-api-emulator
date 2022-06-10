package migrate

import (
	"github.com/Neumann88/payment-api-emulator/pkg/loggin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func InitMigrate(logger loggin.ILogger, dsn string) {
	if len(dsn) == 0 {
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
