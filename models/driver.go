package models

import (
	"time"

	"gorm.io/gorm"
)

type Driver struct {
	ID                uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID            uint           `gorm:"not null;uniqueIndex" json:"user_id"`
	CompanyID         uint           `gorm:"not null" json:"company_id"`
	FullName          string         `gorm:"type:varchar(120);not null" json:"full_name"`
	Email             string         `gorm:"type:varchar(100);not null" json:"email"`
	Phone             string         `gorm:"type:varchar(20);not null" json:"phone"`
	LicenseNumber     string         `gorm:"type:varchar(50);not null" json:"license_number"`
	LicenseExpiryDate time.Time      `gorm:"type:date;not null" json:"license_expiry_date"`
	IsAvailable       bool           `gorm:"default:false" json:"is_available"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	// Relaciones
	User    User    `gorm:"foreignKey:UserID" json:"-"`
	Company Company `gorm:"foreignKey:CompanyID" json:"-"`
}

func (Driver) TableName() string {
	return "drivers"
}
