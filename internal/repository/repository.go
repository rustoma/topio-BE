package repository

import (
	"database/sql"
	"topio/internal/models"
	ai "topio/openAI"
)

type DatabaseRepo interface {
	Connection() *sql.DB
	CreateTables()
	InsertProduct(product ai.ProductWithGeneratedDescription) (int, error)
	InsertPage(page models.Page) (int, error)
}
