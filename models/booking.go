package models

import (
	"time"

	"gorm.io/gorm"
)

type Booking struct {
	ID        uint  `gorm:"primaryKey;autoIncrement" json:"id"`
	ClientID  uint  `gorm:"not null" json:"client_id"`
	CompanyID *uint `json:"company_id"`

	// Datos del Cliente
	ClientName  string `gorm:"type:varchar(100);not null" json:"client_name"`
	ClientPhone string `gorm:"type:varchar(25);not null" json:"client_phone"`

	// Geolocation
	OriginAddress string  `gorm:"type:varchar(255);not null" json:"origin_address"`
	OriginLat     float64 `gorm:"type:decimal(10,8);not null" json:"origin_lat"`
	OriginLng     float64 `gorm:"type:decimal(11,8);not null" json:"origin_lng"`
	DestAddress   string  `gorm:"type:varchar(255);not null" json:"dest_address"`
	DestLat       float64 `gorm:"type:decimal(10,8);not null" json:"dest_lat"`
	DestLng       float64 `gorm:"type:decimal(11,8);not null" json:"dest_lng"`

	// Carga y Pasajeros
	Pax     int  `gorm:"default:1" json:"pax"`
	Maletas int  `gorm:"default:0" json:"maletas"`
	Animal  bool `gorm:"default:false" json:"animal"`

	// Timing Logic
	ScheduledDate time.Time `gorm:"type:date;not null" json:"scheduled_date"`
	ScheduledTime string    `gorm:"type:varchar(10);not null" json:"scheduled_time"`
	ScheduledAt   time.Time `gorm:"index;not null" json:"scheduled_at"`
	IsImmediate   bool      `gorm:"default:false" json:"is_immediate"`

	RequestedVehicleType string `gorm:"type:varchar(20)" json:"requested_vehicle_type"`
	Status               string `gorm:"type:varchar(20);default:'pending'" json:"status"`

	// NUEVO CAMPO: Observaciones
	Notes string `gorm:"type:varchar(255)" json:"notes"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relaciones
	Client  Client   `gorm:"foreignKey:ClientID" json:"-"`
	Company *Company `gorm:"foreignKey:CompanyID" json:"-"`
}

func (Booking) TableName() string {
	return "bookings"
}
