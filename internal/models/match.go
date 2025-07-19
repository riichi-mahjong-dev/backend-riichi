package models

import "time"

type Match struct {
	ApprovedBy *uint64    `json:"approved_by"`
	ApprovedAt *time.Time `json:"approved_at"`
	Approver   *Admin     `gorm:"foreignKey:approved_by" json:"approved_by,omitempty"`
	ID         uint64     `gorm:"primaryKey" json:"id"`
	ParlourID  uint64     `gorm:"not null" json:"parlour_id" validate:"required"`
	CreatedBy  *uint64    `json:"created_by"`
	CreatedAt  time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time  `json:"updated_at" gorm:"autoUpdateTime"`

	// Relations
	Players []MatchPlayer `gorm:"foreignKey:match_id"`
	Parlour Parlour       `gorm:"foreignKey:parlour_id" json:"parlour,omitempty"`
	Creator *Player       `gorm:"foreignKey:created_by" json:"creator,omitempty"`
}

type PointMatchPlayer struct {
	MatchPlayerId *uint64 `json:"match_player_id"`
	Score         *int    `json:"score"`
}

type PlayerMatch struct {
	Player *uint64 `json:"player"`
}

type UpdatePlayerMatch struct {
	Player        *uint64 `json:"player"`
	MatchPlayerID *uint64 `json:"match_player_id"`
}

type UpdateMatchRequest struct {
	ParlourID uint64              `json:"parlour_id" validate:"required"`
	Players   []UpdatePlayerMatch `json:"players"`
}

type MatchRequest struct {
	ParlourID uint64        `json:"parlour_id" validate:"required"`
	Players   []PlayerMatch `json:"players"`
}

type MatchResponse struct {
	ID           uint64           `json:"id"`
	ParlourID    uint64           `json:"parlour_id"`
	CreatedBy    uint64           `json:"created_by"`
	ApprovedBy   *uint64          `json:"approved_by,omitempty"`
	ApprovedAt   *time.Time       `json:"approved_at,omitempty"`
	Parlour      *ParlourResponse `json:"parlour,omitempty"`
	Approver     *Admin           `json:"approver,omitempty"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
	MatchPlayers []MatchPlayer    `json:"match_players"`
	Creator      *Player          `json:"creator_player"`
}

type PointMatchRequest struct {
	PointMatchPlayers []PointMatchPlayer `json:"point_match_players"`
}
