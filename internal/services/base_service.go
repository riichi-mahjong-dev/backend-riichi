package services

import (
	"github.com/riichi-mahjong-dev/backend-riichi/commons"
	"gorm.io/gorm"
)

type BaseService struct {
	DB *gorm.DB
}

// Create a new instance with a transaction DB
func (s *BaseService) WithTx(tx *gorm.DB) *BaseService {
	return &BaseService{DB: tx}
}

// Generic CRUD operations
func (s *BaseService) Create(model any) error {
	return s.DB.Create(model).Error
}

func (s *BaseService) GetByID(model any, id uint64) error {
	return s.DB.First(model, id).Error
}

func (s *BaseService) GetAll(models any, limit, offset int) error {
	query := s.DB
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	return query.Find(models).Error
}

func (s *BaseService) Update(model any, id uint64, updates any) error {
	return s.DB.Model(model).Where("id = ?", id).Updates(updates).Error
}

func (s *BaseService) Delete(model any, id uint64) error {
	return s.DB.Delete(model, id).Error
}

func (s *BaseService) GetWithPreload(model any, id uint64, preloads ...string) error {
	query := s.DB
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	return query.First(model, id).Error
}

func (s *BaseService) GetAllWithPreload(models any, queryPaginate commons.QueryPagination, preloads ...string) error {
	query := s.DB
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	if queryPaginate.Limit > 0 {
		query = query.Limit(queryPaginate.Limit)
	}
	if queryPaginate.Offset > 0 {
		query = query.Offset(queryPaginate.Offset)
	}
	return query.Find(models).Error
}

func (s *BaseService) Count(model any) (int64, error) {
	var count int64
	err := s.DB.Model(model).Count(&count).Error
	return count, err
}

func (s *BaseService) Exists(model any, id uint64) (bool, error) {
	var count int64
	err := s.DB.Model(model).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

func (s *BaseService) Paginate[T any](
	db *gorm.DB,
	model T,
	filters map[string]interface{},
	searchFields []string,
	searchTerm string,
	page int,
	pageSize int,
) ([]T, int64, error) {
	var results []T
	var total int64

	// Default pagination
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	query := db.Model(model)

	// Apply filters
	for field, value := range filters {
		query = query.Where(fmt.Sprintf("%s = ?", field), value)
	}

	// Apply search
	if searchTerm != "" && len(searchFields) > 0 {
		var conditions []string
		var args []interface{}
		for _, field := range searchFields {
			conditions = append(conditions, fmt.Sprintf("%s LIKE ?", field))
			args = append(args, "%"+searchTerm+"%")
		}
		query = query.Where(strings.Join(conditions, " OR "), args...)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if err := query.Limit(pageSize).Offset(offset).Find(&results).Error; err != nil {
		return nil, 0, err
	}

	return results, total, nil
}
