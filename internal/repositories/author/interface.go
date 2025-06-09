package author

import "github.com/maithuc2003/re-book-api/internal/models"

// type AuthorReader interface {
// 	GetByAuthorID(id int) (*models.Author, error)
// 	GetAllAuthor() ([]*models.Author, error)
// }

// type AuthorCreater interface {
// 	CreateAuthor(author *models.Author) error
// }

// type AuthorUpdater interface {
// 	UpdateById(author *models.Author) (*models.Author, error)
// }

// type AuthorDeleter interface {
// 	DeleteById(id int) (*models.Author, error)
// }

type AuthorRepositories interface {
	GetByAuthorID(id int) (*models.Author, error)
	GetAllAuthor() ([]*models.Author, error)
	CreateAuthor(author *models.Author) error
	UpdateById(author *models.Author) (*models.Author, error)
	DeleteById(id int) (*models.Author, error)
}
