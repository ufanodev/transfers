package controllers

import (
	"fmt"
	"net/http"
	"time"
	"transfers/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BookingController struct {
	DB *gorm.DB
}

func (ctrl *BookingController) GetAll(c *gin.Context) {
	var items []models.Booking
	if err := ctrl.DB.Preload("Client").Order("scheduled_at asc").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al listar reservas"})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (ctrl *BookingController) Get(c *gin.Context) {
	var item models.Booking
	if err := ctrl.DB.Preload("Client").First(&item, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reserva no encontrada"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (ctrl *BookingController) POST(c *gin.Context) {
	var item models.Booking
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. LÓGICA DE TIEMPO: Construir scheduled_at
	// Combina la fecha (ScheduledDate) y el string de hora (ScheduledTime)
	timeStr := fmt.Sprintf("%s %s", item.ScheduledDate.Format("2006-01-02"), item.ScheduledTime)
	if t, err := time.Parse("2006-01-02 15:04", timeStr); err == nil {
		item.ScheduledAt = t
	}

	tx := ctrl.DB.Begin()

	// 2. CREACIÓN DE LA RESERVA
	if err := tx.Create(&item).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear reserva"})
		return
	}

	// 3. DISPARAR EVENTO (Encadenamiento automático)
	event := models.BookingEvent{
		BookingID:   item.ID,
		EventType:   "created",
		Description: fmt.Sprintf("Reserva creada para %s - Ruta: %s a %s", item.ClientName, item.OriginAddress, item.DestAddress),
		NewValue:    "pending",
		CreatedBy:   "Admin System", // Aquí podrías obtener el usuario del JWT
	}

	if err := tx.Create(&event).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al registrar evento de auditoría"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusCreated, item)
}

func (ctrl *BookingController) PUT(c *gin.Context) {
	var item models.Booking
	id := c.Param("id")

	if err := ctrl.DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reserva no encontrada"})
		return
	}

	// Guardamos el estado anterior para la auditoría
	oldStatus := item.Status

	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := ctrl.DB.Begin()

	// Si el estado ha cambiado, disparamos un evento de cambio de estado
	if oldStatus != item.Status {
		event := models.BookingEvent{
			BookingID:   item.ID,
			EventType:   "status_change",
			Description: fmt.Sprintf("Estado actualizado manualmente de %s a %s", oldStatus, item.Status),
			OldValue:    oldStatus,
			NewValue:    item.Status,
			CreatedBy:   "Admin Editor",
		}
		tx.Create(&event)
	}

	if err := tx.Save(&item).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, item)
}

func (ctrl *BookingController) DELETE(c *gin.Context) {
	id := c.Param("id")

	// Antes de borrar, registramos la cancelación en eventos
	var item models.Booking
	if err := ctrl.DB.First(&item, id).Error; err == nil {
		event := models.BookingEvent{
			BookingID:   item.ID,
			EventType:   "deleted",
			Description: "La reserva ha sido eliminada del sistema (Soft Delete)",
			CreatedBy:   "Admin Trash",
		}
		ctrl.DB.Create(&event)
	}

	if err := ctrl.DB.Delete(&models.Booking{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Reserva eliminada"})
}
