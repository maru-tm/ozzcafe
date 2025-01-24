package router

import (
	"net/http"
	"ozzcafe/server/dal"
	"ozzcafe/server/handlers"
	"ozzcafe/server/service"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// NewRouter создает новый маршрутизатор и настраивает маршруты
func NewRouter(db *gorm.DB) *mux.Router {
	// Создаем DAL (Data Access Layer)
	userDal := dal.NewUserDal(db)
	emailService := service.NewEmailService() // Используем конструктор без аргументов

	adminDal := dal.NewUserDal(db)
	adminService := service.NewAdminService(adminDal)

	// Создаем Service, передавая DAL
	userService := service.NewUserService(userDal, emailService)

	// Инициализируем роутер
	r := mux.NewRouter()

	// Регистрируем маршруты для статичных страниц
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "pages/index.html")
	}).Methods("GET")
	r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "pages/signup.html")
	}).Methods("GET")
	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "pages/login.html")
	}).Methods("GET")
	r.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "pages/admin_panel.html")
	}).Methods("GET")

	// Регистрируем маршруты для пользователя
	r.HandleFunc("/register", handlers.UserRegistrationHandler(userService, emailService)).Methods("POST")
	r.HandleFunc("/verify", handlers.UserVerificationHandler(userService)).Methods("GET")
	r.HandleFunc("/login", handlers.UserLoginHandler(userService)).Methods("POST")

	// Регистрируем маршруты для повторной отправки письма с подтверждением
	r.HandleFunc("/resend-verification-email", handlers.ResendVerificationEmailHandler(emailService)).Methods("POST")

	// Инициализируем обработчик администратора
	adminHandler := handlers.NewAdminHandler(adminService)

	// Регистрируем маршруты для пользователей
	r.HandleFunc("/api/admin/users", adminHandler.GetAllUsersHandler).Methods("GET")          // Получить всех пользователей
	r.HandleFunc("/api/admin/users/{id}/block", adminHandler.BlockUserHandler).Methods("PUT") // Заблокировать пользователя
	r.HandleFunc("/api/admin/users/{id}", adminHandler.UpdateUserRoleHandler).Methods("PUT")

	// Инициализируем DAL для работы с меню
	menuDal := dal.NewMenuDAL(db)                  // Создаем DAL для работы с меню
	menuService := service.NewMenuService(menuDal) // Создаем сервис для работы с меню

	// Инициализируем обработчик для меню
	menuHandler := handlers.NewMenuHandler(menuService)

	// Регистрируем маршруты для работы с меню
	r.HandleFunc("/api/admin/menu", menuHandler.CreateMenuItemHandler).Methods("POST")        // Добавить новый элемент меню
	r.HandleFunc("/api/admin/menu/{id}", menuHandler.UpdateMenuItemHandler).Methods("PUT")    // Обновить элемент меню
	r.HandleFunc("/api/admin/menu/{id}", menuHandler.DeleteMenuItemHandler).Methods("DELETE") // Удалить элемент меню
	r.HandleFunc("/api/admin/menu", menuHandler.GetAllMenuItemsHandler).Methods("GET")        // Получить все элементы меню

	return r
}
