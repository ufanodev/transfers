package models

import (
	"time"
)

type BookingEvent struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	BookingID   uint      `gorm:"not null;index" json:"booking_id"`
	EventType   string    `gorm:"type:varchar(50);not null" json:"event_type"` // created, status_change, driver_assigned, manual_note
	Description string    `gorm:"type:text" json:"description"`
	OldValue    string    `gorm:"type:varchar(100)" json:"old_value"`
	NewValue    string    `gorm:"type:varchar(100)" json:"new_value"`
	CreatedBy   string    `gorm:"type:varchar(100)" json:"created_by"` // Nombre del admin o sistema
	CreatedAt   time.Time `json:"created_at"`

	// Relación
	Booking Booking `gorm:"foreignKey:BookingID" json:"-"`
}

func (BookingEvent) TableName() string {
	return "booking_events"
}
