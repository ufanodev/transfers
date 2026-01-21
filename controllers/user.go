package controllers

import (
	"net/http"
	"transfers/models"
	"transfers/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

// GetMe: Obtiene el perfil del usuario autenticado actual
// Soluciona el error de los logs: GET /api/v1/users/me
func (ctrl *UserController) GetMe(c *gin.Context) {
	// Extraemos el ID del usuario del contexto (inyectado por JWTAuthMiddleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesión no válida"})
		return
	}

	var user models.User
	// Buscamos por el ID numérico real, no por el string "me"
	if err := ctrl.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetAll: Lista todos los usuarios (Solo para Admin)
func (ctrl *UserController) GetAll(c *gin.Context) {
	var users []models.User
	if err := ctrl.DB.Order("id desc").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al listar usuarios"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// Get: Obtiene un usuario por ID
func (ctrl *UserController) Get(c *gin.Context) {
	var user models.User
	id := c.Param("id")
	if err := ctrl.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no hallado"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// POST: Crea un nuevo usuario
func (ctrl *UserController) POST(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Encriptar contraseña antes de guardar
	hashedPassword, _ := utils.GenerateHashPassword(user.PasswordHash) // Asumiendo que viene plana en el JSON
	user.PasswordHash = hashedPassword

	if err := ctrl.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "El email o usuario ya existe"})
		return
	}
	c.JSON(http.StatusCreated, user)
}

// PUT: Actualiza un usuario
func (ctrl *UserController) PUT(c *gin.Context) {
	var user models.User
	id := c.Param("id")
	if err := ctrl.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	ctrl.DB.Save(&user)
	c.JSON(http.StatusOK, user)
}

// DELETE: Borrado de usuario
func (ctrl *UserController) DELETE(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.DB.Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Usuario eliminado"})
}
