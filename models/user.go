package models

type User struct {
	Id          int    `json:"user_id"`
	Name        string `json:"user_name"`
	Email       string `json:"user_email"`
	Description string `json:"user_description"`
	Avatar      string `json:"user_avatar"`
}
