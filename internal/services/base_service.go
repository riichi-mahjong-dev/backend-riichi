package services

import (
	"fmt"
	"strings"

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

func Paginate[T any](
	db *gorm.DB,
	model T,
	joins []string,
	filters map[string]any,
	includes []string,
	searchFields []string,
	searchTerm string,
	page int,
	pageSize int,
	orderBy string,
	order string,
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

	for _, join := range joins {
		query = query.Joins(join)
	}

	// Apply search
	if searchTerm != "" && len(searchFields) > 0 {
		var likeConditions []string
		var likeArgs []any
		for _, col := range searchFields {
			likeConditions = append(likeConditions, fmt.Sprintf("%s ILIKE ?", col))
			likeArgs = append(likeArgs, "%"+searchTerm+"%")
		}
		query = query.Where(strings.Join(likeConditions, " OR "), likeArgs...)
	}

	// Apply filters
	for field, value := range filters {
		query = query.Where(fmt.Sprintf("%s = ?", field), value)
	}

	for _, include := range includes {
		query = query.Preload(include)
	}

	query = query.Order(orderBy + " " + order)

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
