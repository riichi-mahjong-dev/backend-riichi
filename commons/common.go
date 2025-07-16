package commons

import (
	"github.com/riichi-mahjong-dev/backend-riichi/configs"
	"github.com/riichi-mahjong-dev/backend-riichi/database"
	"github.com/riichi-mahjong-dev/backend-riichi/utils"
)

type AppConfig struct {
	Db     *database.Database
	Mailer *utils.Emailer
	Env    *configs.EnvConfig
}

type QueryPagination struct {
	Search  string `json:"q"`
	SortBy  string `json:"sortBy"`
	Sort    string `json:"sort"`
	Page    int    `json:"page"`
	Limit   int    `json:"limit"`
	Offset  int
	Filters map[string]string `json:"filters"`
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
