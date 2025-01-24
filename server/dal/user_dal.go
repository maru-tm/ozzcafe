// server/dal/user_dal.go
package dal

import (
	"log"
	"ozzcafe/server/models"

	"gorm.io/gorm"
)

// UserDAL структура для работы с данными пользователя
type UserDAL struct {
	DB *gorm.DB
}

// NewUserDal создает новый объект UserDAL
func NewUserDal(db *gorm.DB) *UserDAL {
	return &UserDAL{DB: db}
}

// CreateUser добавляет нового пользователя в базу данных
func (dal *UserDAL) CreateUser(user *models.User) error {
	log.Printf("Inserting new user into database: %s\n", user.Email)
	return dal.DB.Create(user).Error
}

// GetUserByEmail ищет пользователя по email
func (dal *UserDAL) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := dal.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		log.Println("User not found:", email)
	} else {
		log.Printf("User found: %s, ID=%d\n", user.Email, user.ID)
	}
	return &user, err
}

// GetByID ищет пользователя по ID
func (dal *UserDAL) GetByID(userID uint) (*models.User, error) {
	var user models.User
	err := dal.DB.Where("id = ?", userID).First(&user).Error
	if err != nil {
		log.Println("User not found:", userID)
	} else {
		log.Printf("User found: %d, %s\n", user.ID, user.Email)
	}
	return &user, err
}

// UpdateUser обновляет информацию о пользователе
func (dal *UserDAL) UpdateUser(user *models.User) error {
	log.Printf("Updating user: %s, ID=%d\n", user.Email, user.ID)
	return dal.DB.Save(user).Error
}

// Save сохраняет изменения пользователя (аналогично UpdateUser)
func (dal *UserDAL) Save(user *models.User) error {
	log.Printf("Saving user: %s, ID=%d\n", user.Email, user.ID)
	return dal.DB.Save(user).Error
}

// GetAllUsers возвращает список всех пользователей
func (dal *UserDAL) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := dal.DB.Find(&users).Error
	if err != nil {
		log.Println("Error fetching all users:", err)
	}
	return users, err
}

// UpdateUserRole обновляет роль пользователя
func (dal *UserDAL) UpdateUserRole(userID uint, role string) error {
	log.Printf("Updating role for user ID=%d to %s\n", userID, role)
	return dal.DB.Model(&models.User{}).Where("id = ?", userID).Update("role", role).Error
}

// BlockUser блокирует пользователя (например, удаляет доступ)
func (dal *UserDAL) BlockUser(userID uint) error {
	log.Printf("Blocking user ID=%d\n", userID)
	return dal.DB.Model(&models.User{}).Where("id = ?", userID).Update("email_confirmed", false).Error
}
