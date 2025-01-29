package middleware

import (
	"context"
	"log"
	"net/http"
	"ozzcafe/server/models"
	"ozzcafe/server/service"
	"strings"
)

type contextKey string

const userKey contextKey = "user"

// AuthMiddleware аутентифицирует пользователя, извлекая токен либо из заголовка Authorization,
// либо из параметра URL "token".
func AuthMiddleware(userService *service.UserService) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Проверка токена в параметре URL
			token := r.URL.Query().Get("token")
			if token == "" || !strings.HasPrefix(token, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.") {
				log.Println("Invalid or missing token in URL")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Получаем пользователя по токену
			user, err := userService.GetUserByToken(token)
			if err != nil {
				log.Printf("Error retrieving user by token: %v", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			log.Printf("User authenticated: %v", user)

			// Добавляем пользователя в контекст запроса
			ctx := context.WithValue(r.Context(), userKey, user)
			r = r.WithContext(ctx)

			// Передаем управление следующему обработчику
			next(w, r)
		}
	}
}

// AdminAuthMiddleware аутентифицирует администратора, проверяя роль пользователя.
func AdminAuthMiddleware(userService *service.UserService) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Проверка токена в параметре URL
			token := r.URL.Query().Get("token")
			if token == "" || !strings.HasPrefix(token, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.") {
				log.Println("Invalid or missing token in URL")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Получаем пользователя по токену
			user, err := userService.GetUserByToken(token)
			if err != nil {
				log.Printf("Error retrieving user by token: %v", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Проверка роли пользователя на "admin"
			if user.Role != "admin" {
				log.Printf("Access denied: User %v does not have admin privileges", user.ID)
				http.Error(w, "Forbidden: You do not have access", http.StatusForbidden)
				return
			}

			log.Printf("Admin user authenticated: %v", user)

			// Добавляем пользователя в контекст запроса
			ctx := context.WithValue(r.Context(), userKey, user)
			r = r.WithContext(ctx)

			// Передаем управление следующему обработчику
			next(w, r)
		}
	}
}

// GetUserFromContext извлекает пользователя из контекста запроса.
func GetUserFromContext(r *http.Request) *models.User {
	user, _ := r.Context().Value(userKey).(*models.User)
	return user
}
