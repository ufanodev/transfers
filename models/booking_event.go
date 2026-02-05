package models

import (
	"time"
)

// BookingEvent representa un hito histórico en el ciclo de vida de una reserva.
type BookingEvent struct {
	ID        uint `gorm:"primaryKey;autoIncrement" json:"id"`
	BookingID uint `gorm:"not null;index" json:"booking_id"`

	// Tipo de acción: 'created', 'status_change', 'driver_assigned', 'manual_note', etc.
	EventType string `gorm:"type:varchar(50);not null" json:"event_type"`

	// Explicación detallada del suceso.
	Description string `gorm:"type:text" json:"description"`

	// Trazabilidad de cambios (ej: de 'pending' a 'confirmed').
	OldValue string `gorm:"type:varchar(100)" json:"old_value"`
	NewValue string `gorm:"type:varchar(100)" json:"new_value"`

	// Identificador del autor de la acción (Admin o "SYSTEM").
	CreatedBy string `gorm:"type:varchar(100)" json:"created_by"`

	// Fecha y hora exacta del suceso.
	CreatedAt time.Time `json:"created_at"`

	// Relación inversa con el Booking.
	// Usamos json:"-" para evitar ciclos infinitos al serializar desde Booking.
	Booking Booking `gorm:"foreignKey:BookingID" json:"-"`
}

func (BookingEvent) TableName() string {
	return "booking_events"
}
