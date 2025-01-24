package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"ozzcafe/server/dal"
	"ozzcafe/server/models"
	"ozzcafe/server/service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Функция для настройки базы данных
func setupTestDB() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=admin dbname=ozzcafe_test port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Миграция базы данных для тестов
	db.AutoMigrate(&models.MenuItem{})
	return db, nil
}

// Функция для очистки базы данных после тестов
func teardownTestDB(db *gorm.DB) {
	db.Exec("DELETE FROM menu_items")
}

func TestCreateMenuItemHandlerIntegration(t *testing.T) {
	// Настройка тестовой базы данных
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Error setting up test DB: %v", err)
	}
	defer teardownTestDB(db)

	menuDal := dal.NewMenuDAL(db)
	// Создаем сервис с реальной базой данных
	menuService := service.NewMenuService(menuDal)
	menuHandler := NewMenuHandler(menuService)

	// Пример данных для нового блюда
	menuItem := models.MenuItem{
		Name:        "Cappuccino",
		Description: "Delicious coffee with foam",
		Price:       3.0,
		Ingredients: []string{"Espresso", "Milk", "Foam"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Преобразуем данные в JSON
	requestBody, err := json.Marshal(menuItem)
	if err != nil {
		t.Fatalf("Error marshalling request body: %v", err)
	}

	// Создаем новый HTTP-запрос
	req := httptest.NewRequest(http.MethodPost, "/menu-items", bytes.NewBuffer(requestBody))
	w := httptest.NewRecorder()

	// Вызов обработчика
	menuHandler.CreateMenuItemHandler(w, req)

	// Проверка, что статус ответа соответствует ожиданиям
	assert.Equal(t, http.StatusCreated, w.Code)

	// Проверка, что ответ содержит нужное сообщение
	expectedMessage := `{"message": "Menu item created successfully"}`
	assert.JSONEq(t, expectedMessage, w.Body.String())

	// Проверка, что меню действительно добавлено в базу данных
	var createdMenuItem models.MenuItem
	if err := db.First(&createdMenuItem, "name = ?", "Cappuccino").Error; err != nil {
		t.Fatalf("Failed to find the menu item in the database: %v", err)
	}

	// Проверка, что данные в базе соответствуют отправленным
	assert.Equal(t, "Cappuccino", createdMenuItem.Name)
	assert.Equal(t, 3.0, createdMenuItem.Price)
	assert.Equal(t, []string{"Espresso", "Milk", "Foam"}, createdMenuItem.Ingredients)
}
