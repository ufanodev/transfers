package controllers

import (
	"fmt"
	"net/http"
	"time"
	"transfers/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PaymentController struct {
	DB *gorm.DB
}

// GET /api/v1/payments - Monitor Financiero
func (ctrl *PaymentController) GetAll(c *gin.Context) {
	var items []models.Payment
	// Preload Booking y Ride para saber exactamente qué se está pagando
	if err := ctrl.DB.Preload("Booking").Order("id desc").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al recuperar pagos: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GET /api/v1/payments/:id
func (ctrl *PaymentController) Get(c *gin.Context) {
	var item models.Payment
	if err := ctrl.DB.Preload("Booking").First(&item, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pago no encontrado"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// POST /api/v1/payments - HITO 4: Registro y Cierre Financiero
func (ctrl *PaymentController) POST(c *gin.Context) {
	var item models.Payment
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de pago inválidos"})
		return
	}

	tx := ctrl.DB.Begin()

	// 1. Establecer fecha de pago si viene marcado como completado
	if item.Status == "completed" && item.PaidAt == nil {
		now := time.Now()
		item.PaidAt = &now
	}

	// 2. Crear Registro
	if err := tx.Create(&item).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo registrar el pago"})
		return
	}

	// 3. ENCADENAMIENTO: Actualizar Booking y Generar Evento
	if item.BookingID != 0 {
		// Actualizar estado del booking
		statusToUpdate := "paid"
		if item.Status != "completed" {
			statusToUpdate = "payment_pending"
		}
		tx.Model(&models.Booking{}).Where("id = ?", item.BookingID).Update("status", statusToUpdate)

		// Registrar en el log de auditoría
		event := models.BookingEvent{
			BookingID:   item.BookingID,
			EventType:   "payment_registered",
			Description: fmt.Sprintf("Pago de %.2f %s registrado (Método: %s). Estado: %s", item.Amount, item.Currency, item.Method, item.Status),
			NewValue:    statusToUpdate,
			CreatedBy:   "FINANCIAL_MODULE",
		}
		tx.Create(&event)
	}

	tx.Commit()
	c.JSON(http.StatusCreated, item)
}

// PUT /api/v1/payments/:id - Actualizar (ej: de 'pending' a 'completed')
func (ctrl *PaymentController) PUT(c *gin.Context) {
	var item models.Payment
	id := c.Param("id")
	if err := ctrl.DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pago no encontrado"})
		return
	}

	oldStatus := item.Status

	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	tx := ctrl.DB.Begin()

	// Si el pago pasa a completado, actualizamos la fecha y lanzamos evento
	if item.Status == "completed" && oldStatus != "completed" {
		now := time.Now()
		item.PaidAt = &now

		event := models.BookingEvent{
			BookingID:   item.BookingID,
			EventType:   "payment_cleared",
			Description: "El pago ha sido verificado y completado.",
			OldValue:    oldStatus,
			NewValue:    "completed",
			CreatedBy:   "ADMIN_CASHIER",
		}
		tx.Create(&event)
		tx.Model(&models.Booking{}).Where("id = ?", item.BookingID).Update("status", "paid")
	}

	if err := tx.Save(&item).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar pago"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, item)
}

// DELETE /api/v1/payments/:id
func (ctrl *PaymentController) DELETE(c *gin.Context) {
	if err := ctrl.DB.Delete(&models.Payment{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Registro de pago eliminado"})
}
