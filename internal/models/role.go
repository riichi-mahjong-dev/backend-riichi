package models

import "time"

type Role struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:255;not null" json:"name"`
	GuardName string    `gorm:"size:255;not null" json:"guard_name"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
