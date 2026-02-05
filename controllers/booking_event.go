package controllers

import (
	"net/http"
	"transfers/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BookingEventController struct {
	DB *gorm.DB
}

// GET /api/v1/events
// Lista el historial. He añadido Order para que en el monitor veas lo último arriba.
func (ctrl *BookingEventController) GetAll(c *gin.Context) {
	var items []models.BookingEvent

	// Filtro opcional: si pasas ?booking_id=10, solo muestra los de esa reserva
	query := ctrl.DB.Preload("Booking")
	bookingID := c.Query("booking_id")
	if bookingID != "" {
		query = query.Where("booking_id = ?", bookingID)
	}

	if err := query.Order("id desc").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al recuperar auditoría"})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GET /api/v1/events/:id
func (ctrl *BookingEventController) Get(c *gin.Context) {
	var item models.BookingEvent
	if err := ctrl.DB.Preload("Booking").First(&item, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Evento no encontrado"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// POST /api/v1/events
// Este método es el que llaman BookingController y RideController internamente
func (ctrl *BookingEventController) POST(c *gin.Context) {
	var item models.BookingEvent
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Estructura de evento inválida"})
		return
	}

	if err := ctrl.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo registrar el hito"})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// PUT /api/v1/events/:id
// Útil para rectificaciones manuales del Admin sobre notas previas
func (ctrl *BookingEventController) PUT(c *gin.Context) {
	var item models.BookingEvent
	id := c.Param("id")

	if err := ctrl.DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Evento no hallado"})
		return
	}

	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de actualización incorrectos"})
		return
	}

	ctrl.DB.Save(&item)
	c.JSON(http.StatusOK, item)
}

// DELETE /api/v1/events/:id
func (ctrl *BookingEventController) DELETE(c *gin.Context) {
	if err := ctrl.DB.Delete(&models.BookingEvent{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el registro de log"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Entrada de auditoría eliminada"})
}
