package controllers

import (
	"net/http"
	"transfers/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DriverController struct {
	DB *gorm.DB
}

// GET /api/v1/drivers - Listar todos los conductores con su info de usuario y empresa
func (ctrl *DriverController) GetAll(c *gin.Context) {
	var items []models.Driver
	// Preload traemos los datos del User y de la Company asociada
	if err := ctrl.DB.Preload("User").Preload("Company").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al recuperar conductores"})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GET /api/v1/drivers/:id - Obtener un conductor específico
func (ctrl *DriverController) Get(c *gin.Context) {
	var item models.Driver
	if err := ctrl.DB.Preload("User").Preload("Company").First(&item, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Conductor no encontrado"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// POST /api/v1/drivers - Registrar nuevo conductor
func (ctrl *DriverController) POST(c *gin.Context) {
	var item models.Driver
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear el conductor"})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// PUT /api/v1/drivers/:id - Actualizar datos (licencia, disponibilidad, etc.)
func (ctrl *DriverController) PUT(c *gin.Context) {
	var item models.Driver
	if err := ctrl.DB.First(&item, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Conductor no encontrado"})
		return
	}

	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctrl.DB.Save(&item)
	c.JSON(http.StatusOK, item)
}

// DELETE /api/v1/drivers/:id - Eliminar conductor
func (ctrl *DriverController) DELETE(c *gin.Context) {
	if err := ctrl.DB.Delete(&models.Driver{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el conductor"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Conductor eliminado correctamente"})
}
