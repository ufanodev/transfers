package controllers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"transfers/models"
	"transfers/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

// LoginRequest define la estructura estrictamente esperada del cliente
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login: Valida credenciales, comprueba email/password y genera cookie JWT
func (ctrl *AuthController) Login(c *gin.Context) {
	// --- DEPURACIÓN DEL BODY ---
	bodyBytes, _ := io.ReadAll(c.Request.Body)
	fmt.Printf("[AUTH] 🛰️ RAW BODY: '%s'\n", string(bodyBytes))
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var input LoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		fmt.Printf("[AUTH] ❌ Error de Binding JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email y contraseña requeridos"})
		return
	}

	// 1. Buscar al usuario en la DB
	var user models.User
	if err := ctrl.DB.Where("email = ?", input.Username).First(&user).Error; err != nil {
		fmt.Printf("[AUTH] 🚫 Usuario no encontrado: %s\n", input.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales inválidas"})
		return
	}

	// 2. Verificar Password
	if !utils.CheckPasswordHash(input.Password, user.PasswordHash) {
		fmt.Printf("[AUTH] 🔑 Password incorrecta: %s\n", input.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales inválidas"})
		return
	}

	// 3. Generar Cookie JWT
	cookie, err := utils.GenerateAuthCookie(user.ID, user.Role)
	if err != nil {
		fmt.Printf("[AUTH] ❌ Error generando Cookie: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno"})
		return
	}

	// 4. Inyectar Cookie en respuesta
	http.SetCookie(c.Writer, cookie)

	fmt.Printf("[AUTH] ✅ LOGIN EXITOSO: %s | Rol: [%s]\n", user.Email, user.Role)

	c.JSON(http.StatusOK, gin.H{
		"message": "Bienvenido",
		"role":    user.Role,
	})
}

// RedirectByRole: El semáforo que actúa tras el login exitoso (Coordinado con routes.go)
func (ctrl *AuthController) RedirectByRole(c *gin.Context) {
	role, exists := c.Get("userRole")
	if !exists {
		fmt.Println("[SEMÁFORO] ⚠️ Rol no encontrado, regresando al login")
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	fmt.Printf("[SEMÁFORO] 🚦 Redirigiendo según rol: [%v]\n", role)

	// Redirecciones en SINGULAR para coincidir con las carpetas físicas y routes.go
	switch role {
	case "admin":
		c.Redirect(http.StatusSeeOther, "/admin/")
	case "client":
		c.Redirect(http.StatusSeeOther, "/client/")
	case "driver":
		c.Redirect(http.StatusSeeOther, "/driver/")
	case "company":
		c.Redirect(http.StatusSeeOther, "/company/")
	default:
		c.Redirect(http.StatusSeeOther, "/")
	}
}

// Logout: Invalida la cookie y saca al usuario
func (ctrl *AuthController) Logout(c *gin.Context) {
	c.SetCookie(utils.AuthCookieName, "", -1, "/", "", false, true)
	fmt.Println("[AUTH] 🚪 Sesión cerrada.")
	c.Redirect(http.StatusSeeOther, "/")
}
