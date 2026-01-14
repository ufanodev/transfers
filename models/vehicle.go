package models

import (
	"time"

	"gorm.io/gorm"
)

type Vehicle struct {
	ID                uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	CompanyID         uint           `gorm:"not null" json:"company_id"`
	Make              string         `gorm:"type:varchar(50)" json:"make"`
	ModelName         string         `gorm:"type:varchar(50)" json:"model_name"`
	PlateNumber       string         `gorm:"type:varchar(20);not null;uniqueIndex" json:"plate_number"`
	LicenseNumber     string         `gorm:"type:varchar(50);not null" json:"license_number"`
	LicenseExpiryDate time.Time      `gorm:"type:date;not null" json:"license_expiry_date"`
	VehicleType       string         `gorm:"type:enum('economy','comfort','business','vip','minivan');not null" json:"vehicle_type"`
	Category          string         `gorm:"type:varchar(50)" json:"category"`
	CapacityPax       int            `gorm:"default:4" json:"capacity_pax"`
	IsActive          bool           `gorm:"default:true" json:"is_active"`
	BaseFare          float64        `gorm:"type:decimal(10,2)" json:"base_fare"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	// Relación
	Company Company `gorm:"foreignKey:CompanyID" json:"-"`
}

func (Vehicle) TableName() string {
	return "vehicles"
}
