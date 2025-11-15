package database

import "gorm.io/gorm"

// MockConnectDatabase returns a mock Database without touching MySQL
func MockConnectDatabase() *Database {
	return &Database{
		Conn: &gorm.DB{}, // empty gorm DB prevents nil pointer issues
	}
}
