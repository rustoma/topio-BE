package models

import (
	"time"
	"topio/scrapper"
)

type Page struct {
	ID           int                `json:"id"`
	Title        string             `json:"title"`
	ProductName  string             `json:"product_name"`
	Intro        string             `json:"intro"`
	Body         string             `json:"body"`
	Slug         string             `json:"slug"`
	EntryDataUrl string             `json:"entry_data_url"`
	ParentPage   *int               `json:"parent_page"`
	Features     []scrapper.Feature `json:"features"`
	CreatedAt    time.Time          `json:"-"`
	UpdatedAt    time.Time          `json:"-"`
}
