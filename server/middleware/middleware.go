package middleware

// import (
// 	"net/http"
// 	"ozzcafe/server/service"
// 	"strings"

// 	"github.com/dgrijalva/jwt-go"
// )

// func AdminMiddleware(adminService *service.AdminService) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		tokenString := r.Header.Get("Authorization")
// 		if tokenString == "" {
// 			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
// 			return
// 		}

// 		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

// 		claims := &jwt.StandardClaims{}
// 		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
// 			// Получите ваш секретный ключ из конфигурации
// 			return []byte("your_secret_key"), nil
// 		})

// 		if err != nil || !token.Valid {
// 			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
// 			return
// 		}

// 		// Проверяем роль пользователя
// 		userID := claims.Subject // обычно ID пользователя будет в subject
// 		user, err := adminService.GetUserByID(userID)
// 		if err != nil || user.Role != "admin" {
// 			http.Error(w, "Forbidden", http.StatusForbidden)
// 			return
// 		}

// 		// Если все проверки прошли, передаем управление следующему обработчику
// 		next.ServeHTTP(w, r)
// 	}
// }
