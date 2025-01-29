package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"ozzcafe/server/models"
	"ozzcafe/server/service"

	"github.com/dgrijalva/jwt-go"
)

// UserRegistrationHandler обрабатывает запрос на регистрацию пользователя
func UserRegistrationHandler(userService *service.UserService, emailService *service.EmailService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Name            string `json:"name"`
			Email           string `json:"email"`
			Password        string `json:"password"`
			ConfirmPassword string `json:"confirmPassword"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Регистрация пользователя
		user, err := userService.RegisterUser(req.Name, req.Email, req.Password, req.ConfirmPassword)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Отправка письма с подтверждением через emailService
		err = emailService.SendVerificationEmail(user)
		if err != nil {
			http.Error(w, "Failed to send confirmation email", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		response := map[string]string{
			"message":                   "User registered successfully. Please check your email to verify your account.",
			"user_id":                   fmt.Sprintf("%d", user.ID),
			"requiresEmailVerification": "true",
		}
		json.NewEncoder(w).Encode(response)
	}
}

// UserLoginHandler обрабатывает запрос на авторизацию пользователя
func UserLoginHandler(userService *service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		// Пытаемся выполнить вход
		token, err := userService.LoginUser(req.Email, req.Password)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		// Получаем информацию о пользователе (например, роль)
		user, err := userService.UserDAL.GetUserByEmail(req.Email)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "User not found")
			return
		}

		// Логируем токен и роль
		log.Printf("User logged in successfully, token: %s, Role: %s\n", token, user.Role)

		// Возвращаем успешный ответ с токеном и ролью
		respondWithJSON(w, http.StatusOK, map[string]interface{}{
			"message": "User logged in successfully",
			"token":   token,
			"role":    user.Role,
		})
	}
}
func UserLogoutHandler(userService *service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		// Получаем токен из заголовка Authorization
		token := r.Header.Get("Authorization")
		if token == "" {
			respondWithError(w, http.StatusUnauthorized, "No token provided")
			return
		}
		log.Println("Token received:", token)

		// Вызываем функцию для logout
		err := userService.LogoutUser(token)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Logout failed")
			log.Printf("Error during logout: %v", err) // Дополнительное логирование
			return
		}

		// Успешный logout
		respondWithJSON(w, http.StatusOK, map[string]string{
			"message": "User logged out successfully",
		})
		log.Println("User logged out successfully")
	}
}

// respondWithJSON отправляет ответ в формате JSON
func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Println("Error encoding response:", err)
	}
}

// respondWithError отправляет ответ об ошибке в формате JSON
func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	respondWithJSON(w, statusCode, map[string]string{"error": message})
}

// UserVerificationHandler обрабатывает запрос на подтверждение email
func UserVerificationHandler(userService *service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.URL.Query().Get("token")
		email := r.URL.Query().Get("email")

		// Проверка на наличие токена и email
		if tokenStr == "" || email == "" {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		// Пытаемся проверить токен
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret_key"), nil
		})
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Проверяем данные токена
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Проверка email из токена
		if claims["email"] != email {
			http.Error(w, "Email mismatch", http.StatusBadRequest)
			return
		}

		// Получаем пользователя из базы по email
		user, err := userService.UserDAL.GetUserByEmail(email)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		// Проверка, был ли уже подтвержден email
		if user.EmailConfirmed {
			http.Error(w, "Email already confirmed", http.StatusBadRequest)
			return
		}

		// Подтверждаем пользователя
		user.EmailConfirmed = true
		err = userService.UserDAL.UpdateUser(user)
		if err != nil {
			http.Error(w, "Error updating user", http.StatusInternalServerError)
			return
		}

		// Перенаправление на страницу входа
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

// ResendVerificationEmailHandler повторно отправляет email для подтверждения
func ResendVerificationEmailHandler(emailService *service.EmailService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Email string `json:"email"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Отправляем повторное письмо с подтверждением
		err := emailService.SendVerificationEmail(&models.User{Email: req.Email})
		if err != nil {
			http.Error(w, "Error sending verification email", http.StatusInternalServerError)
			return
		}

		// Ответ с подтверждением
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"message": "Verification email sent successfully!"}`)
	}
}
