package models

import "time"

// UserRequest используется для валидации входных данных при регистрации пользователя
type UserRequest struct {
	Name            string `json:"name" validate:"required,min=3,max=100"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=Password"`
}

type Role string

const (
	RoleAdmin   Role = "admin"
	RoleManager Role = "manager"
	RoleUser    Role = "user"
)

// User представляет модель пользователя в базе данных
type User struct {
	ID             uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name           string    `json:"name" gorm:"not null;size:100"`
	Email          string    `json:"email" gorm:"unique;not null;size:100"`
	Password       string    `json:"-" gorm:"not null;size:255"`                  // Пароль должен быть захеширован
	Role           Role      `json:"role" gorm:"type:varchar(20);default:'user'"` // Changed to varchar to avoid enum issue
	EmailConfirmed bool      `json:"email_confirmed" gorm:"default:false"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
