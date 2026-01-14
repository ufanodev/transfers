package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Endpoint de prueba (Health Check)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
			"status":  "ready",
		})
	})

	// Aquí iremos añadiendo grupos de rutas
	// api := r.Group("/api")
	// {
	//    api.POST("/login", controllers.Login)
	// }

	return r
}
