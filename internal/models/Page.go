package models

import "time"

type Page struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	Intro        string    `json:"intro"`
	Body         string    `json:"body"`
	Slug         string    `json:"slug"`
	EntryDataUrl string    `json:"entry_data_url"`
	ParentPage   *int      `json:"parent_page"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
}
