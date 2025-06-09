package repositories

import (
	"github.com/maithuc2003/re-book-api/internal/models"
)

// type OrderReader interface {
// 	GetByOrderID(id int) (*models.Order, error)
// 	GetAllOrders() ([]*models.Order, error)
// }

// type OrderUpdater interface {
// 	UpdateByOrderID(order *models.Order) (*models.Order, error)
// }

// type OrderDeleter interface {
// 	DeleteByOrderID(id int) (*models.Order, error)
// }
// type OrderCreater interface {
// 	Create(order *models.Order) error
// }

type OrderReposiotory interface {
	GetByOrderID(id int) (*models.Order, error)
	GetAllOrders() ([]*models.Order, error)
	UpdateByOrderID(order *models.Order) (*models.Order, error)
	DeleteByOrderID(id int) (*models.Order, error)
	Create(order *models.Order) error

}
