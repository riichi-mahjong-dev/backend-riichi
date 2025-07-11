package services

import (
	"time"
	"gorm.io/gorm"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/models"
)

type MatchService struct {
	BaseService
}

func NewMatchService(db *gorm.DB) *MatchService {
	return &MatchService{
		BaseService: BaseService{DB: db},
	}
}

func (s *MatchService) CreateMatch(req *models.MatchRequest) (*models.Match, error) {
	match := &models.Match{
		Player1ID:    req.Player1ID,
		Player2ID:    req.Player2ID,
		Player3ID:    req.Player3ID,
		Player4ID:    req.Player4ID,
		Player1Score: req.Player1Score,
		Player2Score: req.Player2Score,
		Player3Score: req.Player3Score,
		Player4Score: req.Player4Score,
		ParlourID:    req.ParlourID,
		CreatedBy:    req.CreatedBy,
		ApprovedBy:   req.ApprovedBy,
	}

	err := s.Create(match)
	if err != nil {
		return nil, err
	}
	return match, nil
}

func (s *MatchService) GetMatchByID(id uint64) (*models.Match, error) {
	var match models.Match
	preloads := []string{"Player1", "Player2", "Player3", "Player4", "Parlour", "Parlour.Province"}
	err := s.GetWithPreload(&match, id, preloads...)
	if err != nil {
		return nil, err
	}
	return &match, nil
}

func (s *MatchService) GetAllMatches(limit, offset int) ([]models.Match, error) {
	var matches []models.Match
	preloads := []string{"Player1", "Player2", "Player3", "Player4", "Parlour", "Parlour.Province"}
	err := s.GetAllWithPreload(&matches, limit, offset, preloads...)
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func (s *MatchService) UpdateMatch(id uint64, req *models.MatchRequest) (*models.Match, error) {
	updates := map[string]interface{}{
		"player_1_id":    req.Player1ID,
		"player_2_id":    req.Player2ID,
		"player_3_id":    req.Player3ID,
		"player_4_id":    req.Player4ID,
		"player_1_score": req.Player1Score,
		"player_2_score": req.Player2Score,
		"player_3_score": req.Player3Score,
		"player_4_score": req.Player4Score,
		"parlour_id":     req.ParlourID,
		"created_by":     req.CreatedBy,
		"approved_by":    req.ApprovedBy,
	}

	err := s.Update(&models.Match{}, id, updates)
	if err != nil {
		return nil, err
	}

	return s.GetMatchByID(id)
}

func (s *MatchService) DeleteMatch(id uint64) error {
	return s.Delete(&models.Match{}, id)
}

func (s *MatchService) ApproveMatch(id uint64, approvedBy uint64) (*models.Match, error) {
	now := time.Now()
	updates := map[string]interface{}{
		"approved_by": approvedBy,
		"approved_at": &now,
	}

	err := s.Update(&models.Match{}, id, updates)
	if err != nil {
		return nil, err
	}

	return s.GetMatchByID(id)
}

func (s *MatchService) GetMatchesByParlour(parlourID uint64, limit, offset int) ([]models.Match, error) {
	var matches []models.Match
	query := s.DB.Where("parlour_id = ?", parlourID)
	preloads := []string{"Player1", "Player2", "Player3", "Player4", "Parlour", "Parlour.Province"}
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	err := query.Find(&matches).Error
	if err != nil {
		return nil, err
	}
	return matches, nil
}
