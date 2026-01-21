package models

import (
	"time"

	"gorm.io/gorm"
)

type Ride struct {
	ID        uint `gorm:"primaryKey;autoIncrement" json:"id"`
	BookingID uint `gorm:"not null;uniqueIndex" json:"booking_id"`
	DriverID  uint `gorm:"not null" json:"driver_id"`
	VehicleID uint `gorm:"not null" json:"vehicle_id"`

	// Tiempos Reales
	StartTimeReal   *time.Time `json:"start_time_real"`
	EndTimeReal     *time.Time `json:"end_time_real"`
	WaitTimeMinutes int        `gorm:"type:bigint" json:"wait_time_minutes"`

	// Kilometraje
	KmStart float64 `gorm:"type:decimal(10,2)" json:"km_start"`
	KmEnd   float64 `gorm:"type:decimal(10,2)" json:"km_end"`
	KmTotal float64 `gorm:"type:decimal(10,2)" json:"km_total"`

	// Tarifas y Extras
	IsHoliday    bool    `gorm:"default:false" json:"is_holiday"`
	IsNightShift bool    `gorm:"default:false" json:"is_night_shift"`
	ExtraCharges float64 `gorm:"type:decimal(10,2);default:0.00" json:"extra_charges"`
	TotalAmount  float64 `gorm:"type:decimal(12,2)" json:"total_amount"`

	// Control y Voucher
	VoucherNumber string `gorm:"type:varchar(50);uniqueIndex" json:"voucher_number"`
	IsFinished    bool   `gorm:"default:false" json:"is_finished"`

	// NUEVOS CAMPOS: Geolocalización y Direcciones (Crucial para el Tracking)
	OriginAddress      string  `gorm:"type:varchar(255)" json:"origin_address"`
	OriginLat          float64 `gorm:"type:decimal(9,6)" json:"origin_lat"`
	OriginLng          float64 `gorm:"type:decimal(9,6)" json:"origin_lng"`
	DestinationAddress string  `gorm:"type:varchar(255)" json:"destination_address"`
	DestinationLat     float64 `gorm:"type:decimal(9,6)" json:"destination_lat"`
	DestinationLng     float64 `gorm:"type:decimal(9,6)" json:"destination_lng"`

	// Auditoría GORM
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relaciones (Preload)
	Booking Booking `gorm:"foreignKey:BookingID" json:"-"`
	Driver  Driver  `gorm:"foreignKey:DriverID" json:"-"`
	Vehicle Vehicle `gorm:"foreignKey:VehicleID" json:"-"`
}

func (Ride) TableName() string {
	return "rides"
}
