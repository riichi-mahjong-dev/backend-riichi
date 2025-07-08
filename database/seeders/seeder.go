package seeders

import "gorm.io/gorm"

type Seeder struct {
	db *gorm.DB
}

func SeedDB(db *gorm.DB) *Seeder {
	return &Seeder{
		db: db,
	}
}

func (seeder *Seeder) RunSeeder() {

}
