package controllers

import (
	"encoding/json"
	"net/http"
	"transfers/models"
	"transfers/utils" // Aquí es donde reside GenerateHashPassword

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ClientController struct {
	DB *gorm.DB
}

// POST: /api/v1/clients
func (ctrl *ClientController) POST(c *gin.Context) {
	var item models.Client
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// 1. Verificar disponibilidad de email
	var count int64
	ctrl.DB.Model(&models.User{}).Where("username = ?", item.Email).Count(&count)
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Este email ya tiene una cuenta activa"})
		return
	}

	tx := ctrl.DB.Begin()

	// 2. Crear el Usuario usando el nombre exacto de tu utilidad: GenerateHashPassword
	hashedPassword, err := utils.GenerateHashPassword(item.Phone)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar seguridad"})
		return
	}

	newUser := models.User{
		Username:     item.Email,
		Email:        item.Email,
		PasswordHash: hashedPassword,
		Role:         "client",
		IsActive:     true,
	}

	if err := tx.Create(&newUser).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear el acceso"})
		return
	}

	// 3. Normalizar Preferencias JSON
	var js interface{}
	if err := json.Unmarshal([]byte(item.Preferences), &js); err != nil {
		newData, _ := json.Marshal(map[string]string{"notes": item.Preferences})
		item.Preferences = string(newData)
	}

	// 4. Crear el Cliente vinculado al nuevo Usuario
	item.ID = 0
	item.UserID = newUser.ID

	if err := tx.Create(&item).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear la ficha"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{
		"message": "Cliente registrado",
		"data": gin.H{
			"username": newUser.Username,
			"client":   item.FullName,
		},
	})
}

// GET: /api/v1/clients
func (ctrl *ClientController) GetAll(c *gin.Context) {
	var items []models.Client
	if err := ctrl.DB.Preload("User").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al listar"})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GET: /api/v1/clients/:id
func (ctrl *ClientController) Get(c *gin.Context) {
	var item models.Client
	if err := ctrl.DB.Preload("User").First(&item, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No encontrado"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// PUT: /api/v1/clients/:id
func (ctrl *ClientController) PUT(c *gin.Context) {
	var item models.Client
	if err := ctrl.DB.First(&item, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No encontrado"})
		return
	}

	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// Normalizar preferencias en la actualización
	var js interface{}
	if err := json.Unmarshal([]byte(item.Preferences), &js); err != nil {
		newData, _ := json.Marshal(map[string]string{"notes": item.Preferences})
		item.Preferences = string(newData)
	}

	ctrl.DB.Save(&item)
	c.JSON(http.StatusOK, item)
}

// DELETE: /api/v1/clients/:id
func (ctrl *ClientController) DELETE(c *gin.Context) {
	if err := ctrl.DB.Delete(&models.Client{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Eliminado"})
}
