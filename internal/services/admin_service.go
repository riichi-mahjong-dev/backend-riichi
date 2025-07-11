package services

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/models"
)

type AdminService struct {
	BaseService
}

func NewAdminService(db *gorm.DB) *AdminService {
	return &AdminService{
		BaseService: BaseService{DB: db},
	}
}

func (s *AdminService) CreateAdmin(req *models.AdminRequest) (*models.Admin, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	admin := &models.Admin{
		Username: req.Username,
		Password: string(hashedPassword),
		Role:     req.Role,
	}

	err = s.Create(admin)
	if err != nil {
		return nil, err
	}
	return admin, nil
}

func (s *AdminService) GetAdminByID(id uint64) (*models.Admin, error) {
	var admin models.Admin
	err := s.GetByID(&admin, id)
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (s *AdminService) GetAllAdmins(limit, offset int) ([]models.Admin, error) {
	var admins []models.Admin
	err := s.GetAll(&admins, limit, offset)
	if err != nil {
		return nil, err
	}
	return admins, nil
}

func (s *AdminService) UpdateAdmin(id uint64, req *models.AdminRequest) (*models.Admin, error) {
	updates := map[string]interface{}{
		"username": req.Username,
		"role":     req.Role,
	}

	// Hash password if provided
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		updates["password"] = string(hashedPassword)
	}

	err := s.Update(&models.Admin{}, id, updates)
	if err != nil {
		return nil, err
	}

	return s.GetAdminByID(id)
}

func (s *AdminService) DeleteAdmin(id uint64) error {
	return s.Delete(&models.Admin{}, id)
}

func (s *AdminService) GetAdminByUsername(username string) (*models.Admin, error) {
	var admin models.Admin
	err := s.DB.Where("username = ?", username).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}
