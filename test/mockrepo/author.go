package mockrepo

import (
	"github.com/maithuc2003/re-book-api/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockAuthorRepository struct {
	mock.Mock
}

func (m *MockAuthorRepository) GetAllAuthors() ([]*models.Author, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]*models.Author), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAuthorRepository) GetByAuthorID(id int) (*models.Author, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Author), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAuthorRepository) CreateAuthor(author *models.Author) error {
	args := m.Called(author)
	return args.Error(0)
}

func (m *MockAuthorRepository) UpdateById(author *models.Author) (*models.Author, error) {
	args := m.Called(author)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Author), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAuthorRepository) DeleteById(id int) (*models.Author, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Author), args.Error(1)
	}
	return nil, args.Error(1)
}
