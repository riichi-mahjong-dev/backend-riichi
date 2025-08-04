package services

import (
	"github.com/riichi-mahjong-dev/backend-riichi/commons"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/models"
	"gorm.io/gorm"
)

type ProvinceService struct {
	BaseService
}

func NewProvinceService(db *gorm.DB) *ProvinceService {
	return &ProvinceService{
		BaseService: BaseService{DB: db},
	}
}

func (s *ProvinceService) CreateProvince(req *models.ProvinceRequest) (*models.Province, error) {
	province := &models.Province{
		Name: req.Name,
	}

	err := s.Create(province)
	if err != nil {
		return nil, err
	}
	return province, nil
}

func (s *ProvinceService) GetProvinceByID(id uint64) (*models.Province, error) {
	var province models.Province
	err := s.GetByID(&province, id)
	if err != nil {
		return nil, err
	}
	return &province, nil
}

func (s *ProvinceService) GetAllProvinces(queryPaginate commons.QueryParams) ([]models.Province, int64, error) {
	preloads := map[string]func(*gorm.DB) *gorm.DB{}
	return Paginate(
		s.DB,
		models.Province{},
		[]string{},
		queryPaginate.Filters,
		map[string]func(*gorm.DB, any) *gorm.DB{},
		preloads,
		[]string{"name"},
		queryPaginate,
	)
}

func (s *ProvinceService) UpdateProvince(id uint64, req *models.ProvinceRequest) (*models.Province, error) {
	updates := map[string]any{
		"name": req.Name,
	}

	err := s.Update(&models.Province{}, id, updates)
	if err != nil {
		return nil, err
	}

	return s.GetProvinceByID(id)
}

func (s *ProvinceService) DeleteProvince(id uint64) error {
	return s.Delete(&models.Province{}, id)
}

func (s *ProvinceService) GetProvinceByName(name string) (*models.Province, error) {
	var province models.Province
	err := s.DB.Where("name = ?", name).First(&province).Error
	if err != nil {
		return nil, err
	}
	return &province, nil
}
