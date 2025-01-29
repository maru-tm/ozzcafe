package service

import (
	"errors"
	"fmt"
	"log"
	"ozzcafe/server/dal"
	"ozzcafe/server/models"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// UserService структура для бизнес-логики пользователя
type UserService struct {
	UserDAL      *dal.UserDAL
	EmailService *EmailService
}

// NewUserService создает новый объект UserService
func NewUserService(userDAL *dal.UserDAL, emailService *EmailService) *UserService {
	return &UserService{
		UserDAL:      userDAL,
		EmailService: emailService,
	}
}

// RegisterUser регистрирует нового пользователя
func (s *UserService) RegisterUser(name, email, password, confirmPassword string) (*models.User, error) {
	// Логирование данных
	log.Printf("Registering user: %s, Email=%s\n", name, email)

	// Проверка совпадения паролей
	if password != confirmPassword {
		log.Println("Passwords do not match")
		return nil, fmt.Errorf("passwords do not match")
	}

	// Проверка наличия пользователя с таким email
	_, err := s.UserDAL.GetUserByEmail(email)
	if err == nil {
		log.Println("Email already registered:", email)
		return nil, fmt.Errorf("email is already registered")
	}

	// Хэширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		return nil, err
	}

	// Создание нового пользователя
	user := &models.User{
		Name:      name,
		Email:     email,
		Password:  string(hashedPassword),
		Role:      models.RoleUser, // Роль по умолчанию
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Сохранение пользователя в базе данных
	err = s.UserDAL.CreateUser(user)
	if err != nil {
		log.Println("Error creating user:", err)
		return nil, err
	}

	log.Printf("User created successfully: %s, ID=%d\n", user.Name, user.ID)

	// Отправка email для подтверждения
	err = s.EmailService.SendVerificationEmail(user)
	if err != nil {
		log.Println("Error sending verification email:", err)
		return nil, err
	}

	return user, nil
}

// VerifyEmail проверяет токен и подтверждает email пользователя
func (s *UserService) VerifyEmail(token string) error {
	// Декодируйте или проверьте токен
	email, err := decodeVerificationToken(token)
	if err != nil {
		log.Println("Invalid verification token:", err)
		return fmt.Errorf("invalid token")
	}

	// Найдите пользователя по email
	user, err := s.UserDAL.GetUserByEmail(email)
	if err != nil {
		log.Println("User not found for email:", email)
		return fmt.Errorf("user not found")
	}

	// Если email уже подтвержден, вернуть ошибку
	if user.EmailConfirmed {
		log.Println("Email already confirmed:", email)
		return fmt.Errorf("email already confirmed")
	}

	// Обновите статус подтверждения email
	user.EmailConfirmed = true
	err = s.UserDAL.UpdateUser(user)
	if err != nil {
		log.Println("Error updating user email confirmation status:", err)
		return fmt.Errorf("could not confirm email")
	}

	log.Printf("Email confirmed successfully for user: %s\n", email)
	return nil
}

// generateVerificationToken генерирует токен подтверждения
func generateVerificationToken(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // Токен действителен 24 часа
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secret_key"))
	if err != nil {
		log.Println("Error signing token:", err)
		return "", err
	}
	return tokenString, nil
}

func decodeVerificationToken(token string) (string, error) {
	// В реальном приложении используйте библиотеку JWT для декодирования
	decodedEmail := token // Для примера токен и есть email
	return decodedEmail, nil
}

// LoginUser проверяет учетные данные пользователя и возвращает JWT токен
func (s *UserService) LoginUser(email, password string) (string, error) {
	// Получение пользователя из базы данных
	user, err := s.UserDAL.GetUserByEmail(email)
	if err != nil {
		log.Println("User not found:", email)
		return "", errors.New("invalid email or password")
	}

	// Проверка подтверждения email
	if !user.EmailConfirmed {
		log.Println("Email not confirmed:", email)
		return "", errors.New("email is not verified")
	}

	// Проверка пароля
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Println("Invalid password for email:", email)
		return "", errors.New("invalid email or password")
	}

	// Генерация JWT токена
	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(24 * time.Hour).Unix(), // Токен действует 24 часа
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secret_key")) // Используйте секретный ключ из конфигурации
	if err != nil {
		log.Println("Error generating token for user:", email)
		return "", errors.New("failed to generate token")
	}

	log.Printf("User logged in successfully: %s, ID=%d\n", email, user.ID)
	return tokenString, nil
}

// LogoutUser выполняет логику выхода пользователя
func (s *UserService) LogoutUser(token string) error {
	// Для упрощения просто логируем выход
	// В реальном приложении можно добавить токен в "чёрный список" или выполнять дополнительные действия
	log.Println("User logged out. Token invalidated:", token)
	// Пример: сохранять отозванные токены в базе или кэш-системе
	// err := s.TokenDAL.RevokeToken(token)
	// if err != nil {
	//     return fmt.Errorf("failed to revoke token: %w", err)
	// }
	return nil
}

// GetUserByToken извлекает пользователя из базы данных, основываясь на JWT токене
func (s *UserService) GetUserByToken(tokenString string) (*models.User, error) {
	// Если токен передан через URL параметр, то сразу используем его
	if strings.HasPrefix(tokenString, "Bearer ") {
		tokenString = tokenString[7:] // Удаляем "Bearer " из строки
	}

	// Разбираем токен
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверка метода подписи
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Возвращаем ключ для проверки подписи
		return []byte("secret_key"), nil
	})
	if err != nil || !token.Valid {
		log.Println("Invalid token:", err)
		return nil, fmt.Errorf("invalid token")
	}

	// Извлекаем данные из токена
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		log.Println("Invalid token claims")
		return nil, fmt.Errorf("invalid token claims")
	}

	// Получаем userID из токена
	userID, ok := claims["id"].(float64)
	if !ok {
		log.Println("User ID not found in token")
		return nil, fmt.Errorf("user ID not found in token")
	}

	// Получаем пользователя по ID
	user, err := s.UserDAL.GetByID(uint(userID))
	if err != nil {
		log.Println("User not found:", err)
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}
