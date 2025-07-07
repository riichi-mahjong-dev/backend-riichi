package database

import (
	"fmt"
	"time"

	"github.com/riichi-mahjong-dev/backend-riichi/configs"
	// "github.com/riichi-mahjong-dev/backend-riichi/database/seeders"
	// "github.com/riichi-mahjong-dev/backend-riichi/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
	Conn *gorm.DB
}

func ConnectDatabase(dbConfig *configs.DatabaseConfig) (*Database, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBHost, dbConfig.DBPort, dbConfig.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("unable to load env file: %w", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return &Database{Conn: db}, nil
}

func (database *Database) Migrate() {
	database.Conn.AutoMigrate(&models.User{}, &models.Product{})
	fmt.Println("Database migrated")
}

func (database *Database) Seeder(fresh bool) {
	seeders.SeedDB(database.Conn).RunSeeder(fresh)
	fmt.Println("Seeder done")
}
