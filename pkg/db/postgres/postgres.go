package postgres

import (
	"database/sql"
)

type postgres struct{}

func NewPostgres() *postgres {
	return &postgres{}
}

func (p *postgres) Connect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, err
}
