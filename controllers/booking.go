package controllers

import (
	"net/http"
	"transfers/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BookingController struct {
	DB *gorm.DB
}

func (ctrl *BookingController) GetAll(c *gin.Context) {
	var items []models.Booking
	if err := ctrl.DB.Preload("Client").Order("scheduled_at desc").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al listar reservas"})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (ctrl *BookingController) Get(c *gin.Context) {
	var item models.Booking
	if err := ctrl.DB.First(&item, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reserva no encontrada"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (ctrl *BookingController) POST(c *gin.Context) {
	var item models.Booking
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctrl.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear reserva"})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (ctrl *BookingController) PUT(c *gin.Context) {
	var item models.Booking
	id := c.Param("id")
	if err := ctrl.DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reserva no encontrada"})
		return
	}
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctrl.DB.Save(&item)
	c.JSON(http.StatusOK, item)
}

func (ctrl *BookingController) DELETE(c *gin.Context) {
	if err := ctrl.DB.Delete(&models.Booking{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Reserva eliminada"})
}
