package models

import "time"

type UserType string

const (
	UserTypePlayer UserType = "player"
	UserTypeAdmin  UserType = "admin"
)

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresAt    time.Time `json:"expires_at"`
	User         interface{} `json:"user"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type TokenClaims struct {
	UserID   uint64   `json:"user_id"`
	Username string   `json:"username"`
	UserType UserType `json:"user_type"`
	Role     string   `json:"role,omitempty"` // For admin roles
}

type AuthUser struct {
	ID       uint64   `json:"id"`
	Username string   `json:"username"`
	UserType UserType `json:"user_type"`
	Role     string   `json:"role,omitempty"`
}
