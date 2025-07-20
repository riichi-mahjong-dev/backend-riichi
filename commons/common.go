package commons

import (
	"regexp"
	"strings"

	"github.com/riichi-mahjong-dev/backend-riichi/configs"
	"github.com/riichi-mahjong-dev/backend-riichi/database"
	"github.com/riichi-mahjong-dev/backend-riichi/utils"
)

type AppConfig struct {
	Db     *database.Database
	Mailer *utils.Emailer
	Env    *configs.EnvConfig
}

type QueryParams struct {
	Page     int
	PageSize int
	Search   string
	OrderBy  string
	Order    string
	Filters  map[string]any // Custom filters, e.g., age, status
}

// type PaginationParams struct {
// 	Take   int    `query:"take"`
// 	Skip   int    `query:"skip"`
// 	Search string `query:"search"`
// 	Sort   string `query:"sort"`
// 	SortBy string `query:"sortBy"`
// }

// func (paginationParams *PaginationParams) SetParams(take int, sort, sortBy string) {
// 	if paginationParams.Take == 0 {
// 		paginationParams.Take = 10
// 	}

// 	if paginationParams.Sort == "" {
// 		paginationParams.Sort = sort
// 	}

// 	if paginationParams.SortBy == "" {
// 		paginationParams.SortBy = sortBy
// 	}
// }

func toSnakeCasePlural(name string) string {
	// Convert CamelCase to snake_case
	re := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := re.ReplaceAllString(name, "${1}_${2}")
	snake = strings.ToLower(snake)

	// Naive pluralization: just add 's'
	return snake + "s"
}
