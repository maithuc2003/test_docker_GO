package book

import "github.com/maithuc2003/re-book-api/internal/models"

// type BookReader interface {
// 	GetByBookID(id int) (*models.Book, error)
// 	GetAllBooks() ([]*models.Book, error)
// }

// type BookCreater interface {
// 	Create(book *models.Book) error
// }

// type BookUpdater interface {
// 	UpdateById(book *models.Book) (*models.Book, error)
// }

// type BookDeleter interface {
// 	DeleteById(id int) (*models.Book, error)
// }

// type BookRepositories interface {
// 	BookCreater
// 	BookReader
// 	BookDeleter
// 	BookUpdater
// }

type BookRepositories interface {
	GetByBookID(id int) (*models.Book, error)
	GetAllBooks() ([]*models.Book, error)
	Create(book *models.Book) error
	DeleteById(id int) (*models.Book, error)
	UpdateById(book *models.Book) (*models.Book, error)
}
