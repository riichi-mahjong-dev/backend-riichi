package models

import "time"

type Match struct {
	ID           uint64     `gorm:"primaryKey" json:"id"`
	Player1ID    *uint64    `json:"player_1_id"`
	Player2ID    *uint64    `json:"player_2_id"`
	Player3ID    *uint64    `json:"player_3_id"`
	Player4ID    *uint64    `json:"player_4_id"`
	Player1Score *int       `json:"player_1_score"`
	Player2Score *int       `json:"player_2_score"`
	Player3Score *int       `json:"player_3_score"`
	Player4Score *int       `json:"player_4_score"`
	ParlourID    uint64     `gorm:"not null" json:"parlour_id" validate:"required"`
	CreatedBy    uint64     `gorm:"not null" json:"created_by" validate:"required"`
	ApprovedBy   *uint64    `json:"approved_by"`
	ApprovedAt   *time.Time `json:"approved_at"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"autoUpdateTime"`

	// Relations
	Player1  *Player `gorm:"foreignKey:Player1ID" json:"player_1,omitempty"`
	Player2  *Player `gorm:"foreignKey:Player2ID" json:"player_2,omitempty"`
	Player3  *Player `gorm:"foreignKey:Player3ID" json:"player_3,omitempty"`
	Player4  *Player `gorm:"foreignKey:Player4ID" json:"player_4,omitempty"`
	Parlour  Parlour `gorm:"foreignKey:ParlourID" json:"parlour,omitempty"`
	Creator  Player  `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Approver *Admin  `gorm:"foreignKey:ApprovedBy" json:"approved_by,omitempty"`
}

type MatchRequest struct {
	Player1ID *uint64 `json:"player_1_id"`
	Player2ID *uint64 `json:"player_2_id"`
	Player3ID *uint64 `json:"player_3_id"`
	Player4ID *uint64 `json:"player_4_id"`
	ParlourID uint64  `json:"parlour_id" validate:"required"`
	CreatedBy uint64  `json:"created_by" validate:"required"`
}

type MatchResponse struct {
	ID           uint64           `json:"id"`
	Player1ID    *uint64          `json:"player_1_id"`
	Player2ID    *uint64          `json:"player_2_id"`
	Player3ID    *uint64          `json:"player_3_id"`
	Player4ID    *uint64          `json:"player_4_id"`
	Player1Score *int             `json:"player_1_score"`
	Player2Score *int             `json:"player_2_score"`
	Player3Score *int             `json:"player_3_score"`
	Player4Score *int             `json:"player_4_score"`
	ParlourID    uint64           `json:"parlour_id"`
	CreatedBy    uint64           `json:"created_by"`
	ApprovedBy   *uint64          `json:"approved_by"`
	ApprovedAt   *time.Time       `json:"approved_at"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
	Player1      *PlayerResponse  `json:"player_1,omitempty"`
	Player2      *PlayerResponse  `json:"player_2,omitempty"`
	Player3      *PlayerResponse  `json:"player_3,omitempty"`
	Player4      *PlayerResponse  `json:"player_4,omitempty"`
	Parlour      *ParlourResponse `json:"parlour,omitempty"`
}

type PointMatchRequest struct {
	Player1Score *int `json:"player_1_score"`
	Player2Score *int `json:"player_2_score"`
	Player3Score *int `json:"player_3_score"`
	Player4Score *int `json:"player_4_score"`
}
