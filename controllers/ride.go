package controllers

import (
	"net/http"
	"transfers/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RideController struct {
	DB *gorm.DB
}

// GET /api/v1/rides
// Lista todos los servicios con información cruzada detallada
func (ctrl *RideController) GetAll(c *gin.Context) {
	var items []models.Ride
	// Preload carga las relaciones Driver, Vehicle y Booking para visualización completa
	if err := ctrl.DB.Preload("Driver").Preload("Vehicle").Preload("Booking").Order("id desc").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al recuperar monitor de viajes"})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GET /api/v1/rides/:id
// Obtiene el detalle completo de una carrera específica
func (ctrl *RideController) Get(c *gin.Context) {
	var item models.Ride
	id := c.Param("id")
	if err := ctrl.DB.Preload("Driver").Preload("Vehicle").Preload("Booking").First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Viaje no encontrado"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// POST /api/v1/rides
// Crea la asignación logística inicial
func (ctrl *RideController) POST(c *gin.Context) {
	var item models.Ride
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de asignación inválidos"})
		return
	}

	// Validación de seguridad: Verificar que el BookingID no esté duplicado si se requiere
	tx := ctrl.DB.Begin()

	if err := tx.Create(&item).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusConflict, gin.H{"error": "Error al crear el viaje. Verifique si la reserva ya tiene un servicio asignado."})
		return
	}

	tx.Commit()
	c.JSON(http.StatusCreated, item)
}

// PUT /api/v1/rides/:id
// Actualiza métricas en tiempo real (GPS, Km, Tiempos y Cierre)
func (ctrl *RideController) PUT(c *gin.Context) {
	var item models.Ride
	id := c.Param("id")

	// 1. Verificar existencia
	if err := ctrl.DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Registro de viaje no hallado"})
		return
	}

	// 2. Vincular nuevos datos (JSON -> Model)
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de actualización incorrectos"})
		return
	}

	// 3. Lógica de negocio: Recalcular KM Total si vienen valores parciales
	if item.KmEnd > 0 && item.KmStart > 0 {
		item.KmTotal = item.KmEnd - item.KmStart
	}

	// 4. Guardar cambios (Save actualiza todos los campos incluyendo nulos de geolocalización)
	if err := ctrl.DB.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron sincronizar los cambios en la base de datos"})
		return
	}

	c.JSON(http.StatusOK, item)
}

// DELETE /api/v1/rides/:id
// Borrado lógico (Soft Delete) para mantener historial de auditoría
func (ctrl *RideController) DELETE(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.DB.Delete(&models.Ride{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el registro"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Viaje eliminado correctamente de la lista activa"})
}
