package models

import (
	"time"

	"github.com/lib/pq"
)

type MenuItem struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `json:"description"`
	Price       float64        `gorm:"not null" json:"price"`
	Ingredients pq.StringArray `gorm:"type:text[]" json:"ingredients"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}
