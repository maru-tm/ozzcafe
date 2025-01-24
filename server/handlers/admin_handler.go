package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"ozzcafe/server/service"
	"strconv"

	"github.com/gorilla/mux"
)

// AdminHandler обрабатывает запросы для администраторов
type AdminHandler struct {
	AdminService *service.AdminService
}

// NewAdminHandler создает новый обработчик для администратора
func NewAdminHandler(adminService *service.AdminService) *AdminHandler {
	return &AdminHandler{AdminService: adminService}
}

// GetAllUsersHandler обрабатывает запрос для получения всех пользователей
func (h *AdminHandler) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("GetAllUsersHandler: Fetching all users")
	users, err := h.AdminService.GetAllUsers()
	if err != nil {
		log.Printf("GetAllUsersHandler: Failed to fetch users: %v\n", err)
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		log.Printf("GetAllUsersHandler: Failed to encode response: %v\n", err)
		http.Error(w, "Failed to encode users", http.StatusInternalServerError)
		return
	}

	log.Println("GetAllUsersHandler: Users fetched successfully")
}

// UpdateUserRoleHandler обрабатывает запрос для обновления роли пользователя
func (h *AdminHandler) UpdateUserRoleHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("UpdateUserRoleHandler: Updating user role")

	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("UpdateUserRoleHandler: Invalid user ID: %v\n", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var data struct {
		Role string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Printf("UpdateUserRoleHandler: Failed to decode request body: %v\n", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("UpdateUserRoleHandler: Attempting to update user ID=%d to role=%s\n", userID, data.Role)

	if err := h.AdminService.UpdateUserRole(uint(userID), data.Role); err != nil {
		log.Printf("UpdateUserRoleHandler: Failed to update user role: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("UpdateUserRoleHandler: User ID=%d role updated successfully\n", userID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User role updated successfully"))
}

// BlockUserHandler обрабатывает запрос для блокировки пользователя
func (h *AdminHandler) BlockUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("BlockUserHandler: Blocking user")

	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("BlockUserHandler: Invalid user ID: %v\n", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	log.Printf("BlockUserHandler: Attempting to block user ID=%d\n", userID)

	if err := h.AdminService.BlockUser(uint(userID)); err != nil {
		log.Printf("BlockUserHandler: Failed to block user: %v\n", err)
		http.Error(w, "Failed to block user", http.StatusInternalServerError)
		return
	}

	log.Printf("BlockUserHandler: User ID=%d blocked successfully\n", userID)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User blocked successfully"))
}
