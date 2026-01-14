package models

import (
	"time"
)

type Payment struct {
	ID            uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	BookingID     uint       `gorm:"not null" json:"booking_id"`
	Method        string     `gorm:"type:varchar(20)" json:"method"` // card, paypal, cash
	Amount        float64    `gorm:"type:decimal(10,2);not null" json:"amount"`
	Status        string     `gorm:"type:varchar(20);default:'pending'" json:"status"`
	TransactionID string     `gorm:"type:varchar(100)" json:"transaction_id"`
	PaidAt        *time.Time `json:"paid_at"`

	Booking Booking `gorm:"foreignKey:BookingID" json:"-"`
}

func (Payment) TableName() string {
	return "payments"
}
