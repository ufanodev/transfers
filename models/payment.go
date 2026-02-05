package models

import (
	"time"

	"gorm.io/gorm"
)

// Payment representa el registro financiero de un servicio.
type Payment struct {
	ID        uint  `gorm:"primaryKey;autoIncrement" json:"id"`
	BookingID uint  `gorm:"not null;index" json:"booking_id"`
	RideID    *uint `gorm:"index" json:"ride_id"` // Opcional: vinculación directa a la carrera ejecutada

	// Datos del pago
	Method        string     `gorm:"type:varchar(20)" json:"method"` // card, paypal, cash, transfer
	Amount        float64    `gorm:"type:decimal(10,2);not null" json:"amount"`
	Currency      string     `gorm:"type:varchar(3);default:'EUR'" json:"currency"`
	Status        string     `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, completed, failed, refunded
	TransactionID string     `gorm:"type:varchar(100)" json:"transaction_id"`
	PaidAt        *time.Time `json:"paid_at"`

	// Auditoría
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relaciones
	Booking Booking `gorm:"foreignKey:BookingID" json:"-"`
	Ride    *Ride   `gorm:"foreignKey:RideID" json:"-"`
}

func (Payment) TableName() string {
	return "payments"
}
