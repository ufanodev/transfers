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

// GET /api/v1/vehicles
func (ctrl *VehicleController) GetAll(c *gin.Context) {
	var items []models.Vehicle
	// Preload de Company para ver a quién pertenece el vehículo
	if err := ctrl.DB.Preload("Company").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al listar flota"})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GET /api/v1/vehicles/:id
func (ctrl *VehicleController) Get(c *gin.Context) {
	var item models.Vehicle
	if err := ctrl.DB.Preload("Company").First(&item, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vehículo no encontrado"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// POST /api/v1/vehicles
func (ctrl *VehicleController) POST(c *gin.Context) {
	var item models.Vehicle
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	// Validación de Matrícula Duplicada (Usando PlateNumber según tu modelo)
	var existing models.Vehicle
	result := ctrl.DB.Where("plate_number = ?", item.PlateNumber).First(&existing)

	if result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "La matrícula " + item.PlateNumber + " ya existe"})
		return
	}

	if err := ctrl.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al registrar vehículo"})
		return
	}

	c.JSON(http.StatusCreated, item)
}

// PUT /api/v1/vehicles/:id
func (ctrl *VehicleController) PUT(c *gin.Context) {
	var item models.Vehicle
	id := c.Param("id")

	if err := ctrl.DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vehículo no hallado"})
		return
	}

	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error en los datos"})
		return
	}

	if err := ctrl.DB.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar"})
		return
	}

	c.JSON(http.StatusOK, item)
}

// DELETE /api/v1/vehicles/:id
func (ctrl *VehicleController) DELETE(c *gin.Context) {
	if err := ctrl.DB.Delete(&models.Vehicle{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo eliminar"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Vehículo eliminado"})
}
