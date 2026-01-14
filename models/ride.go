package models

import (
	"time"
)

type Ride struct {
	ID        uint `gorm:"primaryKey;autoIncrement" json:"id"`
	BookingID uint `gorm:"not null;uniqueIndex" json:"booking_id"`
	DriverID  uint `gorm:"not null" json:"driver_id"`
	VehicleID uint `gorm:"not null" json:"vehicle_id"`

	StartTimeReal   *time.Time `json:"start_time_real"`
	EndTimeReal     *time.Time `json:"end_time_real"`
	WaitTimeMinutes int        `json:"wait_time_minutes"`

	KmStart float64 `gorm:"type:decimal(10,2)" json:"km_start"`
	KmEnd   float64 `gorm:"type:decimal(10,2)" json:"km_end"`
	KmTotal float64 `gorm:"type:decimal(10,2)" json:"km_total"`

	IsHoliday    bool `gorm:"default:false" json:"is_holiday"`
	IsNightShift bool `gorm:"default:false" json:"is_night_shift"`

	ExtraCharges float64 `gorm:"type:decimal(10,2);default:0" json:"extra_charges"`
	TotalAmount  float64 `gorm:"type:decimal(12,2)" json:"total_amount"`

	VoucherNumber string `gorm:"type:varchar(50);uniqueIndex" json:"voucher_number"`
	IsFinished    bool   `gorm:"default:false" json:"is_finished"`

	// Relaciones
	Booking Booking `gorm:"foreignKey:BookingID" json:"-"`
	Driver  Driver  `gorm:"foreignKey:DriverID" json:"-"`
	Vehicle Vehicle `gorm:"foreignKey:VehicleID" json:"-"`
}

func (Ride) TableName() string {
	return "rides"
}
