// server/service/menu_service.go
package service

import (
	"ozzcafe/server/dal"
	"ozzcafe/server/models"
)

// MenuService представляет бизнес-логику для работы с меню
type MenuService struct {
	MenuDAL *dal.MenuDAL
}

// NewMenuService создает новый объект MenuService
func NewMenuService(menuDAL *dal.MenuDAL) *MenuService {
	return &MenuService{MenuDAL: menuDAL}
}

// GetAllMenuItems возвращает список всех блюд
func (s *MenuService) GetAllMenuItems() ([]models.MenuItem, error) {
	return s.MenuDAL.GetAllMenuItems()
}

// CreateMenuItem добавляет новое блюдо
func (s *MenuService) CreateMenuItem(menuItem *models.MenuItem) error {
	return s.MenuDAL.CreateMenuItem(menuItem)
}

// UpdateMenuItem обновляет информацию о блюде
func (s *MenuService) UpdateMenuItem(menuItem *models.MenuItem) error {
	return s.MenuDAL.UpdateMenuItem(menuItem)
}

// DeleteMenuItem удаляет блюдо из меню
func (s *MenuService) DeleteMenuItem(menuItemID uint) error {
	return s.MenuDAL.DeleteMenuItem(menuItemID)
}
