package order

import (
	"github.com/maithuc2003/re-book-api/internal/models"
	repositories "github.com/maithuc2003/re-book-api/internal/repositories/order"
)

type OrderService struct {
	repo repositories.OrderReposiotory
}

func NewOrderService(repo repositories.OrderReposiotory) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrder(order *models.Order) error {
	return s.repo.Create(order)
}

func (s *OrderService) GetAllOrders() ([]*models.Order, error) {
	return s.repo.GetAllOrders()
}

func (s *OrderService) GetByOrderID(id int) (*models.Order, error) {
	return s.repo.GetByOrderID(id)
}

func (s *OrderService) DeleteByOrderID(id int) (*models.Order, error) {
	return s.repo.DeleteByOrderID(id)
}

func (s *OrderService) UpdateByOrderID(order *models.Order) (*models.Order, error) {
	return s.repo.UpdateByOrderID(order)
}
