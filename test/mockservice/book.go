package mockservice

import (
	"github.com/maithuc2003/re-book-api/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockBookService struct {
	mock.Mock
}

func (m *MockBookService) CreateBook(book *models.Book) error {
	args := m.Called(book)
	return args.Error(0)
}

func (m *MockBookService) GetAllBooks() ([]*models.Book, error) {
	args := m.Called()
	return args.Get(0).([]*models.Book), args.Error(1)
}

func (m *MockBookService) GetByBookID(id int) (*models.Book, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Book), args.Error(1)
}

func (m *MockBookService) DeleteById(id int) (*models.Book, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Book), args.Error(1)
}

func (m *MockBookService) UpdateById(book *models.Book) (*models.Book, error) {
	args := m.Called(book)
	return args.Get(0).(*models.Book), args.Error(1)
}
