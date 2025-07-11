package main

import (
	"log"
	"github.com/riichi-mahjong-dev/backend-riichi/configs"
	"github.com/riichi-mahjong-dev/backend-riichi/database"
	"github.com/riichi-mahjong-dev/backend-riichi/database/seeders"
)

func main() {
	env := configs.LoadEnv()
	dbConfig := env.LoadDatabaseConfig()
	db, err := database.ConnectDatabase(dbConfig)

	if err != nil {
		log.Fatalf("Failed to connect to database %v", err)
		return
	}

	adminSeeder := seeders.NewAdminSeeder(db.Conn)
	adminSeeder.SeedDefaultAdmin()
	
	log.Println("Admin seeding completed")
}
