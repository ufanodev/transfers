package controllers

import (
	"net/http"
	"transfers/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PaymentController struct {
	DB *gorm.DB
}

// GET /api/v1/payments - Listar todos los pagos
func (ctrl *PaymentController) GetAll(c *gin.Context) {
	var items []models.Payment

	// Corregido: Usamos "id desc" en lugar de "created_at" para evitar el error 1054
	// si la tabla no tiene el campo timestamp explícito.
	if err := ctrl.DB.Preload("Booking").Order("id desc").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al recuperar pagos: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GET /api/v1/payments/:id - Detalle de un pago
func (ctrl *PaymentController) Get(c *gin.Context) {
	var item models.Payment
	if err := ctrl.DB.Preload("Booking").First(&item, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pago no encontrado"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// POST /api/v1/payments - Registrar nuevo pago
func (ctrl *PaymentController) POST(c *gin.Context) {
	var item models.Payment
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de pago inválidos"})
		return
	}

	// Transacción para asegurar integridad
	tx := ctrl.DB.Begin()

	if err := tx.Create(&item).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo registrar el pago"})
		return
	}

	// Actualizar estado de la reserva vinculada si existe
	if item.BookingID != 0 {
		tx.Model(&models.Booking{}).Where("id = ?", item.BookingID).Update("status", "paid")
	}

	tx.Commit()
	c.JSON(http.StatusCreated, item)
}

// PUT /api/v1/payments/:id - Actualizar pago
func (ctrl *PaymentController) PUT(c *gin.Context) {
	var item models.Payment
	if err := ctrl.DB.First(&item, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pago no encontrado"})
		return
	}

	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	ctrl.DB.Save(&item)
	c.JSON(http.StatusOK, item)
}

// DELETE /api/v1/payments/:id - Eliminar registro
func (ctrl *PaymentController) DELETE(c *gin.Context) {
	if err := ctrl.DB.Delete(&models.Payment{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Pago eliminado"})
}
