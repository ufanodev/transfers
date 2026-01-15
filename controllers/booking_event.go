package controllers

import (
	"net/http"
	"transfers/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BookingEventController struct{ DB *gorm.DB }

func (ctrl *BookingEventController) GetAll(c *gin.Context) {
	var items []models.BookingEvent
	// En eventos solemos precargar el Booking relacionado
	ctrl.DB.Preload("Booking").Find(&items)
	c.JSON(http.StatusOK, items)
}

func (ctrl *BookingEventController) Get(c *gin.Context) {
	var item models.BookingEvent
	if err := ctrl.DB.Preload("Booking").First(&item, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Evento no encontrado"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (ctrl *BookingEventController) POST(c *gin.Context) {
	var item models.BookingEvent
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctrl.DB.Create(&item)
	c.JSON(http.StatusCreated, item)
}

func (ctrl *BookingEventController) PUT(c *gin.Context) {
	var item models.BookingEvent
	if err := ctrl.DB.First(&item, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Evento no encontrado"})
		return
	}
	c.ShouldBindJSON(&item)
	ctrl.DB.Save(&item)
	c.JSON(http.StatusOK, item)
}

func (ctrl *BookingEventController) DELETE(c *gin.Context) {
	// Importante: Usar el modelo correcto para el borrado
	if err := ctrl.DB.Delete(&models.BookingEvent{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Evento eliminado"})
}
