package book

import (
	"github.com/maithuc2003/re-book-api/internal/models"
	repositories "github.com/maithuc2003/re-book-api/internal/repositories/book"
)

type BookService struct {
	repo repositories.BookRepoInterface
}

func NewBookService(repo repositories.BookRepoInterface) *BookService {
	return &BookService{repo: repo}
}

func (s *BookService) CreateBook(book *models.Book) error {
	return s.repo.Create(book)
}

func (s *BookService) GetAllBooks() ([]*models.Book, error) {
	return s.repo.GetAllBooks()
}

func (s *BookService) GetByBookID(id int) (*models.Book, error) {
	return s.repo.GetByBookID(id)
}

func (s *BookService) DeleteById(id int) (*models.Book, error) {
	return s.repo.DeleteById(id)
}

func (s *BookService) UpdateById(book *models.Book) (*models.Book, error) {
	return s.repo.UpdateById(book)
}
