package services

import (
	"github.com/riichi-mahjong-dev/backend-riichi/internal/models"
	"gorm.io/gorm"
)

type RoleService struct {
	BaseService
}

func NewRoleService(db *gorm.DB) *RoleService {
	return &RoleService{
		BaseService: BaseService{DB: db},
	}
}

func (s *RoleService) CreateRole(req *models.RoleRequest) (*models.Role, error) {
	role := &models.Role{
		Name:      req.Name,
		GuardName: req.GuardName,
	}

	err := s.Create(role)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (s *RoleService) GetRoleByID(id uint64) (*models.Role, error) {
	var role models.Role
	err := s.GetByID(&role, id)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (s *RoleService) GetAllRoles(limit, offset int) ([]models.Role, error) {
	var roles []models.Role
	err := s.GetAll(&roles, limit, offset)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (s *RoleService) UpdateRole(id uint64, req *models.RoleRequest) (*models.Role, error) {
	updates := map[string]any{
		"name":       req.Name,
		"guard_name": req.GuardName,
	}

	err := s.Update(&models.Role{}, id, updates)
	if err != nil {
		return nil, err
	}

	return s.GetRoleByID(id)
}

func (s *RoleService) DeleteRole(id uint64) error {
	return s.Delete(&models.Role{}, id)
}

func (s *RoleService) GetRoleByName(name string) (*models.Role, error) {
	var role models.Role
	err := s.DB.Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}
