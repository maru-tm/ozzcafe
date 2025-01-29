package database

import (
	"fmt"
	"log"
	"ozzcafe/server/config"
	"ozzcafe/server/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Global variable to hold the database connection
var DB *gorm.DB

// ConnectDatabase initializes the database connection
func ConnectDatabase() {
	// Формируем строку подключения с данными из конфигурации
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.AppConfig.DBHost,
		config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBName,
		config.AppConfig.DBPort)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	// Миграция таблиц (если нужно)
	MigrateTables()

	log.Println("Successfully connected to the database")
}

// MigrateTables performs migrations for necessary tables
func MigrateTables() {
	if err := DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Error creating User table: %v", err)
	}
	log.Println("User table has been migrated")

	if err := DB.AutoMigrate(&models.MenuItem{}); err != nil {
		log.Fatalf("Error migrating MenuItem or MenuIngredient models: %v", err)
	}
	log.Println("MenuItem table has been migrated")
}

// GetDB returns the current database connection
func GetDB() *gorm.DB {
	return DB
}
