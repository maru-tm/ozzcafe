// main.go
package main

import (
	"log"
	"net/http"
	"ozzcafe/server/config"
	"ozzcafe/server/database"
	"ozzcafe/server/router"
)

func main() {
	// Загрузка конфигурации
	config.LoadConfig()

	// Инициализация подключения к базе данных
	database.ConnectDatabase()

	// Создание маршрутизатора с подключением к базе данных
	r := router.NewRouter(database.GetDB())

	// Запуск сервера
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
