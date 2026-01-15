package main

import (
	"fmt"
	"log"
	"os"
	"transfers/config"
	"transfers/models"
	"transfers/routes"
	"transfers/utils" // 👈 Importamos tus nuevas utilidades

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	// 1. Encabezado de inicio
	fmt.Println("###########################################")
	fmt.Println("#       TRANSFERS SYSTEM - BACKEND        #")
	fmt.Println("###########################################")

	// 2. Carga del archivo .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  [Config] .env no detectado, usando variables de entorno o defaults.")
	}

	// 3. Inicialización de la Base de Datos
	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatalf("❌ [FATAL] Error de conexión: %v", err)
	}

	// 4. Semilla de Administrador Inicial
	// Usamos la lógica de utils para que el hash sea idéntico al del Login
	createInitialAdmin(db)

	// 5. Configurar el Router de Gin
	r := routes.SetupRouter(db)

	// 6. Preparación y Arranque del Servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 [Server] Backend listo y escuchando en http://localhost:%s", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("❌ [FATAL] No se pudo iniciar el servidor: %v", err)
	}
}

func createInitialAdmin(db *gorm.DB) {
	var count int64
	db.Model(&models.User{}).Count(&count)

	if count == 0 {
		log.Println("👤 [Setup] Base de datos vacía. Creando administrador inicial...")

		// 🔑 Usamos tu utilidad centralizada de hashing
		pass := "Ies11887010!"
		hashedPassword, err := utils.GenerateHashPassword(pass)
		if err != nil {
			log.Printf("❌ [Error] Falló encriptación de pass: %v", err)
			return
		}

		admin := models.User{
			Username:     "admin",
			Email:        "admin@transfers.com",
			PasswordHash: hashedPassword,
			Role:         "admin",
			IsActive:     true,
		}

		if err := db.Create(&admin).Error; err != nil {
			log.Printf("❌ [Error] No se pudo crear el admin: %v", err)
		} else {
			log.Println("✅ [Setup] Admin creado con éxito.")
			log.Println("📧 Usuario: admin@transfers.com | 🔑 Pass: Ies11887010!")
		}
	} else {
		log.Printf("ℹ️  [Setup] Sistema listo con %d usuarios registrados.", count)
	}
}
