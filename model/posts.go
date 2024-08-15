package model

import "time"

type PostsModel struct {
	Id          int        `json:"id" db:"id"`
	Title       string     `json:"title" db:"title"`
	Content     string     `json:"content" db:"content"`
	Category    string     `json:"category" db:"category"`
	Status      string     `json:"status" db:"status"`
	CreatedDate *time.Time `json:"created_date" db:"created_date"`
	UpdatedDate *time.Time `json:"updated_date" db:"updated_date"`
}

type RequestPostModel struct {
	Status   string `json:"status" validate:"required,oneof=publish draft thrash"`
	Title    string `json:"title"  validate:"required,min=20"`
	Category string `json:"category"  validate:"required,min=3"`
	Content  string `json:"content"  validate:"required,min=200"`
}
