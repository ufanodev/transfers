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

// GET /api/v1/payments - Listar todos los pagos registrados
func (ctrl *PaymentController) GetAll(c *gin.Context) {
	var items []models.Payment
	// Preload("Booking") para saber qué reserva se está pagando
	if err := ctrl.DB.Preload("Booking").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al recuperar pagos"})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GET /api/v1/payments/:id - Obtener detalle de un pago
func (ctrl *PaymentController) Get(c *gin.Context) {
	var item models.Payment
	if err := ctrl.DB.Preload("Booking").First(&item, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pago no encontrado"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// POST /api/v1/payments - Registrar un nuevo pago (Cierre de transacción)
func (ctrl *PaymentController) POST(c *gin.Context) {
	var item models.Payment
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo registrar el pago"})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// PUT /api/v1/payments/:id - Actualizar estado del pago (ej: de 'pending' a 'paid')
func (ctrl *PaymentController) PUT(c *gin.Context) {
	var item models.Payment
	if err := ctrl.DB.First(&item, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pago no encontrado"})
		return
	}

	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctrl.DB.Save(&item)
	c.JSON(http.StatusOK, item)
}

// DELETE /api/v1/payments/:id - Eliminar registro de pago
func (ctrl *PaymentController) DELETE(c *gin.Context) {
	if err := ctrl.DB.Delete(&models.Payment{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el registro de pago"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Registro de pago eliminado correctamente"})
}
