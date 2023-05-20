package dbrepo

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"
	"topio/internal/models"
	ai "topio/openAI"
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

	//createImageTableStmt := `CREATE TABLE IF NOT EXISTS public.image (
	//						id SERIAL PRIMARY KEY,
	//						url VARCHAR(255) NOT NULL,
	//						created_at TIMESTAMPTZ,
	//						updated_at TIMESTAMPTZ
	//						)`

	createPageTableStmt := `CREATE TABLE IF NOT EXISTS public.page (
							id SERIAL PRIMARY KEY,
                            entry_data_url VARCHAR(255),
							title VARCHAR(255) NOT NULL,
							product_name VARCHAR(255) NOT NULL,
							intro TEXT,
							body TEXT,
							slug VARCHAR(255) UNIQUE NOT NULL,
							parent_page INT REFERENCES page(id),
                            created_at TIMESTAMPTZ,
							updated_at TIMESTAMPTZ
						)`

	createProductTableStmt := `CREATE TABLE IF NOT EXISTS public.product (
							id SERIAL PRIMARY KEY,
							name VARCHAR(255),
							description TEXT,
							url VARCHAR(255),
							order_weight INT,
							score REAL,
							main_image TEXT,
							page_id INT REFERENCES public.page(id),
							features JSONB,
							created_at TIMESTAMPTZ,
							updated_at TIMESTAMPTZ
							)`

	//_, err := m.DB.ExecContext(ctx, createImageTableStmt)
	_, err := m.DB.ExecContext(ctx, createPageTableStmt)
	_, err = m.DB.ExecContext(ctx, createProductTableStmt)

	if err != nil {
		log.Fatal("[SQL ERROR]: Cannot create table", err)
	}
}

func (m *PostgresDBRepo) InsertProduct(product ai.ProductWithGeneratedDescription) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `insert into product (name,description,url,order_weight,score,main_image,page_id,features,created_at,updated_at)
			values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) returning id`

	var newID int

	features, err := json.Marshal(product.Features)
	if err != nil {
		fmt.Println("Failed to marshal JSON:", err)
		features = make([]byte, 0)
	}

	err = m.DB.QueryRowContext(ctx, stmt,
		product.Name,
		product.GeneratedDescription,
		product.Url,
		product.OrderWeight,
		0,
		nil,
		product.RelatedPage,
		features,
		time.Now().UTC(),
		time.Now().UTC(),
	).Scan(&newID)

	if err != nil {
		fmt.Printf("[InsertPage page]: %v\n", err)
		return 0, err
	}

	return newID, nil
}

func (m *PostgresDBRepo) InsertPage(page models.Page) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `insert into page (title,product_name,intro,body,slug,entry_data_url,parent_page,created_at,updated_at)
			values ($1,$2,$3,$4,$5,$6,$7,$8,$9) returning id`

	var newID int
	err := m.DB.QueryRowContext(ctx, stmt,
		page.Title,
		page.ProductName,
		page.Intro,
		page.Body,
		page.Slug,
		page.EntryDataUrl,
		&page.ParentPage,
		time.Now().UTC(),
		time.Now().UTC(),
	).Scan(&newID)

	if err != nil {
		fmt.Printf("[InsertPage page]: %v\n", err)
		return 0, err
	}

	return newID, nil
}
