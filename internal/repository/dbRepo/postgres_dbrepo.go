package dbrepo

import (
	"context"
	"database/sql"
	"log"
	"time"
)

const dbTimeout = time.Second * 3

type PostgresDBRepo struct {
	DB *sql.DB
}

func (m *PostgresDBRepo) Connection() *sql.DB {
	return m.DB
}

func (m *PostgresDBRepo) CreateTables() {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	createImageTableStmt := `CREATE TABLE IF NOT EXISTS public.image (
							id SERIAL PRIMARY KEY,
							url VARCHAR(255),
							created_at TIMESTAMPTZ,
							updated_at TIMESTAMPTZ
							)`

	createProductTableStmt := `CREATE TABLE IF NOT EXISTS public.product (
							id SERIAL PRIMARY KEY,
							main_category VARCHAR(255),
							sub_category VARCHAR(255),
							description TEXT,
							url VARCHAR(255),
							order_weight INT,
							score REAL,
							main_image INT REFERENCES public.image(ID),
							created_at TIMESTAMPTZ,
							updated_at TIMESTAMPTZ
							)`

	createPageTableStmt := `CREATE TABLE IF NOT EXISTS public.page (
							id SERIAL PRIMARY KEY,
							title VARCHAR(255) NOT NULL,
							intro TEXT,
							body TEXT,
							products INT REFERENCES public.product(ID),
							slug VARCHAR(255) UNIQUE NOT NULL
						)`

	_, err := m.DB.ExecContext(ctx, createImageTableStmt)
	_, err = m.DB.ExecContext(ctx, createProductTableStmt)
	_, err = m.DB.ExecContext(ctx, createPageTableStmt)

	if err != nil {
		log.Fatal("[SQL ERROR]: Cannot create table", err)
	}
}
