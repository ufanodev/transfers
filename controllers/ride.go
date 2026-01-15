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

// GET /api/v1/rides - Listar todos los viajes ejecutados o en curso
func (ctrl *RideController) GetAll(c *gin.Context) {
	var items []models.Ride
	// Preload de Booking para saber el origen/destino y Driver para saber quién conduce
	if err := ctrl.DB.Preload("Booking").Preload("Driver.User").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al recuperar viajes"})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GET /api/v1/rides/:id - Detalle de un viaje específico
func (ctrl *RideController) Get(c *gin.Context) {
	var item models.Ride
	if err := ctrl.DB.Preload("Booking").Preload("Driver.User").First(&item, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Viaje no encontrado"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// POST /api/v1/rides - Iniciar un viaje (Cuando el cliente sube al vehículo)
func (ctrl *RideController) POST(c *gin.Context) {
	var item models.Ride
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo registrar el inicio del viaje"})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// PUT /api/v1/rides/:id - Actualizar datos del viaje (ej: cambios en la ruta o estado)
func (ctrl *RideController) PUT(c *gin.Context) {
	var item models.Ride
	if err := ctrl.DB.First(&item, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Viaje no encontrado"})
		return
	}

	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctrl.DB.Save(&item)
	c.JSON(http.StatusOK, item)
}

// DELETE /api/v1/rides/:id - Eliminar registro de viaje
func (ctrl *RideController) DELETE(c *gin.Context) {
	if err := ctrl.DB.Delete(&models.Ride{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el registro"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Registro de viaje eliminado"})
}
