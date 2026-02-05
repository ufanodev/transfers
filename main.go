package main

import (
	"fmt"
	"log"
	"os"
	"transfers/config"
	"transfers/models"
	"transfers/routes"
	"transfers/utils"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	// 1. Encabezado de inicio
	fmt.Println("###########################################")
	fmt.Println("#       TAKEUS SYSTEM - LOGISTICS         #")
	fmt.Println("###########################################")

	// 2. Carga del archivo .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  [Config] .env no detectado, usando variables de entorno.")
	}

	// 3. Inicialización de la Base de Datos
	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatalf("❌ [FATAL] Error de conexión: %v", err)
	}

	// 4. AUTOMIGRATE (Orden Crítico para Foreign Keys)
	// Primero los maestros, luego la lógica, al final finanzas y auditoría.
	log.Println("🗂️  [Database] Sincronizando tablas...")
	err = db.AutoMigrate(
		&models.User{},
		&models.Company{},
		&models.Client{},
		&models.Driver{},
		&models.Vehicle{},
		&models.Booking{},      // Debe ir antes que Events y Rides
		&models.BookingEvent{}, // Depende de Booking
		&models.Ride{},         // Depende de Booking, Driver y Vehicle
		&models.Payment{},      // Depende de Booking
		&models.Rating{},
	)
	if err != nil {
		log.Fatalf("❌ [FATAL] Error en migración: %v", err)
	}
	log.Println("✅ [Database] Tablas sincronizadas con éxito.")

	// 5. Semilla de Administrador Inicial
	createInitialAdmin(db)

	// 6. Configurar el Router de Gin
	r := routes.SetupRouter(db)

	// 7. Preparación y Arranque del Servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 [Server] Backend TakeUs listo en http://localhost:%s", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("❌ [FATAL] No se pudo iniciar el servidor: %v", err)
	}
}

func createInitialAdmin(db *gorm.DB) {
	var count int64
	db.Model(&models.User{}).Count(&count)

	if count == 0 {
		log.Println("👤 [Setup] Creando administrador por defecto...")

		pass := "Ies11887010!"
		hashedPassword, err := utils.GenerateHashPassword(pass)
		if err != nil {
			log.Printf("❌ [Error] Falló encriptación: %v", err)
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
			log.Println("✅ [Setup] Admin inicial: admin@transfers.com / Ies11887010!")
		}
	}
}
