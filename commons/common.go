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

func toSnakeCasePlural(name string) string {
	// Convert CamelCase to snake_case
	re := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := re.ReplaceAllString(name, "${1}_${2}")
	snake = strings.ToLower(snake)

	// Naive pluralization: just add 's'
	return snake + "s"
}
