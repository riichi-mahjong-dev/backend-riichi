package jobs

import (
	"encoding/json"
	"time"

	"github.com/riichi-mahjong-dev/backend-riichi/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func EnqueueJob(db *gorm.DB, jobType string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	job := models.Job{
		JobType: jobType,
		Payload: data,
		Status:  "queued",
	}

	return db.Create(&job).Error
}

func FetchPendingJob(db *gorm.DB) (*models.Job, error) {
	var job models.Job

	tx := db.Begin()
	if err := tx.Clauses(
		clause.Locking{Strength: "UPDATE"},
	).Where("status = ?", "queued").Order("created_at").First(&job).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Model(&job).Where("id = ?", job.ID).Update("status", "processing").Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &job, tx.Commit().Error
}

func MarkJobDone(db *gorm.DB, id uint64) error {
	return db.Model(&models.Job{}).Where("id = ?", id).Updates(map[string]any{
		"status":     "done",
		"updated_at": time.Now(),
	}).Error
}

func MarkJobFailed(db *gorm.DB, id uint64, reason string) error {
	return db.Model(&models.Job{}).Where("id = ?", id).Updates(map[string]any{
		"status":     "error",
		"reason":     reason,
		"updated_at": time.Now(),
	}).Error
}
