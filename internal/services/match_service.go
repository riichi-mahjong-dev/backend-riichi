package services

import (
	"fmt"
	"strings"
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
		BaseService: BaseService{
			DB: db,
		},
	}
}

func (s *MatchService) CreateMatch(req *models.MatchRequest, userId uint64, role string) (*models.Match, error) {
	match := &models.Match{
		ParlourID: req.ParlourID,
		PlayingAt: req.PlayingAt,
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
	preloads := []string{"Parlour", "Parlour.Province", "Creator", "MatchPlayers.Player"}
	err := s.GetWithPreload(&match, id, preloads...)
	if err != nil {
		return nil, err
	}
	return &match, nil
}

func (s *MatchService) GetAllMatches(queryPaginate commons.QueryParams) ([]models.Match, int64, error) {
	preloads := map[string]func(*gorm.DB) *gorm.DB{"Parlour": nil, "Creator": nil, "MatchPlayers.Player": nil}
	return Paginate(
		s.DB,
		models.Match{},
		[]string{},
		queryPaginate.Filters,
		map[string]func(*gorm.DB, any) *gorm.DB{
			"playing_between": func(d *gorm.DB, a any) *gorm.DB {
				val, ok := a.([]string)
				if !ok {
					return d
				}

				if len(val) != 2 {
					return d
				}

				startDate := strings.TrimSpace(val[0])
				endDate := strings.TrimSpace(val[1])
				return d.Where("playing_at BETWEEN ? AND ?", startDate, endDate)
			},
			"created_between": func(d *gorm.DB, a any) *gorm.DB {
				val, ok := a.([]string)
				if !ok {
					return d
				}

				if len(val) != 2 {
					return d
				}

				startDate := strings.TrimSpace(val[0])
				endDate := strings.TrimSpace(val[1])
				return d.Where("created_at BETWEEN ? AND ?", startDate, endDate)
			},
		},
		preloads,
		[]string{},
		queryPaginate.Search,
		queryPaginate.Page,
		queryPaginate.PageSize,
		queryPaginate.OrderBy,
		queryPaginate.Order,
	)
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
		"playing_at": req.PlayingAt,
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

func (s *MatchService) checkAdminPermission(adminId uint64, provinceId uint64, parlourId uint64) error {
	var adminPermission models.AdminPermission
	return s.DB.
		Where("admin_id = ?", adminId).
		Where("province_id = ?", provinceId).
		Where("parlour_id = ?", parlourId).
		First(&adminPermission).Error
}
