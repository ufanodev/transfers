package models

import (
	"time"

	"gorm.io/gorm"
)

type Ride struct {
	ID        uint `gorm:"primaryKey;autoIncrement" json:"id"`
	BookingID uint `gorm:"not null;index" json:"booking_id"`
	DriverID  uint `gorm:"not null" json:"driver_id"`
	VehicleID uint `gorm:"not null" json:"vehicle_id"`

	// Datos del Cliente y Carga (NUEVOS)
	ClientName   string `gorm:"type:varchar(100)" json:"client_name"`
	ClientPhone  string `gorm:"type:varchar(20)" json:"client_phone"`
	Pax          int    `gorm:"default:1" json:"pax"`
	Animals      bool   `gorm:"default:false" json:"animals"`
	LuggageCount int    `gorm:"default:0" json:"luggage_count"`

	// Ruta y Geolocalización (NUEVOS)
	PickupAddress  string  `gorm:"type:varchar(255)" json:"pickup_address"`
	DropoffAddress string  `gorm:"type:varchar(255)" json:"dropoff_address"`
	OriginLat      float64 `gorm:"type:decimal(9,6)" json:"origin_lat"`
	OriginLng      float64 `gorm:"type:decimal(9,6)" json:"origin_lng"`
	DestinationLat float64 `gorm:"type:decimal(9,6)" json:"destination_lat"`
	DestinationLng float64 `gorm:"type:decimal(9,6)" json:"destination_lng"`

	// Métricas Operativas
	StartTimeReal   *time.Time `json:"start_time_real"`
	EndTimeReal     *time.Time `json:"end_time_real"`
	WaitTimeMinutes int        `json:"wait_time_minutes"`
	KmStart         float64    `gorm:"type:decimal(10,2)" json:"km_start"`
	KmEnd           float64    `gorm:"type:decimal(10,2)" json:"km_end"`
	KmTotal         float64    `gorm:"type:decimal(10,2)" json:"km_total"`

	// Estado y Finanzas
	Status        string  `gorm:"type:varchar(20);default:'scheduled'" json:"status"` // scheduled, in_progress, completed, cancelled
	IsHoliday     bool    `json:"is_holiday"`
	IsNightShift  bool    `json:"is_night_shift"`
	ExtraCharges  float64 `gorm:"type:decimal(10,2)" json:"extra_charges"`
	TotalAmount   float64 `gorm:"type:decimal(12,2)" json:"total_amount"`
	VoucherNumber string  `gorm:"type:varchar(50)" json:"voucher_number"`
	IsFinished    bool    `json:"is_finished"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relaciones
	Booking Booking `gorm:"foreignKey:BookingID" json:"-"`
	Driver  Driver  `gorm:"foreignKey:DriverID" json:"-"`
	Vehicle Vehicle `gorm:"foreignKey:VehicleID" json:"-"`
}
