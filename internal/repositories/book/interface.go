package book

import "github.com/maithuc2003/re-book-api/internal/models"

// internal/repositories/book/interface.go
type BookRepoInterface interface {
	Create(book *models.Book) error
	GetAllBooks() ([]*models.Book, error)
	GetByBookID(id int) (*models.Book, error)
	DeleteById(id int) (*models.Book, error)
	UpdateById(book *models.Book) (*models.Book, error)
}
