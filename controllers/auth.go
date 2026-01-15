/**
 * @file auth.go
 * @description Controlador de autenticación. Gestiona el inicio de sesión,
 * la creación de cookies de sesión y la redirección de vistas por rol.
 */

package controllers

import (
	"net/http"
	"transfers/models"
	"transfers/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

// Login procesa las credenciales del usuario y establece la Cookie HttpOnly
func (ctrl *AuthController) Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// 1. Validar formato de entrada JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada inválidos"})
		return
	}

	// 2. Buscar usuario en la base de datos por email
	var user models.User
	if err := ctrl.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales incorrectas"})
		return
	}

	// 3. Verificar si el usuario está activo
	if !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cuenta de usuario desactivada"})
		return
	}

	// 4. Comparar hashes de contraseña usando la utilidad centralizada
	if !utils.CheckPasswordHash(input.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales incorrectas"})
		return
	}

	// 5. Generar la Cookie con el Token JWT (Contiene ID y Rol)
	cookie, err := utils.GenerateAuthCookie(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al generar la sesión segura"})
		return
	}

	// 6. Establecer la cookie en la respuesta HTTP
	http.SetCookie(c.Writer, cookie)

	// Respuesta de éxito para el JS (login.js)
	c.JSON(http.StatusOK, gin.H{
		"message": "Autenticación exitosa",
		"role":    user.Role,
	})
}

// RedirectByRole sirve el archivo HTML físico basándose en el rol del usuario.
// Este método es llamado por la ruta /dashboard/ después de pasar el Middleware.
func (ctrl *AuthController) RedirectByRole(c *gin.Context) {
	// Obtenemos el rol inyectado por el Middleware utils.JWTAuthMiddleware()
	role, exists := c.Get("userRole")

	if !exists {
		// Si no hay rol en el contexto, forzamos vuelta al login
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	// Servimos el archivo HTML correspondiente desde la carpeta static/views
	// Go lee el archivo del disco y lo envía directamente al navegador
	switch role {
	case "admin":
		c.File("./static/views/admin.html")
	case "client":
		c.File("./static/views/clients.html")
	case "company":
		c.File("./static/views/companies.html")
	case "driver":
		c.File("./static/views/drivers.html")
	default:
		// Rol desconocido: limpieza de seguridad y redirección
		utils.ClearAuthCookie(c)
		c.Redirect(http.StatusSeeOther, "/")
	}
}
