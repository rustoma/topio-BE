package models

import "time"

type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Url         string    `json:"url"`
	OrderWeight int       `json:"order_weight"`
	Score       float32   `json:"score"`
	MainImage   Image     `json:"main_image"`
	Page        Page      `json:"page"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

type Image struct {
	ID        int       `json:"id"`
	Url       string    `json:"url"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
