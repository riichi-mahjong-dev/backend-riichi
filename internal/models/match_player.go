package models

import "time"

type MatchPlayer struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	MatchID   uint64    `gorm:"not null" json:"match_id" validate:"required"`
	PlayerID  uint64    `gorm:"not null" json:"player_id" validate:"required"`
	Point     *int      `json:"point"`
	MmrDelta  *int      `json:"mmr_delta"`
	Pinalty   *int      `json:"pinalty"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relations
	Match  Match  `gorm:"foreignKey:MatchID" json:"-"`
	Player Player `gorm:"foreignKey:PlayerID" json:"player,omitempty"`
}
