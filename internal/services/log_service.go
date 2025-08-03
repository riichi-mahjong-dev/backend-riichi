package services

import (
	"strings"

	"github.com/riichi-mahjong-dev/backend-riichi/commons"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/models"
	"gorm.io/gorm"
)

type LogService struct {
	BaseService
}

func NewLogService(db *gorm.DB) *LogService {
	return &LogService{
		BaseService: BaseService{DB: db},
	}
}

func (s *LogService) GetLogById(id uint64) (*models.Job, error) {
	var job models.Job
	err := s.GetByID(&job, id)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (s *LogService) GetAllLogs(queryPaginate commons.QueryParams) ([]models.Job, int64, error) {
	preloads := map[string]func(*gorm.DB) *gorm.DB{}
	return Paginate(
		s.DB,
		models.Job{},
		[]string{},
		queryPaginate.Filters,
		map[string]func(*gorm.DB, any) *gorm.DB{
			"created_between": func(d *gorm.DB, a any) *gorm.DB {
				val, ok := a.([]string)
				if !ok {
					return d
				}

				if len(val) != 2 {
					return d
				}

				startDate := strings.TrimSpace(val[0])
				endDate := strings.TrimSpace(val[1])
				return d.Where("created_at BETWEEN ? AND ?", startDate, endDate)
			},
		},
		preloads,
		[]string{"job_type", "status", "reason"},
		queryPaginate.Search,
		queryPaginate.Page,
		queryPaginate.PageSize,
		queryPaginate.OrderBy,
		queryPaginate.Order,
	)
}
