package handlers

import (
	"encoding/json"
	"net/http"
	"ozzcafe/server/models"
	"ozzcafe/server/service"
	"strconv"

	"log"

	"github.com/gorilla/mux"
)

// MenuHandler обрабатывает запросы для управления меню
type MenuHandler struct {
	MenuService *service.MenuService
}

// NewMenuHandler создает новый объект MenuHandler
func NewMenuHandler(menuService *service.MenuService) *MenuHandler {
	return &MenuHandler{MenuService: menuService}
}

// GetAllMenuItemsHandler возвращает список всех блюд
func (h *MenuHandler) GetAllMenuItemsHandler(w http.ResponseWriter, r *http.Request) {
	menuItems, err := h.MenuService.GetAllMenuItems()
	if err != nil {
		log.Printf("Error fetching menu items: %v", err)
		http.Error(w, "Failed to fetch menu items", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(menuItems); err != nil {
		log.Printf("Error encoding menu items response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// CreateMenuItemHandler добавляет новое блюдо
func (h *MenuHandler) CreateMenuItemHandler(w http.ResponseWriter, r *http.Request) {
	var menuItem models.MenuItem
	// Декодируем JSON из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&menuItem); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Убедитесь, что ingredients парсятся как массив строк
	log.Printf("Decoded Ingredients: %v", menuItem.Ingredients)

	// Теперь сохранение в базу данных
	if err := h.MenuService.CreateMenuItem(&menuItem); err != nil {
		log.Printf("Error creating menu item: %v", err)
		http.Error(w, "Failed to create menu item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "Menu item created successfully"}`))
}

// UpdateMenuItemHandler обновляет информацию о блюде
func (h *MenuHandler) UpdateMenuItemHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	menuItemID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Invalid menu item ID: %v", err)
		http.Error(w, "Invalid menu item ID", http.StatusBadRequest)
		return
	}

	var menuItem models.MenuItem
	if err := json.NewDecoder(r.Body).Decode(&menuItem); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	menuItem.ID = uint(menuItemID)

	if err := h.MenuService.UpdateMenuItem(&menuItem); err != nil {
		log.Printf("Error updating menu item (ID %d): %v", menuItemID, err)
		http.Error(w, "Failed to update menu item", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(`{"message": "Menu item updated successfully"}`))
}

// DeleteMenuItemHandler удаляет блюдо
func (h *MenuHandler) DeleteMenuItemHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	menuItemID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Invalid menu item ID: %v", err)
		http.Error(w, "Invalid menu item ID", http.StatusBadRequest)
		return
	}

	if err := h.MenuService.DeleteMenuItem(uint(menuItemID)); err != nil {
		log.Printf("Error deleting menu item (ID %d): %v", menuItemID, err)
		http.Error(w, "Failed to delete menu item", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(`{"message": "Menu item deleted successfully"}`))
}
