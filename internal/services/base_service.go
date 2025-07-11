package services

import (
	"gorm.io/gorm"
)

type BaseService struct {
	DB *gorm.DB
}

// Generic CRUD operations
func (s *BaseService) Create(model interface{}) error {
	return s.DB.Create(model).Error
}

func (s *BaseService) GetByID(model interface{}, id uint64) error {
	return s.DB.First(model, id).Error
}

func (s *BaseService) GetAll(models interface{}, limit, offset int) error {
	query := s.DB
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	return query.Find(models).Error
}

func (s *BaseService) Update(model interface{}, id uint64, updates interface{}) error {
	return s.DB.Model(model).Where("id = ?", id).Updates(updates).Error
}

func (s *BaseService) Delete(model interface{}, id uint64) error {
	return s.DB.Delete(model, id).Error
}

func (s *BaseService) GetWithPreload(model interface{}, id uint64, preloads ...string) error {
	query := s.DB
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	return query.First(model, id).Error
}

func (s *BaseService) GetAllWithPreload(models interface{}, limit, offset int, preloads ...string) error {
	query := s.DB
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	return query.Find(models).Error
}

func (s *BaseService) Count(model interface{}) (int64, error) {
	var count int64
	err := s.DB.Model(model).Count(&count).Error
	return count, err
}

func (s *BaseService) Exists(model interface{}, id uint64) (bool, error) {
	var count int64
	err := s.DB.Model(model).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}
