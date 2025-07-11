package models

import "time"

type Player struct {
	ID         uint64    `gorm:"primaryKey" json:"id"`
	ProvinceID uint64    `gorm:"not null" json:"province_id"`
	Rank       int       `gorm:"not null" json:"rank"`
	Name       string    `gorm:"size:255;not null" json:"name"`
	Country    string    `gorm:"size:255;not null" json:"country"`
	Username   string    `gorm:"size:255;not null" json:"username" validate:"required,min:2"`
	Password   string    `gorm:"size:255;not null" json:"-"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relations
	Province Province `gorm:"foreignKey:ProvinceID" json:"province,omitempty"`
}

type PlayerRequest struct {
	ProvinceID uint64 `json:"province_id" validate:"required"`
	Rank       int    `json:"rank" validate:"required"`
	Name       string `json:"name" validate:"required,min:2"`
	Country    string `json:"country" validate:"required"`
	Username   string `json:"username" validate:"required,min:2"`
	Password   string `json:"password" validate:"required,min:6"`
}

type PlayerResponse struct {
	ID         uint64            `json:"id"`
	ProvinceID uint64            `json:"province_id"`
	Rank       int               `json:"rank"`
	Name       string            `json:"name"`
	Country    string            `json:"country"`
	Username   string            `json:"username"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
	Province   *ProvinceResponse `json:"province,omitempty"`
}
