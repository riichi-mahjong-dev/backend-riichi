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
