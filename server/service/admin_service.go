package service

import (
	"fmt"
	"log"
	"ozzcafe/server/dal"
	"ozzcafe/server/models"
)

// AdminService представляет бизнес-логику для администратора
type AdminService struct {
	UserDAL *dal.UserDAL
}

// NewAdminService создает новый объект AdminService
func NewAdminService(userDAL *dal.UserDAL) *AdminService {
	return &AdminService{UserDAL: userDAL}
}

// GetAllUsers возвращает список всех пользователей
func (s *AdminService) GetAllUsers() ([]models.User, error) {
	log.Println("AdminService: Fetching all users")
	users, err := s.UserDAL.GetAllUsers()
	if err != nil {
		log.Printf("AdminService: Failed to fetch all users: %v\n", err)
	}
	log.Printf("AdminService: Fetched %d users\n", len(users))
	return users, err
}

// UpdateUserRole обновляет роль пользователя
func (s *AdminService) UpdateUserRole(userID uint, role string) error {
	log.Printf("AdminService: Attempting to update user ID=%d role to %s\n", userID, role)

	// Проверка допустимых значений роли
	var newRole models.Role
	switch role {
	case "user":
		newRole = models.RoleUser
	case "manager":
		newRole = models.RoleManager
	case "admin":
		newRole = models.RoleAdmin
	default:
		log.Printf("AdminService: Invalid role provided: %s\n", role)
		return fmt.Errorf("invalid role")
	}

	// Получение пользователя из базы данных
	user, err := s.UserDAL.GetByID(userID)
	if err != nil {
		log.Printf("AdminService: User ID=%d not found\n", userID)
		return fmt.Errorf("user not found")
	}

	// Обновление роли пользователя
	user.Role = newRole
	if err := s.UserDAL.Save(user); err != nil {
		log.Printf("AdminService: Failed to update user role for user ID=%d: %v\n", userID, err)
		return fmt.Errorf("failed to update user role: %v", err)
	}

	log.Printf("AdminService: User ID=%d role updated successfully to %s\n", userID, role)
	return nil
}

// BlockUser блокирует пользователя
func (s *AdminService) BlockUser(userID uint) error {
	log.Printf("AdminService: Attempting to block user ID=%d\n", userID)
	err := s.UserDAL.BlockUser(userID)
	if err != nil {
		log.Printf("AdminService: Failed to block user ID=%d: %v\n", userID, err)
		return fmt.Errorf("failed to block user")
	}

	log.Printf("AdminService: User ID=%d blocked successfully\n", userID)
	return nil
}
