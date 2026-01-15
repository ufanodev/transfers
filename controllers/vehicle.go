package controllers

import (
	"net/http"
	"transfers/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type VehicleController struct {
	DB *gorm.DB
}

// GET /api/v1/vehicles - Listar toda la flota
func (ctrl *VehicleController) GetAll(c *gin.Context) {
	var items []models.Vehicle
	// Traemos la información de la empresa dueña del vehículo
	if err := ctrl.DB.Preload("Company").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al recuperar vehículos"})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GET /api/v1/vehicles/:id - Detalle de un vehículo específico
func (ctrl *VehicleController) Get(c *gin.Context) {
	var item models.Vehicle
	if err := ctrl.DB.Preload("Company").First(&item, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vehículo no encontrado"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// POST /api/v1/vehicles - Registrar un nuevo vehículo en la flota
func (ctrl *VehicleController) POST(c *gin.Context) {
	var item models.Vehicle
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo registrar el vehículo"})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// PUT /api/v1/vehicles/:id - Actualizar datos (matrícula, modelo, estado)
func (ctrl *VehicleController) PUT(c *gin.Context) {
	var item models.Vehicle
	if err := ctrl.DB.First(&item, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vehículo no encontrado"})
		return
	}

	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctrl.DB.Save(&item)
	c.JSON(http.StatusOK, item)
}

// DELETE /api/v1/vehicles/:id - Dar de baja un vehículo
func (ctrl *VehicleController) DELETE(c *gin.Context) {
	if err := ctrl.DB.Delete(&models.Vehicle{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el vehículo"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Vehículo eliminado correctamente"})
}
