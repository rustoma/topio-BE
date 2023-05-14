package models

import "time"

type Product struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	MainCategory string    `json:"main_category"`
	SubCategory  string    `json:"sub_category"`
	Description  string    `json:"description"`
	Url          string    `json:"url"`
	OrderWeight  int       `json:"order_weight"`
	Score        float32   `json:"score"`
	MainImage    Image     `json:"main_image"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
}

type Image struct {
	ID        int       `json:"id"`
	Url       string    `json:"url"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
