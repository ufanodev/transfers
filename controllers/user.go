package controllers

import (
	"net/http"
	"transfers/models"
	"transfers/utils" // Importante para el hashing

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserController define el controlador para el modelo User
type UserController struct {
	DB *gorm.DB
}

// GetAll - GET /api/v1/users
func (ctrl *UserController) GetAll(c *gin.Context) {
	var items []models.User
	if err := ctrl.DB.Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al recuperar usuarios"})
		return
	}
	c.JSON(http.StatusOK, items)
}

// Get - GET /api/v1/users/:id
func (ctrl *UserController) Get(c *gin.Context) {
	var item models.User
	if err := ctrl.DB.First(&item, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// POST - POST /api/v1/users (Creación)
func (ctrl *UserController) POST(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
		Role     string `json:"role" binding:"required"`
		IsActive bool   `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Encriptar contraseña usando tu utilidad
	hashedPassword, err := utils.GenerateHashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar contraseña"})
		return
	}

	newUser := models.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: hashedPassword,
		Role:         input.Role,
		IsActive:     input.IsActive,
	}

	if err := ctrl.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear el usuario"})
		return
	}

	c.JSON(http.StatusCreated, newUser)
}

// PUT - PUT /api/v1/users/:id (Actualización)
func (ctrl *UserController) PUT(c *gin.Context) {
	var user models.User
	if err := ctrl.DB.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"` // Opcional en update
		Role     string `json:"role"`
		IsActive *bool  `json:"is_active"` // Puntero para detectar si viene el campo false
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Actualizar campos básicos si vienen en el JSON
	if input.Username != "" {
		user.Username = input.Username
	}
	if input.Email != "" {
		user.Email = input.Email
	}
	if input.Role != "" {
		user.Role = input.Role
	}
	if input.IsActive != nil {
		user.IsActive = *input.IsActive
	}

	// Si el usuario envía una nueva contraseña, la encriptamos
	if input.Password != "" {
		hashed, _ := utils.GenerateHashPassword(input.Password)
		user.PasswordHash = hashed
	}

	if err := ctrl.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DELETE - DELETE /api/v1/users/:id
func (ctrl *UserController) DELETE(c *gin.Context) {
	if err := ctrl.DB.Delete(&models.User{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Usuario eliminado"})
}
