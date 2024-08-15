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

type LogPostsModel struct {
	Id             int    `json:"id" db:"id"`
	ArticleId      int    `json:"article_id" db:"article_id"`
	DataBefore     string `json:"data_before" db:"data_before"`
	DataAfter      string `json:"data_after" db:"data_after"`
	CategoryStatus string `json:"category_status" db:"category_status"`

	CreatedDate *time.Time `json:"created_date" db:"created_date"`
	UpdatedDate *time.Time `json:"updated_date" db:"updated_date"`
}

type ResponseLogPostsModel struct {
	Id             int         `json:"id" db:"id"`
	ArticleId      int         `json:"article_id" db:"article_id"`
	DataBefore     interface{} `json:"data_before" db:"data_before"`
	DataAfter      interface{} `json:"data_after" db:"data_after"`
	CategoryStatus string      `json:"category_status" db:"category_status"`

	CreatedDate *time.Time `json:"created_date" db:"created_date"`
	UpdatedDate *time.Time `json:"updated_date" db:"updated_date"`
}

type RequestPostModel struct {
	Status   string `json:"status" validate:"required,oneof=publish draft"`
	Title    string `json:"title"  validate:"required,min=20"`
	Category string `json:"category"  validate:"required,min=3"`
	Content  string `json:"content"  validate:"required,min=200"`
}

type RequestUpdatePostModel struct {
	Id       int    `uri:"id" validate:"required,number,min=1"`
	Status   string `json:"status" validate:"required,oneof=publish draft"`
	Title    string `json:"title"  validate:"required,min=20"`
	Category string `json:"category"  validate:"required,min=3"`
	Content  string `json:"content"  validate:"required,min=200"`
}

type ResponsePostModel struct {
	Id       int    `json:"id"`
	Status   string `json:"status"`
	Title    string `json:"title"  `
	Category string `json:"category"  `
	Content  string `json:"content"  `
}

type RequestGetListPostModel struct {
	Status string `form:"status" validate:"required,oneof=publish draft thrash"`
	Search string `form:"search"`
	Limit  uint   `uri:"limit" validate:"required,number,min=1"`
	Offset uint   `uri:"offset" validate:"required,number,min=1"`
}

type RequestGetDetailPostModel struct {
	Id uint `uri:"id" validate:"required,number,min=1"`
}
