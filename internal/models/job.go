package models

import (
	"encoding/json"
	"time"
)

type Job struct {
	ID        uint64          `gorm:"primaryKey" json:"id"`
	JobType   string          `json:"job_type"`
	Payload   json.RawMessage `json:"payload" gorm:"type:json"`
	Status    string          `json:"status" gorm:"type:enum('queued', 'processing', 'done', 'error')"`
	Reason    *string         `json:"reason"`
	CreatedAt time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
}
