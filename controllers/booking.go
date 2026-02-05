package controllers

import (
	"fmt"
	"net/http"
	"time"
	"transfers/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BookingController struct {
	DB *gorm.DB
}

func (ctrl *BookingController) GetAll(c *gin.Context) {
	var items []models.Booking
	if err := ctrl.DB.Preload("Client").Order("scheduled_at asc").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al listar reservas"})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (ctrl *BookingController) Get(c *gin.Context) {
	var item models.Booking
	if err := ctrl.DB.Preload("Client").First(&item, c.Param("id")).Error; err != nil {
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

	timeStr := fmt.Sprintf("%s %s", item.ScheduledDate.Format("2006-01-02"), item.ScheduledTime)
	if t, err := time.Parse("2006-01-02 15:04", timeStr); err == nil {
		item.ScheduledAt = t
	}

	tx := ctrl.DB.Begin()
	if err := tx.Create(&item).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear reserva"})
		return
	}

	tx.Create(&models.BookingEvent{
		BookingID:   item.ID,
		EventType:   "created",
		Description: fmt.Sprintf("Reserva creada para %s", item.ClientName),
		NewValue:    "pending",
		CreatedBy:   "Admin System",
	})

	tx.Commit()
	c.JSON(http.StatusCreated, item)
}

func (ctrl *BookingController) PUT(c *gin.Context) {
	var item models.Booking
	id := c.Param("id")

	if err := ctrl.DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reserva no encontrada"})
		return
	}

	oldStatus := item.Status
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := ctrl.DB.Begin()

	// --- TRASPASO AUTOMÁTICO A OPERATIVA (RIDE) ---
	if oldStatus != "confirmed" && item.Status == "confirmed" {
		var count int64
		tx.Model(&models.Ride{}).Where("booking_id = ?", item.ID).Count(&count)

		if count == 0 {
			newRide := models.Ride{
				BookingID:      item.ID,
				ClientName:     item.ClientName,
				ClientPhone:    item.ClientPhone,
				PickupAddress:  item.OriginAddress,
				DropoffAddress: item.DestAddress,
				OriginLat:      item.OriginLat,
				OriginLng:      item.OriginLng,
				DestinationLat: item.DestLat,
				DestinationLng: item.DestLng,
				Pax:            item.Pax,
				LuggageCount:   item.Maletas,
				Animals:        item.Animal,
				Status:         "scheduled",
			}
			// DriverID y VehicleID quedan como nil (NULL en DB)
			if err := tx.Create(&newRide).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Fallo al crear Ride: " + err.Error()})
				return
			}

			tx.Create(&models.BookingEvent{
				BookingID:   item.ID,
				EventType:   "transfer_to_ops",
				Description: "Reserva confirmada. Enviada al monitor de viajes.",
				NewValue:    "scheduled",
				CreatedBy:   "Logistics System",
			})
		}
	}

	if oldStatus != item.Status {
		tx.Create(&models.BookingEvent{
			BookingID:   item.ID,
			EventType:   "status_change",
			Description: fmt.Sprintf("Cambio de estado: %s -> %s", oldStatus, item.Status),
			OldValue:    oldStatus,
			NewValue:    item.Status,
			CreatedBy:   "Admin Editor",
		})
	}

	if err := tx.Save(&item).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, item)
}

func (ctrl *BookingController) DELETE(c *gin.Context) {
	id := c.Param("id")
	ctrl.DB.Delete(&models.Booking{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Reserva eliminada"})
}
