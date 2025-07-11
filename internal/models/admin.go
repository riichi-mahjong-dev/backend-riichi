package models

import "time"

type AdminRole string

const (
	AdminRoleSuperAdmin AdminRole = "super-admin"
	AdminRoleStaff     AdminRole = "staff"
)

type Admin struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"size:255;not null" json:"username" validate:"required,min:2"`
	Password  string    `gorm:"size:255;not null" json:"-"`
	Role      AdminRole `gorm:"size:20;not null" json:"role" validate:"required"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type AdminRequest struct {
	Username string    `json:"username" validate:"required,min:2"`
	Password string    `json:"password" validate:"required,min:6"`
	Role     AdminRole `json:"role" validate:"required"`
}

type AdminResponse struct {
	ID        uint64    `json:"id"`
	Username  string    `json:"username"`
	Role      AdminRole `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
