package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type PostgresStorage struct{
	db *sql.DB
}

func NewPostgresStorage (connStr string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &PostgresStorage{db: db}, nil
}

func (ps *PostgresStorage) Close() error {
    return ps.db.Close()
}