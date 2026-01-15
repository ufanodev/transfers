/**
 * @file auth.go
 * @package utils
 * @description Utilidades de seguridad. Contiene la lógica de encriptación bcrypt,
 * generación/limpieza de cookies JWT y el Middleware de autenticación.
 */

package utils

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Nombre de la cookie que almacenará el token
const AuthCookieName = "transfers_session"

// ---------------------------------------------------------------------
// --- Utilidades de Hashing
// ---------------------------------------------------------------------

// GenerateHashPassword convierte texto plano en hash bcrypt
func GenerateHashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compara contraseña vs hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ---------------------------------------------------------------------
// --- Gestión de Cookies y JWT
// ---------------------------------------------------------------------

// GenerateAuthCookie crea una cookie HttpOnly con claims de usuario y rol
func GenerateAuthCookie(userID uint, role string) (*http.Cookie, error) {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return nil, fmt.Errorf("JWT_SECRET no configurado en .env")
	}

	expirationTime := time.Now().Add(time.Hour * 24) // 24 horas de validez

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     expirationTime.Unix(),
	})

	tokenString, err := claims.SignedString([]byte(secretKey))
	if err != nil {
		return nil, err
	}

	return &http.Cookie{
		Name:     AuthCookieName,
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: true,  // Protege contra XSS
		Secure:   false, // Cambiar a true en producción con HTTPS
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}, nil
}

// ClearAuthCookie elimina la cookie de sesión
func ClearAuthCookie(c *gin.Context) {
	c.SetCookie(AuthCookieName, "", -1, "/", "", false, true)
}

// ---------------------------------------------------------------------
// --- Middleware de Autenticación
// ---------------------------------------------------------------------

// JWTAuthMiddleware valida la cookie en cada petición protegida
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		secretKey := os.Getenv("JWT_SECRET")

		cookie, err := c.Request.Cookie(AuthCookieName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Sesión no iniciada"})
			return
		}

		// Parsear y validar el token
		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("método de firma inesperado")
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Sesión inválida o expirada"})
			return
		}

		// Inyectar claims en el contexto de Gin para los controladores
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("userID", uint(claims["user_id"].(float64)))
			c.Set("userRole", claims["role"].(string))
			c.Next()
		}
	}
}
