package postgres

import (
	"database/sql"
)

type postgres struct {
	options DBOptions
}

func NewPostgres(options DBOptions) *postgres {
	return &postgres{
		options: options,
	}
}

func (p *postgres) Connect() (*sql.DB, error) {
	dsn := getDSN(p.options)

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
