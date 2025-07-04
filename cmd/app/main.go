package main

import (
	"log"

	"github.com/riichi-mahjong-dev/backend/commons"
	"github.com/riichi-mahjong-dev/backend/configs"
	"github.com/riichi-mahjong-dev/backend/database"
	"github.com/riichi-mahjong-dev/backend/internal/app"
	"github.com/riichi-mahjong-dev/backend/utils"
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
		Db:     db.Conn,
		Mailer: mailer,
		Env:    env,
	}

	app.CreateApp(appConfig)
}
