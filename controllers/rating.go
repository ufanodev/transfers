package controllers

import (
	"net/http"
	"transfers/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RatingController struct {
	DB *gorm.DB
}

// GET /api/v1/ratings - Listar todas las valoraciones
func (ctrl *RatingController) GetAll(c *gin.Context) {
	var items []models.Rating
	// Preload("Booking") para saber qué viaje se está valorando
	if err := ctrl.DB.Preload("Booking").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al recuperar valoraciones"})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GET /api/v1/ratings/:id - Obtener una valoración específica
func (ctrl *RatingController) Get(c *gin.Context) {
	var item models.Rating
	if err := ctrl.DB.Preload("Booking").First(&item, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Valoración no encontrada"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// POST /api/v1/ratings - Crear una nueva valoración
func (ctrl *RatingController) POST(c *gin.Context) {
	var item models.Rating
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validación básica de rango de estrellas (1-5)
	if item.Score < 1 || item.Score > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "La puntuación debe estar entre 1 y 5"})
		return
	}

	if err := ctrl.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo guardar la valoración"})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// PUT /api/v1/ratings/:id - Editar una valoración o comentario
func (ctrl *RatingController) PUT(c *gin.Context) {
	var item models.Rating
	if err := ctrl.DB.First(&item, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Valoración no encontrada"})
		return
	}

	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctrl.DB.Save(&item)
	c.JSON(http.StatusOK, item)
}

// DELETE /api/v1/ratings/:id - Eliminar una valoración
func (ctrl *RatingController) DELETE(c *gin.Context) {
	if err := ctrl.DB.Delete(&models.Rating{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar la valoración"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Valoración eliminada correctamente"})
}
