package models

type Parlour struct {
	ID         uint64 `gorm:"primaryKey" json:"id"`
	Name       string `gorm:"size:255;not null" json:"name" validate:"required,min:2"`
	Country    string `gorm:"size:255;not null" json:"country" validate:"required"`
	ProvinceID uint64 `gorm:"not null" json:"province_id" validate:"required"`
	Address    string `gorm:"type:text" json:"address"`
	
	// Relations
	Province Province `gorm:"foreignKey:ProvinceID" json:"province,omitempty"`
}

type ParlourRequest struct {
	Name       string `json:"name" validate:"required,min:2"`
	Country    string `json:"country" validate:"required"`
	ProvinceID uint64 `json:"province_id" validate:"required"`
	Address    string `json:"address"`
}

type ParlourResponse struct {
	ID         uint64            `json:"id"`
	Name       string            `json:"name"`
	Country    string            `json:"country"`
	ProvinceID uint64            `json:"province_id"`
	Address    string            `json:"address"`
	Province   *ProvinceResponse `json:"province,omitempty"`
}
