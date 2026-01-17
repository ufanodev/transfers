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

const AuthCookieName = "transfers_session"

// --- Seguridad Base ---

func GenerateHashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// --- Generador de Token y Cookie ---
// Esta es la función que faltaba y causaba el error undefined
func GenerateAuthCookie(userID uint, role string) (*http.Cookie, error) {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return nil, fmt.Errorf("JWT_SECRET no configurado")
	}

	expirationTime := time.Now().Add(time.Hour * 24)

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
		HttpOnly: true,
		Secure:   false, // Cambiar a true si usas HTTPS en producción
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}, nil
}

// --- Middlewares de Auditoría y Control ---

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		secretKey := os.Getenv("JWT_SECRET")
		now := time.Now().Format("2006-01-02 15:04:05")
		path := c.Request.URL.Path
		method := c.Request.Method

		cookie, err := c.Request.Cookie(AuthCookieName)
		if err != nil {
			fmt.Printf("[%s] ⚠️  ACCESO ANÓNIMO | %s %s | Motivo: No hay cookie\n", now, method, path)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Inicie sesión"})
			return
		}

		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			fmt.Printf("[%s] 🚫 TOKEN INVÁLIDO | %s %s | Error: %v\n", now, method, path, err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Sesión inválida"})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			uid := uint(claims["user_id"].(float64))
			role := claims["role"].(string)

			fmt.Printf("[%s] ✅ LLAMADA | %s %s | UserID: %d | Rol: [%s]\n", now, method, path, uid, role)

			c.Set("userID", uid)
			c.Set("userRole", role)
			c.Next()
		}
	}
}

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, _ := c.Get("userRole")
		userID, _ := c.Get("userID")
		now := time.Now().Format("15:04:05")
		path := c.Request.URL.Path

		isAllowed := false
		for _, role := range allowedRoles {
			if role == userRole {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			fmt.Printf("[%s] 🛑 BLOQUEO DE ROL | User: %v [%s] intentó entrar a %s (Permitido solo para: %v)\n", now, userID, userRole, path, allowedRoles)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Acceso denegado"})
			return
		}
		c.Next()
	}
}
