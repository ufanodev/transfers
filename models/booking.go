package models

import (
	"time"

	"gorm.io/gorm"
)

// Booking representa la entidad principal de una reserva en el sistema.
type Booking struct {
	ID        uint  `gorm:"primaryKey;autoIncrement" json:"id"`
	ClientID  uint  `gorm:"not null;index" json:"client_id"`
	CompanyID *uint `gorm:"index" json:"company_id"` // Opcional, para clientes corporativos

	// --- DATOS DEL PASAJERO ---
	// Se guardan aquí también por si el cliente cambia o se borra,
	// manteniendo la integridad del registro histórico del viaje.
	ClientName  string `gorm:"type:varchar(100);not null" json:"client_name"`
	ClientPhone string `gorm:"type:varchar(25);not null" json:"client_phone"`

	// --- GEOLOCALIZACIÓN ---
	// Uso de decimal para máxima precisión en mapas (Google Maps/Leaflet)
	OriginAddress string  `gorm:"type:varchar(255);not null" json:"origin_address"`
	OriginLat     float64 `gorm:"type:decimal(10,8);not null" json:"origin_lat"`
	OriginLng     float64 `gorm:"type:decimal(11,8);not null" json:"origin_lng"`

	DestAddress string  `gorm:"type:varchar(255);not null" json:"dest_address"`
	DestLat     float64 `gorm:"type:decimal(10,8);not null" json:"dest_lat"`
	DestLng     float64 `gorm:"type:decimal(11,8);not null" json:"dest_lng"`

	// --- CARGA Y PASAJEROS ---
	Pax     int  `gorm:"default:1" json:"pax"`
	Maletas int  `gorm:"default:0" json:"maletas"`
	Animal  bool `gorm:"default:false" json:"animal"`

	// --- LÓGICA TEMPORAL ---
	ScheduledDate time.Time `gorm:"type:date;not null" json:"scheduled_date"`
	ScheduledTime string    `gorm:"type:varchar(10);not null" json:"scheduled_time"`
	ScheduledAt   time.Time `gorm:"index;not null" json:"scheduled_at"` // Timestamp completo para ordenamiento
	IsImmediate   bool      `gorm:"default:false" json:"is_immediate"`  // Para servicios "Asap"

	// --- ESTADO Y PREFERENCIAS ---
	RequestedVehicleType string `gorm:"type:varchar(20)" json:"requested_vehicle_type"`
	Status               string `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, confirmed, dispatched, completed, cancelled

	// --- NOTAS ---
	Notes string `gorm:"type:varchar(255)" json:"notes"` // Observaciones para el conductor

	// --- AUDITORÍA ---
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// --- RELACIONES (GORM) ---
	Client  Client   `gorm:"foreignKey:ClientID" json:"-"`
	Company *Company `gorm:"foreignKey:CompanyID" json:"-"`
}

// TableName especifica el nombre de la tabla en la base de datos.
func (Booking) TableName() string {
	return "bookings"
}
