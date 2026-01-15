package controllers

import (
	"net/http"
	"transfers/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BookingController struct{ DB *gorm.DB }

func (ctrl *BookingController) GetAll(c *gin.Context) {
	var items []models.Booking
	ctrl.DB.Preload("Client").Find(&items) // Preload ayuda a traer info relacionada
	c.JSON(http.StatusOK, items)
}

func (ctrl *BookingController) Get(c *gin.Context) {
	var item models.Booking
	if err := ctrl.DB.Preload("Client").First(&item, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No encontrado"})
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
	ctrl.DB.Create(&item)
	c.JSON(http.StatusCreated, item)
}

func (ctrl *BookingController) PUT(c *gin.Context) {
	var item models.Booking
	if err := ctrl.DB.First(&item, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No encontrado"})
		return
	}
	c.ShouldBindJSON(&item)
	ctrl.DB.Save(&item)
	c.JSON(http.StatusOK, item)
}

func (ctrl *BookingController) DELETE(c *gin.Context) {
	ctrl.DB.Delete(&models.Booking{}, c.Param("id"))
	c.JSON(http.StatusOK, gin.H{"message": "Reserva eliminada"})
}
