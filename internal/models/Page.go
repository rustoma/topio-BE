package models

type Page struct {
	ID       int      `json:"id"`
	Title    string   `json:"title"`
	Intro    string   `json:"intro"`
	Body     string   `json:"body"`
	Products *Product `json:"products"`
	Slug     string   `json:"slug"`
}
