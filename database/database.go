package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/pressly/goose/v3"
	"github.com/riichi-mahjong-dev/backend-riichi/configs"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/riichi-mahjong-dev/backend-riichi/database/migrations"
	"github.com/riichi-mahjong-dev/backend-riichi/database/seeders"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
	Conn          *gorm.DB
	migrationConn *sql.DB
}

func ConnectDatabase(dbConfig *configs.DatabaseConfig) (*Database, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBHost, dbConfig.DBPort, dbConfig.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("unable to load env file: %w", err)
	}

	sqlConn, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, fmt.Errorf("unable to connect sql", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return &Database{Conn: db, migrationConn: sqlConn}, nil
}

func (database *Database) Migrate() {
	goose.SetDialect("mysql")
	if err := goose.Up(database.migrationConn, "database/migrations"); err != nil {
		log.Error(err)
		return
	}
	fmt.Println("Database migrated")
}

func (database *Database) Seeder(fresh bool) {
	seeders.SeedDB(database.Conn).RunSeeder()
	fmt.Println("Seeder done")
}
