package dbrepo

import "database/sql"

type PostgresDBRepo struct {
	DB *sql.DB
}

func (m *PostgresDBRepo) Connection() *sql.DB {
	return m.DB
}
