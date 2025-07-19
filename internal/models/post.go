package models

import "time"

type Post struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"size:255;not null" json:"title" validate:"required,min:2"`
	Slug      string    `gorm:"size:255;not null" json:"slug" validate:"required,min:2"`
	Content   string    `gorm:"type:text;not null" json:"content" validate:"required"`
	CreatedBy uint64    `gorm:"not null" json:"created_by" validate:"required"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type PostRequest struct {
	Title     string `json:"title" validate:"required,min:2"`
	Slug      string `json:"slug" validate:"required,min:2"`
	Content   string `json:"content" validate:"required"`
	CreatedBy uint64 `json:"created_by" validate:"required"`
}

type PostResponse struct {
	ID        uint64    `json:"id"`
	Title     string    `json:"title"`
	Slug      string    `json:"slug"`
	Content   string    `json:"content"`
	CreatedBy uint64    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
