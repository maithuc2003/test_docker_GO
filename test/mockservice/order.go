package  mockservice

import (
	"github.com/stretchr/testify/mock"
	"github.com/maithuc2003/re-book-api/internal/models"
)

// MockOrderService mocks the OrderServiceInterface
type MockOrderService struct {
	mock.Mock
}

func (m *MockOrderService) CreateOrder(order *models.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockOrderService) GetAllOrders() ([]*models.Order, error) {
	args := m.Called()
	return args.Get(0).([]*models.Order), args.Error(1)
}

func (m *MockOrderService) GetByOrderID(id int) (*models.Order, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderService) DeleteByOrderID(id int) (*models.Order, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderService) UpdateByOrderID(order *models.Order) (*models.Order, error) {
	args := m.Called(order)
	return args.Get(0).(*models.Order), args.Error(1)
}
