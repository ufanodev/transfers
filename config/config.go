package config

import (
	"fmt"
	"log"
	"os"
	"transfers/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() (*gorm.DB, error) {
	// Intentar cargar .env pero no morir si falla
	_ = godotenv.Load()

	// Obtener variables con valores de respaldo (Coincidentes con tu Docker)
	dbUser := getEnv("DB_USER", "transfers_user")
	dbPass := getEnv("DB_PASSWORD", "Ies11887010!")
	dbHost := getEnv("DB_HOST", "127.0.0.1") // localhost
	dbPort := getEnv("DB_PORT", "3306")
	dbName := getEnv("DB_NAME", "transfers")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	log.Printf("🔗 [DB] Intentando conectar a %s:%s...", dbHost, dbPort)

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Migraciones automáticas
	log.Println("🛠️ [DB] Sincronizando tablas...")
	database.AutoMigrate(
		&models.User{}, &models.Client{}, &models.Company{},
		&models.Driver{}, &models.Vehicle{}, &models.Booking{},
		&models.BookingEvent{}, &models.Ride{}, &models.Payment{}, &models.Rating{},
	)

	DB = database
	return database, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
