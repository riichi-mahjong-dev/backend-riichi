package models

import "time"

type Player struct {
	ID         uint64    `gorm:"primaryKey" json:"id"`
	ProvinceId uint64    `gorm:"not null" json:"province_id"`
	Rank       int       `gorm:"not null" json:"rank"`
	Name       string    `gorm:"size:255;not null" json:"name"`
	Country    string    `gorm:"size:255;not null" json:"country"`
	Username   string    `gorm:"size:255;not null" json:"username" validate:"required,min:2"`
	Password   string    `gorm:"size:36;not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
