package postgres

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type postgres struct {
	options      DBOptions
	connAttempts int
	connTimeout  time.Duration
}

func NewPostgres(options DBOptions, connAttempts int, connTimeout time.Duration) *postgres {
	return &postgres{
		options:      options,
		connAttempts: connAttempts,
		connTimeout:  connTimeout,
	}
}

func (p *postgres) Connect() (*sql.DB, error) {
	dsn := getDSN(p.options)

	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}

	for p.connAttempts > 0 {
		err = db.PingContext(context.Background())

		if err == nil {
			break
		}

		log.Printf("Postgres is trying to connect, attempts left: %d", p.connAttempts)

		time.Sleep(p.connTimeout)

		p.connAttempts--
	}

	return db, err
}
