package services

import (
	"github.com/riichi-mahjong-dev/backend-riichi/commons"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/models"
	"gorm.io/gorm"
)

type ParlourService struct {
	BaseService
}

func NewParlourService(db *gorm.DB) *ParlourService {
	return &ParlourService{
		BaseService: BaseService{DB: db},
	}
}

func (s *ParlourService) CreateParlour(req *models.ParlourRequest) (*models.Parlour, error) {
	parlour := &models.Parlour{
		Name:       req.Name,
		Country:    req.Country,
		ProvinceID: req.ProvinceID,
		Address:    req.Address,
	}

	err := s.Create(parlour)
	if err != nil {
		return nil, err
	}
	return parlour, nil
}

func (s *ParlourService) GetParlourByID(id uint64) (*models.Parlour, error) {
	var parlour models.Parlour
	err := s.GetWithPreload(&parlour, id, "Province")
	if err != nil {
		return nil, err
	}
	return &parlour, nil
}

func (s *ParlourService) GetAllParlours(queryPaginate commons.QueryPagination) ([]models.Parlour, error) {
	var parlours []models.Parlour
	err := s.GetAllWithPreload(&parlours, queryPaginate, "Province")
	if err != nil {
		return nil, err
	}
	return parlours, nil
}

func (s *ParlourService) UpdateParlour(id uint64, req *models.ParlourRequest) (*models.Parlour, error) {
	updates := map[string]any{
		"name":        req.Name,
		"country":     req.Country,
		"province_id": req.ProvinceID,
		"address":     req.Address,
	}

	err := s.Update(&models.Parlour{}, id, updates)
	if err != nil {
		return nil, err
	}

	return s.GetParlourByID(id)
}

func (s *ParlourService) DeleteParlour(id uint64) error {
	return s.Delete(&models.Parlour{}, id)
}
