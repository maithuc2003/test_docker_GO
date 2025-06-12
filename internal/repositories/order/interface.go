package repositories

import (
	"github.com/maithuc2003/re-book-api/internal/models"
)



type OrderReposiotoryInterface interface {
	GetByOrderID(id int) (*models.Order, error)
	GetAllOrders() ([]*models.Order, error)
	UpdateByOrderID(order *models.Order) (*models.Order, error)
	DeleteByOrderID(id int) (*models.Order, error)
	Create(order *models.Order) error

}
