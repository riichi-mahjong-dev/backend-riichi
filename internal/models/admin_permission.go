package models

type AdminPermission struct {
	AdminID    uint64   `gorm:"not null" json:"admin_id"`
	ProvinceID uint64   `gorm:"not null" json:"province_id"`
	ParlourID  uint64   `gorm:"not null" json:"parlour_id"`
	Admin      Admin    `gorm:"foreignKey:admin_id" json:"admin"`
	Province   Province `gorm:"foreignKey:province_id" json:"province"`
	Parlour    Parlour  `gorm:"foreignKey:parlour_id" json:"parlour"`
}
