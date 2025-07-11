package models

import "time"

type Role struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:255;not null" json:"name"`
	GuardName string    `gorm:"size:255;not null" json:"guard_name"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type RoleRequest struct {
	Name      string `json:"name" validate:"required,min:2"`
	GuardName string `json:"guard_name" validate:"required"`
}

type RoleResponse struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	GuardName string    `json:"guard_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
