package services

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/models"
)

type PlayerService struct {
	BaseService
}

func NewPlayerService(db *gorm.DB) *PlayerService {
	return &PlayerService{
		BaseService: BaseService{DB: db},
	}
}

func (s *PlayerService) CreatePlayer(req *models.PlayerRequest) (*models.Player, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	player := &models.Player{
		ProvinceID: req.ProvinceID,
		Rank:       req.Rank,
		Name:       req.Name,
		Country:    req.Country,
		Username:   req.Username,
		Password:   string(hashedPassword),
	}

	err = s.Create(player)
	if err != nil {
		return nil, err
	}
	return player, nil
}

func (s *PlayerService) GetPlayerByID(id uint64) (*models.Player, error) {
	var player models.Player
	err := s.GetWithPreload(&player, id, "Province")
	if err != nil {
		return nil, err
	}
	return &player, nil
}

func (s *PlayerService) GetAllPlayers(limit, offset int) ([]models.Player, error) {
	var players []models.Player
	err := s.GetAllWithPreload(&players, limit, offset, "Province")
	if err != nil {
		return nil, err
	}
	return players, nil
}

func (s *PlayerService) UpdatePlayer(id uint64, req *models.PlayerRequest) (*models.Player, error) {
	updates := map[string]interface{}{
		"province_id": req.ProvinceID,
		"rank":        req.Rank,
		"name":        req.Name,
		"country":     req.Country,
		"username":    req.Username,
	}

	// Hash password if provided
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		updates["password"] = string(hashedPassword)
	}

	err := s.Update(&models.Player{}, id, updates)
	if err != nil {
		return nil, err
	}

	return s.GetPlayerByID(id)
}

func (s *PlayerService) DeletePlayer(id uint64) error {
	return s.Delete(&models.Player{}, id)
}

func (s *PlayerService) GetPlayerByUsername(username string) (*models.Player, error) {
	var player models.Player
	err := s.DB.Where("username = ?", username).First(&player).Error
	if err != nil {
		return nil, err
	}
	return &player, nil
}
