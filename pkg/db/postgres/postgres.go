package postgres

import (
	"github.com/jmoiron/sqlx"
)

type postgres struct{}

func (p *postgres) Connect(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, err
}

func NewPostgres() *postgres {
	return &postgres{}
}
