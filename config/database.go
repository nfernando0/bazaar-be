package config

import (
	"fmt"
	"log"
	"os"

	"bazar-be/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// ConnectDB initializes the database connection and runs auto migrations
func ConnectDB() {
	// Load .env file if available
	if err := godotenv.Load(); err != nil {
		log.Println("Note: .env file not found, reading from system environment variables")
	}

	dbHost := getEnv("DB_HOST", "127.0.0.1")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	dbPass := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "bazar_db")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v\nDSN: %s", err, dsn)
	}

	log.Println("Database connection established successfully.")

	// Auto Migration
	if err := DB.AutoMigrate(
		&models.Bazaar{},
		&models.Outlet{},
		&models.Vendor{},
		&models.VendorOutlet{},
		&models.Category{},
		&models.Product{},
		&models.User{},
		&models.UserToken{},
		&models.Transaction{},
		&models.TransactionItem{},
	); err != nil {
		log.Fatalf("Auto migration failed: %v", err)
	}

	log.Println("Database migration completed.")
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
