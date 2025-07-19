package jobs

import (
	"context"
	"log"
	"time"

	"github.com/riichi-mahjong-dev/backend-riichi/internal/models"
	"gorm.io/gorm"
)

func HandleMMR(_ context.Context, db *gorm.DB, job models.Job) error {
	log.Printf("Calculate MMR %d", job.ID)
	time.Sleep(2 * time.Second)
	return nil
}
