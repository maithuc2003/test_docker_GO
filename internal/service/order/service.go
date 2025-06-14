package order

import (
	"errors"
	"strings"
	"time"

	"github.com/maithuc2003/re-book-api/internal/models"
	repositories "github.com/maithuc2003/re-book-api/internal/repositories/order"
)

type OrderService struct {
	repo repositories.OrderReposiotoryInterface
}

func NewOrderService(repo repositories.OrderReposiotoryInterface) *OrderService {
	return &OrderService{repo: repo}
}

// CreateOrder kiểm tra dữ liệu đầu vào trước khi tạo
func (s *OrderService) CreateOrder(order *models.Order) error {
	if order == nil {
		return errors.New("order is nil")
	}
	if order.BookID <= 0 {
		return errors.New("invalid book ID")
	}
	if order.UserID <= 0 {
		return errors.New("invalid user ID")
	}
	if order.Quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}
	if strings.TrimSpace(order.Status) == "" {
		return errors.New("status is required")
	}

	order.OrderedAt = time.Now()
	order.UpdatedAt = time.Now()

	return s.repo.Create(order)
}

// GetAllOrders kiểm tra lỗi khi lấy danh sách
func (s *OrderService) GetAllOrders() ([]*models.Order, error) {
	orders, err := s.repo.GetAllOrders()
	if err != nil {
		return nil, err
	}
	if len(orders) == 0 {
		return nil, errors.New("no orders found")
	}
	return orders, nil
}

// GetByOrderID kiểm tra ID hợp lệ
func (s *OrderService) GetByOrderID(id int) (*models.Order, error) {
	if id <= 0 {
		return nil, errors.New("invalid order ID")
	}
	return s.repo.GetByOrderID(id)
}

// DeleteByOrderID kiểm tra ID hợp lệ
func (s *OrderService) DeleteByOrderID(id int) (*models.Order, error) {
	if id <= 0 {
		return nil, errors.New("invalid order ID")
	}
	return s.repo.DeleteByOrderID(id)
}

// UpdateByOrderID kiểm tra dữ liệu trước khi cập nhật
func (s *OrderService) UpdateByOrderID(order *models.Order) (*models.Order, error) {
	if order == nil {
		return nil, errors.New("order is nil")
	}
	if order.ID <= 0 {
		return nil, errors.New("invalid order ID")
	}
	if order.BookID <= 0 {
		return nil, errors.New("invalid book ID")
	}
	if order.UserID <= 0 {
		return nil, errors.New("invalid user ID")
	}
	if order.Quantity <= 0 {
		return nil, errors.New("quantity must be greater than zero")
	}
	if strings.TrimSpace(order.Status) == "" {
		return nil, errors.New("status is required")
	}

	order.UpdatedAt = time.Now()

	return s.repo.UpdateByOrderID(order)
}
