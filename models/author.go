package models

type Author struct {
	Id          int    `json:"author_id"`
	Name        string `json:"author_name"`
	Description string `json:"author_description"`
	Image       string `json:"author_image"`
	AuthorURL   string `json:"author_url"`
}
