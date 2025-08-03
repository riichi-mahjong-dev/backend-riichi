package jobs

import (
	"context"
	"encoding/json"
	"math"
	"sort"
	"time"

	"github.com/riichi-mahjong-dev/backend-riichi/internal/models"
	"gorm.io/gorm"
)

type DataCalculate struct {
	MatchPlayerID uint64
	PlayerID      uint64
	OldRank       int
	NewRank       int
	Decay         int
	Score         int
	GamesPlayed   int
}

func HandleMMR(_ context.Context, db *gorm.DB, job models.Job) error {
	var jobMatch models.JobMatch
	if err := json.Unmarshal(job.Payload, &jobMatch); err != nil {
		return err
	}

	var match models.Match
	if err := db.Preload("MatchPlayers.Player").First(&match, jobMatch.ID).Error; err != nil {
		return err
	}

	dataCalculate := []DataCalculate{}

	for _, matchPlayer := range match.MatchPlayers {
		oldRank := matchPlayer.Player.Rank
		if matchPlayer.Pinalty == nil {
			var lastPlayed time.Time
			if matchPlayer.Player.LastMatchAt != nil {
				lastPlayed = *matchPlayer.Player.LastMatchAt
			} else {
				lastPlayed = time.Now()
			}
			decayRank := int(decayInactiveR(matchPlayer.Player.Rank, lastPlayed))
			dataCalculate = append(dataCalculate, DataCalculate{
				MatchPlayerID: matchPlayer.ID,
				PlayerID:      matchPlayer.PlayerID,
				OldRank:       oldRank,
				NewRank:       decayRank,
				Decay:         oldRank - decayRank,
				Score:         *matchPlayer.Point,
				GamesPlayed:   matchPlayer.Player.GameCount,
			})
		} else {
			dataCalculate = append(dataCalculate, DataCalculate{
				MatchPlayerID: matchPlayer.ID,
				PlayerID:      matchPlayer.PlayerID,
				OldRank:       oldRank,
				NewRank:       oldRank - *matchPlayer.MmrDelta,
				Decay:         *matchPlayer.Pinalty,
				Score:         *matchPlayer.Point,
				GamesPlayed:   *matchPlayer.GamesPlayed,
			})
		}
	}

	avg := averageScore(dataCalculate)
	avgR := averageR(dataCalculate)

	sort.SliceStable(dataCalculate, func(i, j int) bool {
		return dataCalculate[i].Score > dataCalculate[j].Score
	})

	for index, calculate := range dataCalculate {
		placement := index + 1
		change := calculateRChange(calculate, avg, avgR, placement)

		matchPlayerUpdate := map[string]any{
			"mmr_delta":    change,
			"pinalty":      calculate.Decay,
			"games_played": calculate.GamesPlayed,
		}
		if err := db.Model(&models.MatchPlayer{}).Where("id = ?", calculate.MatchPlayerID).Updates(matchPlayerUpdate).Error; err != nil {
			return err
		}

		playerUpdate := map[string]any{
			"rank": calculate.NewRank + int(change),
		}

		if calculate.Decay == 0 {
			playerUpdate["last_match_at"] = time.Now()
		}

		if err := db.Model(&models.Player{}).Where("id = ?", calculate.PlayerID).Updates(playerUpdate).Error; err != nil {
			return err
		}
	}

	return nil
}

func placementScore(place int) float64 {
	switch place {
	case 1:
		return 1.0
	case 2:
		return 0.3
	case 3:
		return -0.3
	case 4:
		return -1.0
	default:
		panic("Invalid placement")
	}
}

func coefficient(gamesPlayed int) float64 {
	if gamesPlayed < 100 {
		return 40.0
	} else if gamesPlayed < 400 {
		return 32.0 - float64(gamesPlayed)/20.0
	}
	return 16.0
}

func averageScore(players []DataCalculate) int {
	var sum int
	for _, p := range players {
		sum += p.Score
	}
	return sum / len(players)
}

func averageR(players []DataCalculate) int {
	var sum int
	for _, p := range players {
		sum += p.NewRank
	}
	return sum / len(players)
}

func calculateRChange(player DataCalculate, avgScore, avgR, place int) float64 {
	P := placementScore(place)
	scoreBonus := float64(player.Score-avgScore) / 1000.0
	ratingAdj := float64(avgR-player.NewRank) / 100.0
	C := coefficient(player.GamesPlayed)

	delta := (P + scoreBonus + ratingAdj) * C
	return math.Round(delta)
}

func decayInactiveR(r int, lastPlayed time.Time) float64 {
	days := time.Since(lastPlayed).Hours() / 24
	if days <= 30 {
		return float64(r)
	}
	decayDays := days - 30
	penalty := math.Min(decayDays/10, 50)
	return float64(r) - penalty
}
