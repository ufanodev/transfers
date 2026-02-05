package controllers

import (
	"fmt"
	"net/http"
	"transfers/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RideController struct {
	DB *gorm.DB
}

// GET /api/v1/rides - Monitor de flota en tiempo real
func (ctrl *RideController) GetAll(c *gin.Context) {
	var items []models.Ride
	if err := ctrl.DB.Preload("Driver").Preload("Vehicle").Preload("Booking").Order("id desc").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al recuperar monitor de viajes"})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GET /api/v1/rides/:id
func (ctrl *RideController) Get(c *gin.Context) {
	var item models.Ride
	id := c.Param("id")
	if err := ctrl.DB.Preload("Driver").Preload("Vehicle").Preload("Booking").First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Viaje no encontrado"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// POST /api/v1/rides - PASO 3: Despacho y Asignación
func (ctrl *RideController) POST(c *gin.Context) {
	var item models.Ride
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de asignación inválidos"})
		return
	}

	tx := ctrl.DB.Begin()

	// 1. HERENCIA INTELIGENTE: Si faltan datos, succionarlos del Booking
	var booking models.Booking
	if err := tx.First(&booking, item.BookingID).Error; err == nil {
		if item.ClientName == "" {
			item.ClientName = booking.ClientName
		}
		if item.ClientPhone == "" {
			item.ClientPhone = booking.ClientPhone
		}
		if item.PickupAddress == "" {
			item.PickupAddress = booking.OriginAddress
		}
		if item.DropoffAddress == "" {
			item.DropoffAddress = booking.DestAddress
		}
		if item.Pax == 0 {
			item.Pax = booking.Pax
		}
		if item.LuggageCount == 0 {
			item.LuggageCount = booking.Maletas
		}
		item.Animals = booking.Animal

		// Actualizar el estado de la reserva original
		tx.Model(&booking).Update("status", "dispatched")
	}

	// 2. CREAR EL RIDE (La ejecución)
	if err := tx.Create(&item).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusConflict, gin.H{"error": "Esta reserva ya tiene un servicio asignado"})
		return
	}

	// 3. PASO 2: LOG AUTOMÁTICO
	event := models.BookingEvent{
		BookingID:   item.BookingID,
		EventType:   "driver_assigned",
		Description: fmt.Sprintf("Conductor asignado. Ride ID: %d", item.ID),
		NewValue:    "dispatched",
		CreatedBy:   "SISTEMA_DESPACHO",
	}
	tx.Create(&event)

	tx.Commit()
	c.JSON(http.StatusCreated, item)
}

// PUT /api/v1/rides/:id - PASO 4: Cierre y Disparo de Pago
func (ctrl *RideController) PUT(c *gin.Context) {
	var item models.Ride
	id := c.Param("id")

	if err := ctrl.DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Registro de viaje no hallado"})
		return
	}

	oldFinished := item.IsFinished

	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos incorrectos"})
		return
	}

	// Cálculo de Kilometraje en caliente
	if item.KmEnd > 0 && item.KmStart > 0 {
		item.KmTotal = item.KmEnd - item.KmStart
	}

	tx := ctrl.DB.Begin()

	// SI EL SERVICIO SE CIERRA AHORA (HITO 4)
	if item.IsFinished && !oldFinished {
		// A. Log de Auditoría
		event := models.BookingEvent{
			BookingID:   item.BookingID,
			EventType:   "ride_completed",
			Description: fmt.Sprintf("Servicio finalizado con %.2f KM totales.", item.KmTotal),
			NewValue:    "completed",
			CreatedBy:   "SISTEMA_CIERRE",
		}
		tx.Create(&event)

		// B. Sincronizar estado del Booking
		tx.Model(&models.Booking{}).Where("id = ?", item.BookingID).Update("status", "completed")

		// C. DISPARAR PAGO (Generar registro de Payment)
		payment := models.Payment{
			BookingID: item.BookingID,
			RideID:    &item.ID,
			Amount:    item.TotalAmount,
			Status:    "pending", // El pago queda pendiente de cobro real
			Method:    "cash",    // Por defecto, se puede cambiar luego
		}
		if err := tx.Create(&payment).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al generar registro de pago"})
			return
		}
	}

	if err := tx.Save(&item).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fallo al sincronizar viaje"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, item)
}

// DELETE /api/v1/rides/:id
func (ctrl *RideController) DELETE(c *gin.Context) {
	id := c.Param("id")
	var item models.Ride

	if err := ctrl.DB.First(&item, id).Error; err == nil {
		// Log de cancelación
		event := models.BookingEvent{
			BookingID:   item.BookingID,
			EventType:   "ride_cancelled",
			Description: "Asignación de viaje cancelada por administración.",
			CreatedBy:   "ADMIN_OPERACIONES",
		}
		ctrl.DB.Create(&event)
		// Restaurar estado del booking
		ctrl.DB.Model(&models.Booking{}).Where("id = ?", item.BookingID).Update("status", "confirmed")
	}

	ctrl.DB.Delete(&models.Ride{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Ride eliminado"})
}
