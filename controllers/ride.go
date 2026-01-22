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

// GET /api/v1/rides
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

// POST /api/v1/rides
// Crea el despacho y registra el evento automáticamente
func (ctrl *RideController) POST(c *gin.Context) {
	var item models.Ride
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de asignación inválidos"})
		return
	}

	tx := ctrl.DB.Begin()

	// 1. HERENCIA: Si los datos del cliente vienen vacíos, traerlos del Booking
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

		// Actualizar estado del Booking a 'dispatched'
		tx.Model(&booking).Update("status", "dispatched")
	}

	// 2. Crear el Ride
	if err := tx.Create(&item).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusConflict, gin.H{"error": "La reserva ya tiene un viaje asignado"})
		return
	}

	// 3. AUDITORÍA: Registrar creación de Ride en BookingEvents
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

// PUT /api/v1/rides/:id
// Actualiza métricas y registra eventos de cambio de estado (ej: Finalizado)
func (ctrl *RideController) PUT(c *gin.Context) {
	var item models.Ride
	id := c.Param("id")

	if err := ctrl.DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Registro de viaje no hallado"})
		return
	}

	// Guardamos el estado anterior para comparar
	oldStatus := item.IsFinished

	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de actualización incorrectos"})
		return
	}

	// Lógica de Kilometraje
	if item.KmEnd > 0 && item.KmStart > 0 {
		item.KmTotal = item.KmEnd - item.KmStart
	}

	tx := ctrl.DB.Begin()

	// AUDITORÍA: Si se marca como finalizado, registrar en eventos
	if item.IsFinished && !oldStatus {
		event := models.BookingEvent{
			BookingID:   item.BookingID,
			EventType:   "ride_completed",
			Description: fmt.Sprintf("Carrera finalizada. Km Totales: %.2f", item.KmTotal),
			NewValue:    "completed",
			CreatedBy:   "DRIVER_APP",
		}
		tx.Create(&event)

		// También actualizamos el Booking original a 'completed'
		tx.Model(&models.Booking{}).Where("id = ?", item.BookingID).Update("status", "completed")
	}

	if err := tx.Save(&item).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar cambios"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, item)
}

// DELETE /api/v1/rides/:id
func (ctrl *RideController) DELETE(c *gin.Context) {
	id := c.Param("id")

	// Antes de borrar, recuperamos el Ride para saber el BookingID y registrar el evento
	var item models.Ride
	if err := ctrl.DB.First(&item, id).Error; err == nil {
		event := models.BookingEvent{
			BookingID:   item.BookingID,
			EventType:   "ride_cancelled",
			Description: "Asignación de conductor eliminada. El booking vuelve a estar pendiente.",
			CreatedBy:   "ADMIN_DESPACHO",
		}
		ctrl.DB.Create(&event)

		// Devolvemos el Booking a estado 'confirmed' o 'pending'
		ctrl.DB.Model(&models.Booking{}).Where("id = ?", item.BookingID).Update("status", "confirmed")
	}

	if err := ctrl.DB.Delete(&models.Ride{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el registro"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Viaje eliminado correctamente"})
}
