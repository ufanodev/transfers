package controllers

import (
	"fmt"
	"net/http"
	"transfers/models"
	"transfers/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DriverController struct {
	DB *gorm.DB
}

// GET /api/v1/drivers - Listar todos los conductores
func (ctrl *DriverController) GetAll(c *gin.Context) {
	var items []models.Driver
	// Preload de User para login y Company para saber a qué flota pertenece
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

// POST /api/v1/drivers - Crear conductor + Crear acceso de usuario
func (ctrl *DriverController) POST(c *gin.Context) {
	var item models.Driver
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	tx := ctrl.DB.Begin()

	// 1. Verificar si el usuario ya existe por email
	var existingUser models.User
	err := tx.Where("email = ?", item.Email).First(&existingUser).Error

	var targetUserID uint

	if err == nil {
		targetUserID = existingUser.ID
		fmt.Println("[DRIVER] Usando usuario existente ID:", targetUserID)
	} else {
		// Crear usuario nuevo (Password inicial = Teléfono)
		hashedPassword, _ := utils.GenerateHashPassword(item.Phone)
		newUser := models.User{
			Username:     item.Email,
			Email:        item.Email,
			PasswordHash: hashedPassword,
			Role:         "driver",
			IsActive:     true,
		}

		if err := tx.Create(&newUser).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear credenciales"})
			return
		}
		targetUserID = newUser.ID
	}

	// 2. Vincular y crear ficha de conductor
	item.ID = 0
	item.UserID = targetUserID

	if err := tx.Create(&item).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusConflict, gin.H{"error": "Este conductor ya está registrado en el sistema"})
		return
	}

	tx.Commit()
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

	// Actualizar solo la tabla drivers
	if err := ctrl.DB.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar"})
		return
	}

	c.JSON(http.StatusOK, item)
}

// DELETE /api/v1/drivers/:id - Eliminar conductor
func (ctrl *DriverController) DELETE(c *gin.Context) {
	id := c.Param("id")

	// Borrado lógico del conductor
	if err := ctrl.DB.Delete(&models.Driver{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Conductor eliminado correctamente"})
}
