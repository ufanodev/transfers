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

// GET /api/v1/rides - Listar todos los servicios/viajes
func (ctrl *RideController) GetAll(c *gin.Context) {
	var items []models.Ride
	// Preload de todo lo necesario para la logística
	if err := ctrl.DB.Preload("Booking").Preload("Driver").Preload("Vehicle").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al recuperar servicios"})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GET /api/v1/rides/:id - Detalle de un viaje específico
func (ctrl *RideController) Get(c *gin.Context) {
	var item models.Ride
	if err := ctrl.DB.Preload("Booking").Preload("Driver").Preload("Vehicle").First(&item, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Servicio no encontrado"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// POST /api/v1/rides - Crear un servicio (Asignación logística)
func (ctrl *RideController) POST(c *gin.Context) {
	var item models.Ride
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de asignación inválidos"})
		return
	}

	// Lógica extra opcional: Podrías marcar el vehículo o conductor como "ocupado" aquí
	if err := ctrl.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear la asignación de viaje"})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// PUT /api/v1/rides/:id - Actualizar estado (Iniciado, Finalizado, Cancelado)
func (ctrl *RideController) PUT(c *gin.Context) {
	var item models.Ride
	if err := ctrl.DB.First(&item, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Servicio no encontrado"})
		return
	}

	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al procesar actualización"})
		return
	}

	ctrl.DB.Save(&item)
	c.JSON(http.StatusOK, item)
}

// DELETE /api/v1/rides/:id - Eliminar registro de servicio
func (ctrl *RideController) DELETE(c *gin.Context) {
	if err := ctrl.DB.Delete(&models.Ride{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar servicio"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Servicio eliminado correctamente"})
}
