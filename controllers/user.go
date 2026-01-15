package controllers

import (
	"net/http"
	"transfers/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

// GET /api/v1/users - Listar todos los usuarios
func (ctrl *UserController) GetAll(c *gin.Context) {
	var items []models.User
	if err := ctrl.DB.Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al recuperar usuarios"})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GET /api/v1/users/:id - Obtener un usuario por ID
func (ctrl *UserController) Get(c *gin.Context) {
	var item models.User
	if err := ctrl.DB.First(&item, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// POST /api/v1/users - Crear usuario con encriptación de contraseña
func (ctrl *UserController) POST(c *gin.Context) {
	var item models.User
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Encriptar la contraseña antes de guardar
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(item.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar contraseña"})
		return
	}
	item.PasswordHash = string(hashedPassword)

	if err := ctrl.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear el usuario"})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// PUT /api/v1/users/:id - Actualizar usuario
func (ctrl *UserController) PUT(c *gin.Context) {
	var item models.User
	if err := ctrl.DB.First(&item, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	var updateData struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
		IsActive bool   `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Actualizar campos básicos
	item.Username = updateData.Username
	item.Email = updateData.Email
	item.Role = updateData.Role
	item.IsActive = updateData.IsActive

	// Si se envía una nueva contraseña, la encriptamos
	if updateData.Password != "" {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(updateData.Password), bcrypt.DefaultCost)
		item.PasswordHash = string(hashed)
	}

	ctrl.DB.Save(&item)
	c.JSON(http.StatusOK, item)
}

// DELETE /api/v1/users/:id - Borrado físico o lógico según el modelo
func (ctrl *UserController) DELETE(c *gin.Context) {
	if err := ctrl.DB.Delete(&models.User{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar usuario"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Usuario eliminado correctamente"})
}
