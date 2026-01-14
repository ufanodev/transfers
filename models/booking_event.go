package models

import (
	"time"
)

type BookingEvent struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	BookingID uint      `gorm:"not null" json:"booking_id"`
	EventType string    `gorm:"type:varchar(50)" json:"event_type"` // status_change, driver_assigned, etc.
	OldStatus string    `gorm:"type:varchar(20)" json:"old_status"`
	NewStatus string    `gorm:"type:varchar(20)" json:"new_status"`
	ActorID   uint      `json:"actor_id"`                 // ID del usuario que disparó el evento
	Payload   string    `gorm:"type:json" json:"payload"` // Datos extra (coordenadas, motivos)
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relación
	Booking Booking `gorm:"foreignKey:BookingID" json:"-"`
}

func (BookingEvent) TableName() string {
	return "booking_events"
}
