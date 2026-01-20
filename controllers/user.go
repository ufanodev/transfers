package controllers

import (
	"net/http"
	"transfers/models"
	"transfers/utils" // Asegúrate de que aquí esté GenerateHashPassword

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserController gestiona las peticiones HTTP para el modelo User
type UserController struct {
	DB *gorm.DB
}

// GetMe - GET /api/v1/users/me
// Obtiene el perfil del usuario logueado mediante el token JWT
func (ctrl *UserController) GetMe(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesión expirada"})
		return
	}

	var user models.User
	if err := ctrl.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":        user.ID,
		"username":  user.Username,
		"email":     user.Email,
		"role":      user.Role,
		"is_active": user.IsActive,
	})
}

// GetAll - GET /api/v1/users
// Lista todos los usuarios de la base de datos
func (ctrl *UserController) GetAll(c *gin.Context) {
	var users []models.User
	if err := ctrl.DB.Order("id desc").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al recuperar usuarios"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// Get - GET /api/v1/users/:id
// Obtiene un usuario específico por ID para cargar el formulario de edición
func (ctrl *UserController) Get(c *gin.Context) {
	var user models.User
	id := c.Param("id")

	if err := ctrl.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no hallado"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// POST - POST /api/v1/users
// Crea un nuevo usuario con contraseña encriptada
func (ctrl *UserController) POST(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
		Role     string `json:"role" binding:"required"`
		IsActive bool   `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos incompletos o inválidos"})
		return
	}

	// Verificar si el email ya existe
	var existing models.User
	if err := ctrl.DB.Where("email = ?", input.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "El email ya está registrado"})
		return
	}

	// Encriptar password
	hashed, err := utils.GenerateHashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fallo de seguridad"})
		return
	}

	newUser := models.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: hashed,
		Role:         input.Role,
		IsActive:     input.IsActive,
	}

	if err := ctrl.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear el registro"})
		return
	}

	c.JSON(http.StatusCreated, newUser)
}

// PUT - PUT /api/v1/users/:id
// Actualiza un usuario existente. Si el password viene vacío, no se modifica.
func (ctrl *UserController) PUT(c *gin.Context) {
	var user models.User
	id := c.Param("id")

	if err := ctrl.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
		IsActive *bool  `json:"is_active"` // Puntero para detectar booleanos false
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al procesar los datos"})
		return
	}

	// Actualización selectiva
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

	// Si el admin escribió una nueva contraseña, se re-encripta
	if input.Password != "" {
		hashed, _ := utils.GenerateHashPassword(input.Password)
		user.PasswordHash = hashed
	}

	if err := ctrl.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar cambios"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DELETE - DELETE /api/v1/users/:id
// Elimina un usuario de la base de datos
func (ctrl *UserController) DELETE(c *gin.Context) {
	id := c.Param("id")

	// Usamos Unscoped() si quieres borrarlo físicamente, o normal para soft-delete de GORM
	if err := ctrl.DB.Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo eliminar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuario eliminado correctamente"})
}
