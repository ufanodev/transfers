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

// GET /api/v1/drivers - Listar todos los conductores
func (ctrl *DriverController) GetAll(c *gin.Context) {
	var items []models.Driver
	// Preload de User y Company para mostrar nombres en la tabla
	if err := ctrl.DB.Preload("User").Preload("Company").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al recuperar conductores"})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GET /api/v1/drivers/:id - Obtener detalle
func (ctrl *DriverController) Get(c *gin.Context) {
	var item models.Driver
	if err := ctrl.DB.Preload("User").Preload("Company").First(&item, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Conductor no encontrado"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// POST /api/v1/drivers - Crear conductor vinculado a Usuario y Empresa
func (ctrl *DriverController) POST(c *gin.Context) {
	var item models.Driver
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos del formulario inválidos: " + err.Error()})
		return
	}

	// Validación de llaves foráneas manual (opcional pero recomendada)
	if item.UserID == 0 || item.CompanyID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Debe seleccionar un Usuario y una Empresa obligatoriamente"})
		return
	}

	if err := ctrl.DB.Create(&item).Error; err != nil {
		// Error común: el UserID ya está asignado a otro Driver (unique index)
		c.JSON(http.StatusConflict, gin.H{"error": "El usuario seleccionado ya tiene una ficha de conductor activa"})
		return
	}

	c.JSON(http.StatusCreated, item)
}

// PUT /api/v1/drivers/:id - Actualizar conductor
func (ctrl *DriverController) PUT(c *gin.Context) {
	var item models.Driver
	id := c.Param("id")

	if err := ctrl.DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Conductor no encontrado"})
		return
	}

	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// Forzamos el ID de la URL para evitar que cambie el ID primario
	if err := ctrl.DB.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar la ficha"})
		return
	}

	c.JSON(http.StatusOK, item)
}

// DELETE /api/v1/drivers/:id - Borrado lógico
func (ctrl *DriverController) DELETE(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.DB.Delete(&models.Driver{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo eliminar el registro"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Conductor dado de baja"})
}
