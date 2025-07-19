package services

import (
	"fmt"
	"time"

	"github.com/riichi-mahjong-dev/backend-riichi/commons"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/jobs"
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

func (s *MatchService) CreateMatch(req *models.MatchRequest, userId uint64, role string) (*models.Match, error) {
	match := &models.Match{
		ParlourID: req.ParlourID,
	}

	if role == "player" {
		match.CreatedBy = &userId
	} else {
		now := time.Now()
		match.ApprovedBy = &userId
		match.ApprovedAt = &now
	}

	err := s.Create(match)
	if err != nil {
		return nil, err
	}

	matchPlayers := []models.MatchPlayer{}

	for _, player := range req.Players {
		matchPlayers = append(matchPlayers, models.MatchPlayer{
			MatchID:  match.ID,
			PlayerID: *player.Player,
		})
	}

	err = s.Create(matchPlayers)

	if err != nil {
		return nil, err
	}

	match.Players = matchPlayers

	return match, nil
}

func (s *MatchService) PointMatch(id uint64, req *models.PointMatchRequest, userId uint64) (*models.Match, error) {
	match, err := s.GetMatchByID(id)

	if err != nil {
		return nil, err
	}

	if err := s.checkAdminPermission(userId, match.Parlour.ProvinceID, match.ParlourID); err != nil {
		return nil, fmt.Errorf("you dont't have authority to input point this match")
	}

	err = s.DB.Transaction(func(tx *gorm.DB) error {
		txService := s.WithTx(tx)
		for _, pointMatch := range req.PointMatchPlayers {
			updates := map[string]any{
				"point": pointMatch.Score,
			}
			err := txService.Update(&models.MatchPlayer{}, *pointMatch.MatchPlayerId, updates)
			if err != nil {
				return err
			}
		}

		err := jobs.EnqueueJob(s.DB, "calculate_mmr", map[string]any{
			"id": id,
		})

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.GetMatchByID(id)
}

func (s *MatchService) GetMatchByID(id uint64) (*models.Match, error) {
	var match models.Match
	preloads := []string{"Parlour", "Parlour.Province", "Players.Player", "Creator"}
	err := s.GetWithPreload(&match, id, preloads...)
	if err != nil {
		return nil, err
	}
	return &match, nil
}

func (s *MatchService) GetAllMatches(queryPaginate commons.QueryPagination) ([]models.Match, error) {
	var matches []models.Match
	preloads := []string{"Parlour", "Parlour.Province", "Players.Player", "Creator"}
	err := s.GetAllWithPreload(&matches, queryPaginate, preloads...)
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func (s *MatchService) UpdateMatch(id uint64, req *models.UpdateMatchRequest, userId uint64, role string) (*models.Match, error) {
	match, err := s.GetMatchByID(id)

	if err != nil {
		return nil, err
	}

	if role == "player" && match.CreatedBy != &userId {
		return nil, fmt.Errorf("cannot update match, you are not the creator of this match")
	}

	if role == "player" && (match.ApprovedAt != nil || match.ApprovedBy != nil) {
		return nil, fmt.Errorf("match is already approved, cannot changed anymore")
	}

	if role == "admin" {
		if err := s.checkAdminPermission(userId, match.Parlour.ProvinceID, match.ParlourID); err != nil {
			return nil, fmt.Errorf("you dont't have authority to update this match")
		}
	}

	updates := map[string]any{
		"parlour_id": req.ParlourID,
	}

	err = s.Update(&models.Match{}, id, updates)
	if err != nil {
		return nil, err
	}

	for _, player := range req.Players {
		// if player.MatchPlayerID != nil {
		playerUpdate := map[string]any{
			"player_id": player.Player,
		}

		err = s.DB.Model(&models.MatchPlayer{}).Where("match_id = ?", id).Where("id = ?", player.MatchPlayerID).Updates(playerUpdate).Error

		if err != nil {
			return nil, err
		}
		// continue
		// }

		// err = s.Create(&models.MatchPlayer{
		// 	MatchID:  id,
		// 	PlayerID: *player.Player,
		// })

		// if err != nil {
		// 	return nil, err
		// }
	}

	match, err = s.GetMatchByID(id)

	if err != nil {
		return nil, err
	}

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

	if err := s.checkAdminPermission(approvedBy, match.Parlour.ProvinceID, match.ParlourID); err != nil {
		return nil, fmt.Errorf("you don't have authority to approve this match")
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

func (s *MatchService) GetAllMatchesByAdmin(queryPaginate commons.QueryPagination, userId uint64) ([]models.Match, error) {
	var adminPermissions []models.AdminPermission
	var parlourIds []uint64
	err := s.DB.Where("admin_id = ?", userId).Find(adminPermissions).Error

	for _, adminPermission := range adminPermissions {
		parlourIds = append(parlourIds, adminPermission.ParlourID)
	}

	var matches []models.Match
	preloads := []string{"Parlour", "Parlour.Province", "Players.Player", "Creator"}

	query := s.DB

	for _, preload := preloads {
		query = query.Preload(preload)
	}

	if queryPaginate.Limit > 0 {
		query = query.Limit(queryPaginate.Limit)
	}
	if queryPaginate.Offset > 0 {
		query = query.Offset(queryPaginate.Offset)
	}
	
	err = query.Where("parlour_id IN ?", parlourIds).Find(matches).Error

	if err != nil {
		return nil, err
	}

	return matches, nil
}

func (s *MatchService) GetMatchesByParlour(parlourID uint64, queryPaginate commons.QueryPagination) ([]models.Match, error) {
	var matches []models.Match
	query := s.DB.Where("parlour_id = ?", parlourID)
	preloads := []string{"Parlour", "Parlour.Province"}
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

func (s *MatchService) checkAdminPermission(adminId uint64, provinceId uint64, parlourId uint64) error {
	var adminPermission models.AdminPermission
	return s.DB.
		Where("admin_id = ?", adminId).
		Where("province_id = ?", provinceId).
		Where("parlour_id = ?", parlourId).
		First(&adminPermission).Error
}
