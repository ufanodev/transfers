/**
 * @file auth.go
 * @description Controlador de autenticación. Gestiona el acceso y la
 * redirección directa a rutas limpias (/admin, /drivers, /clients).
 */

package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"transfers/models"
	"transfers/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

// Login valida las credenciales y establece la cookie de sesión segura
func (ctrl *AuthController) Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// 1. Validar JSON entrante
	if err := c.ShouldBindJSON(&input); err != nil {
		fmt.Println("[LOGIN] Error: Datos de entrada inválidos")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de datos inválido"})
		return
	}

	// 2. Buscar usuario por Email
	var user models.User
	if err := ctrl.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		fmt.Printf("[LOGIN] Fallido: Email %s no existe\n", input.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales incorrectas"})
		return
	}

	// 3. Verificar estado de cuenta
	if !user.IsActive {
		fmt.Printf("[LOGIN] Bloqueado: Usuario %s desactivado\n", user.Username)
		c.JSON(http.StatusForbidden, gin.H{"error": "Cuenta de usuario desactivada"})
		return
	}

	// 4. Validar Password
	if !utils.CheckPasswordHash(input.Password, user.PasswordHash) {
		fmt.Printf("[LOGIN] Fallido: Password incorrecta para %s\n", user.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales incorrectas"})
		return
	}

	// 5. Generar Cookie JWT
	cookie, err := utils.GenerateAuthCookie(user.ID, user.Role)
	if err != nil {
		fmt.Println("[LOGIN] Error generando cookie:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	// 6. Enviar cookie al cliente
	http.SetCookie(c.Writer, cookie)
	fmt.Printf("[LOGIN] Éxito: %s conectado como [%s]\n", user.Username, user.Role)

	c.JSON(http.StatusOK, gin.H{
		"message": "Autenticación exitosa",
		"role":    user.Role,
	})
}

/**
 * RedirectByRole: Redirige el navegador a la URL final limpia.
 * Se activa al acceder a /dashboard.
 */
func (ctrl *AuthController) RedirectByRole(c *gin.Context) {
	role, exists := c.Get("userRole")

	if !exists {
		fmt.Println("[ROUTING] No se encontró rol. Redirigiendo a Login.")
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	// Normalizar el rol para el switch
	roleStr := strings.ToLower(role.(string))
	fmt.Printf("[ROUTING] Usuario con rol [%s] -> Redirigiendo a su ruta raíz\n", roleStr)

	// REDIRECCIONES FÍSICAS A RUTAS LIMPIAS
	switch roleStr {
	case "admin":
		fmt.Println("[ROUTING] -> http://localhost:8080/admin")
		c.Redirect(http.StatusSeeOther, "/admin")

	case "driver":
		fmt.Println("[ROUTING] -> http://localhost:8080/drivers")
		c.Redirect(http.StatusSeeOther, "/drivers")

	case "client":
		fmt.Println("[ROUTING] -> http://localhost:8080/clients")
		c.Redirect(http.StatusSeeOther, "/clients")

	case "company":
		fmt.Println("[ROUTING] -> http://localhost:8080/companies")
		c.Redirect(http.StatusSeeOther, "/companies")

	default:
		fmt.Printf("[ROUTING] Rol desconocido: %s. Expulsando.\n", roleStr)
		utils.ClearAuthCookie(c)
		c.Redirect(http.StatusSeeOther, "/")
	}
}

// Logout elimina la sesión y redirige al inicio
func (ctrl *AuthController) Logout(c *gin.Context) {
	fmt.Println("[AUTH] Logout: Limpiando sesión.")
	utils.ClearAuthCookie(c)
	c.Redirect(http.StatusSeeOther, "/")
}
