package controllers

import (
	"net/http"
	"transfers/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RideController struct {
	DB *gorm.DB
}

// GET /api/v1/rides
// Lista todos los servicios con información cruzada de Driver, Vehicle y Booking
func (ctrl *RideController) GetAll(c *gin.Context) {
	var items []models.Ride
	// Preload carga las relaciones para que el frontend vea nombres y matrículas
	if err := ctrl.DB.Preload("Driver").Preload("Vehicle").Preload("Booking").Order("id desc").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al recuperar monitor de viajes"})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GET /api/v1/rides/:id
func (ctrl *RideController) Get(c *gin.Context) {
	var item models.Ride
	if err := ctrl.DB.Preload("Driver").Preload("Vehicle").Preload("Booking").First(&item, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Viaje no encontrado"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// POST /api/v1/rides
// Crea la asignación logística (Paso de Booking a Ride)
func (ctrl *RideController) POST(c *gin.Context) {
	var item models.Ride
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de asignación inválidos"})
		return
	}

	// Transacción: Crear Ride y marcar Conductor/Vehículo como ocupados (opcional)
	tx := ctrl.DB.Begin()

	if err := tx.Create(&item).Error; err != nil {
		tx.Rollback()
		// Error común: BookingID ya asignado (si mantienes el uniqueIndex)
		c.JSON(http.StatusConflict, gin.H{"error": "Esta reserva ya tiene un viaje asignado"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusCreated, item)
}

// PUT /api/v1/rides/:id
// Actualización PRO: Permite cambiar estados (arriving, started, finished) y capturar coords
func (ctrl *RideController) PUT(c *gin.Context) {
	var item models.Ride
	id := c.Param("id")

	if err := ctrl.DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Registro de viaje no hallado"})
		return
	}

	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de actualización incorrectos"})
		return
	}

	// Lógica automática: Si el estado cambia a 'finished', marcar IsFinished true
	// Nota: Esto depende de si decides implementar el campo 'Status' recomendado
	// if item.Status == "finished" { item.IsFinished = true }

	if err := ctrl.DB.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron guardar los cambios logísticos"})
		return
	}

	c.JSON(http.StatusOK, item)
}

// DELETE /api/v1/rides/:id
// Borrado lógico (Soft Delete)
func (ctrl *RideController) DELETE(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.DB.Delete(&models.Ride{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el registro"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Viaje eliminado/cancelado correctamente"})
}
