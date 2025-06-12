package mockservice

import (
	"github.com/maithuc2003/re-book-api/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockAuthorService struct {
	mock.Mock
}

func (m *MockAuthorService) GetAllAuthors() ([]*models.Author, error) {
	args := m.Called()
	return args.Get(0).([]*models.Author), args.Error(1)
}


func (m *MockAuthorService) GetByAuthorID(id int) (*models.Author, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Author), args.Error(1)
}

func (m *MockAuthorService) CreateAuthor(author *models.Author) error {
	args := m.Called(author)
	return args.Error(0)
}

func (m *MockAuthorService) DeleteById(id int) (*models.Author, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Author), args.Error(1)
}

func (m *MockAuthorService) UpdateById(author *models.Author) (*models.Author, error) {
	args := m.Called(author)
	return args.Get(0).(*models.Author), args.Error(1)
}