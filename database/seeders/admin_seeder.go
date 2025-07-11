package seeders

import (
	"log"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AdminSeeder struct {
	DB *gorm.DB
}

func NewAdminSeeder(db *gorm.DB) *AdminSeeder {
	return &AdminSeeder{DB: db}
}

func (s *AdminSeeder) SeedDefaultAdmin() {
	// Check if any admin exists
	var count int64
	s.DB.Model(&models.Admin{}).Count(&count)
	
	if count > 0 {
		log.Println("Admins already exist, skipping default admin creation")
		return
	}

	// Create default super-admin
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return
	}

	admin := &models.Admin{
		Username: "admin",
		Password: string(hashedPassword),
		Role:     models.AdminRoleSuperAdmin,
	}

	if err := s.DB.Create(admin).Error; err != nil {
		log.Printf("Error creating default admin: %v", err)
		return
	}

	log.Println("Default super-admin created successfully")
	log.Println("Username: admin")
	log.Println("Password: admin123")
	log.Println("Role: super-admin")
}
