package book

import (
	"errors"
	"strings"

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
	if book == nil {
		return errors.New("book is nil")
	}
	if strings.TrimSpace(book.Title) == "" {
		return errors.New("book title is required")
	}
	if book.AuthorID <= 0 {
		return errors.New("book author ID is required")
	}
	if book.Stock < 0 {
		return errors.New("book quantity cannot be negative")
	}

	return s.repo.Create(book)
}

// GetAllBooks trả về lỗi nếu không có sách nào
func (s *BookService) GetAllBooks() ([]*models.Book, error) {
	books, err := s.repo.GetAllBooks()
	if err != nil {
		return nil, err
	}
	if len(books) == 0 {
		return nil, errors.New("no books found")
	}
	return books, nil
}

// GetByBookID kiểm tra ID hợp lệ
func (s *BookService) GetByBookID(id int) (*models.Book, error) {
	if id <= 0 {
		return nil, errors.New("invalid book ID")
	}
	return s.repo.GetByBookID(id)
}

// DeleteById kiểm tra ID hợp lệ
func (s *BookService) DeleteById(id int) (*models.Book, error) {
	if id <= 0 {
		return nil, errors.New("invalid book ID")
	}
	return s.repo.DeleteById(id)
}

// UpdateById kiểm tra dữ liệu trước khi cập nhật
func (s *BookService) UpdateById(book *models.Book) (*models.Book, error) {
	if book == nil {
		return nil, errors.New("book is nil")
	}
	if book.ID <= 0 {
		return nil, errors.New("invalid book ID")
	}
	if strings.TrimSpace(book.Title) == "" {
		return nil, errors.New("book title is required")
	}
	if book.AuthorID <= 0 {
		return nil, errors.New("book author ID is required")
	}
	if book.Stock < 0 {
		return nil, errors.New("book quantity cannot be negative")
	}

	return s.repo.UpdateById(book)
}
