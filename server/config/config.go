package config

import (
	"os"
)

// Конфигурация для подключения к базе данных
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

var AppConfig Config

// Загрузка конфигурации из переменных окружения
func LoadConfig() {
	AppConfig = Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "ozzcafe_db_user"),
		DBPassword: getEnv("DB_PASSWORD", "<password>"), // Вставь свой пароль
		DBName:     getEnv("DB_NAME", "ozzcafe_db"),
	}
}

// Получение значения из переменной окружения или использование значения по умолчанию
func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
