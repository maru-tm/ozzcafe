// server/dal/menu_dal.go
package dal

import (
	"ozzcafe/server/models"

	"gorm.io/gorm"
)

// MenuDAL структура для работы с меню
type MenuDAL struct {
	DB *gorm.DB
}

// NewMenuDAL создает новый объект MenuDAL
func NewMenuDAL(db *gorm.DB) *MenuDAL {
	return &MenuDAL{DB: db}
}

// GetAllMenuItems возвращает список всех блюд
func (dal *MenuDAL) GetAllMenuItems() ([]models.MenuItem, error) {
	var menuItems []models.MenuItem
	err := dal.DB.Find(&menuItems).Error
	return menuItems, err
}

// CreateMenuItem добавляет новое блюдо в меню
func (dal *MenuDAL) CreateMenuItem(menuItem *models.MenuItem) error {
	return dal.DB.Create(menuItem).Error
}

// UpdateMenuItem обновляет информацию о блюде
func (dal *MenuDAL) UpdateMenuItem(menuItem *models.MenuItem) error {
	return dal.DB.Save(menuItem).Error
}

// DeleteMenuItem удаляет блюдо из меню
func (dal *MenuDAL) DeleteMenuItem(menuItemID uint) error {
	return dal.DB.Delete(&models.MenuItem{}, menuItemID).Error
}
