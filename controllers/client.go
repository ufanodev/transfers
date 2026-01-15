package controllers

import (
	"net/http"
	"transfers/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ClientController struct {
	DB *gorm.DB
}

// GET /api/v1/clients
func (ctrl *ClientController) GetAll(c *gin.Context) {
	var items []models.Client
	// Preload("User") trae la información de la tabla users asociada
	if err := ctrl.DB.Preload("User").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al recuperar clientes"})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GET /api/v1/clients/:id
func (ctrl *ClientController) Get(c *gin.Context) {
	var item models.Client
	if err := ctrl.DB.Preload("User").First(&item, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cliente no encontrado"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// POST /api/v1/clients
func (ctrl *ClientController) POST(c *gin.Context) {
	var item models.Client
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear el cliente"})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// PUT /api/v1/clients/:id
func (ctrl *ClientController) PUT(c *gin.Context) {
	var item models.Client
	if err := ctrl.DB.First(&item, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cliente no encontrado"})
		return
	}

	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctrl.DB.Save(&item)
	c.JSON(http.StatusOK, item)
}

// DELETE /api/v1/clients/:id
func (ctrl *ClientController) DELETE(c *gin.Context) {
	if err := ctrl.DB.Delete(&models.Client{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el cliente"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Cliente eliminado correctamente"})
}
