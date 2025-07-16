package services

import (
	"fmt"
	"time"

	"github.com/riichi-mahjong-dev/backend-riichi/commons"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/models"
	"gorm.io/gorm"
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
		Player1ID: req.Player1ID,
		Player2ID: req.Player2ID,
		Player3ID: req.Player3ID,
		Player4ID: req.Player4ID,
		ParlourID: req.ParlourID,
		CreatedBy: req.CreatedBy,
	}

	err := s.Create(match)
	if err != nil {
		return nil, err
	}
	return match, nil
}

func (s *MatchService) PointMatch(id uint64, req *models.PointMatchRequest) (*models.Match, error) {
	updates := map[string]any{
		"player_1_score": req.Player1Score,
		"player_2_score": req.Player2Score,
		"player_3_score": req.Player3Score,
		"player_4_score": req.Player4Score,
	}

	err := s.Update(&models.Match{}, id, updates)
	if err != nil {
		return nil, err
	}

	return s.GetMatchByID(id)
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

func (s *MatchService) GetAllMatches(queryPaginate commons.QueryPagination) ([]models.Match, error) {
	var matches []models.Match
	preloads := []string{"Player1", "Player2", "Player3", "Player4", "Parlour", "Parlour.Province"}
	err := s.GetAllWithPreload(&matches, queryPaginate, preloads...)
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func (s *MatchService) UpdateMatch(id uint64, req *models.MatchRequest) (*models.Match, error) {
	match, err := s.GetMatchByID(id)

	if err != nil {
		return nil, err
	}

	if match.ApprovedAt != nil || match.ApprovedBy != nil {
		return nil, fmt.Errorf("match is already approved, cannot changed anymore")
	}

	updates := map[string]any{
		"player_1_id": req.Player1ID,
		"player_2_id": req.Player2ID,
		"player_3_id": req.Player3ID,
		"player_4_id": req.Player4ID,
		"parlour_id":  req.ParlourID,
	}

	err = s.Update(&models.Match{}, id, updates)
	if err != nil {
		return nil, err
	}

	match.Player1ID = req.Player1ID
	match.Player2ID = req.Player2ID
	match.Player3ID = req.Player3ID
	match.Player4ID = req.Player4ID
	match.ParlourID = req.ParlourID

	return match, nil
}

func (s *MatchService) DeleteMatch(id uint64) error {
	return s.Delete(&models.Match{}, id)
}

func (s *MatchService) ApproveMatch(id uint64, approvedBy uint64) (*models.Match, error) {
	match, err := s.GetMatchByID(id)

	if err != nil {
		return nil, err
	}

	if match.ApprovedBy != nil || match.ApprovedAt != nil {
		return nil, fmt.Errorf("match is already approved")
	}

	now := time.Now()
	updates := map[string]any{
		"approved_by": approvedBy,
		"approved_at": &now,
	}

	err = s.Update(&models.Match{}, id, updates)
	if err != nil {
		return nil, err
	}

	match.ApprovedAt = &now
	match.ApprovedBy = &approvedBy

	return match, nil
}

func (s *MatchService) GetMatchesByParlour(parlourID uint64, queryPaginate commons.QueryPagination) ([]models.Match, error) {
	var matches []models.Match
	query := s.DB.Where("parlour_id = ?", parlourID)
	preloads := []string{"Player1", "Player2", "Player3", "Player4", "Parlour", "Parlour.Province"}
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	if queryPaginate.Limit > 0 {
		query = query.Limit(queryPaginate.Limit)
	}
	if queryPaginate.Offset > 0 {
		query = query.Offset(queryPaginate.Offset)
	}
	err := query.Find(&matches).Error
	if err != nil {
		return nil, err
	}
	return matches, nil
}
