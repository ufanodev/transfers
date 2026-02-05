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

// GetMe: Obtiene el perfil del usuario autenticado actual mediante el token JWT
func (ctrl *UserController) GetMe(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesión no válida o expirada"})
		return
	}

	var user models.User
	if err := ctrl.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	// Omitimos el hash de la contraseña por seguridad en la respuesta
	user.PasswordHash = ""
	c.JSON(http.StatusOK, user)
}

// GetAll: Lista todos los usuarios (Reservado para Admin)
func (ctrl *UserController) GetAll(c *gin.Context) {
	var users []models.User
	if err := ctrl.DB.Order("id desc").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al listar usuarios"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// Get: Obtiene un usuario específico por ID
func (ctrl *UserController) Get(c *gin.Context) {
	var user models.User
	id := c.Param("id")
	if err := ctrl.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no hallado"})
		return
	}
	user.PasswordHash = "" // Seguridad
	c.JSON(http.StatusOK, user)
}

// POST: Crea un nuevo usuario y encripta su contraseña
func (ctrl *UserController) POST(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de usuario inválidos"})
		return
	}

	// Encriptar contraseña: utils.GenerateHashPassword debe manejar el coste de Bcrypt
	hashedPassword, err := utils.GenerateHashPassword(user.PasswordHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar credenciales"})
		return
	}
	user.PasswordHash = hashedPassword

	if err := ctrl.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "El email o nombre de usuario ya está registrado"})
		return
	}

	user.PasswordHash = ""
	c.JSON(http.StatusCreated, user)
}

// PUT: Actualiza datos de usuario (sin modificar password a menos que se solicite)
func (ctrl *UserController) PUT(c *gin.Context) {
	var user models.User
	id := c.Param("id")

	if err := ctrl.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de actualización incorrectos"})
		return
	}

	// Lógica de actualización selectiva para no pisar campos vacíos o el ID
	ctrl.DB.Model(&user).Updates(models.User{
		Username: input.Username,
		Email:    input.Email,
		Role:     input.Role,
	})

	c.JSON(http.StatusOK, user)
}

// DELETE: Borrado físico o lógico (según configuración de models.User)
func (ctrl *UserController) DELETE(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.DB.Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo eliminar el usuario"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Usuario eliminado del sistema"})
}
