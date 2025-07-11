package main

import (
	"log"

	"github.com/riichi-mahjong-dev/backend-riichi/commons"
	"github.com/riichi-mahjong-dev/backend-riichi/configs"
	"github.com/riichi-mahjong-dev/backend-riichi/database"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/app"
	"github.com/riichi-mahjong-dev/backend-riichi/utils"
)

func main() {
	env := configs.LoadEnv()
	dbConfig := env.LoadDatabaseConfig()
	emailConfig := env.LoadEmailConfig()
	db, err := database.ConnectDatabase(dbConfig)

	if err != nil {
		log.Fatalf("Failed to connect to database %v", err)
		return
	}

	mailer, err := utils.InitializeEmailer(emailConfig)

	if err != nil {
		log.Fatalf("Failed to initialize email %v", err)
		return
	}

	appConfig := commons.AppConfig{
		Db:     db,
		Mailer: mailer,
		Env:    env,
	}

	app.CreateApp(appConfig)
}
