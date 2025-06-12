package author

import (
	"github.com/maithuc2003/re-book-api/internal/models"
	repositories "github.com/maithuc2003/re-book-api/internal/repositories/author"
)

type AuthorService struct {
	repo repositories.AuthorRepositoriesInterface
}

func NewAuthorService(repo repositories.AuthorRepositoriesInterface) *AuthorService {
	return &AuthorService{repo: repo}
}

func (s *AuthorService) CreateAuthor(author *models.Author) error {
	return s.repo.CreateAuthor(author)
}
func (s *AuthorService) GetAllAuthors() ([]*models.Author, error) {
	return s.repo.GetAllAuthors()
}

func (s *AuthorService) GetByAuthorID(id int) (*models.Author, error) {
	return s.repo.GetByAuthorID(id)
}

func (s *AuthorService) DeleteById(id int) (*models.Author, error) {
	return s.repo.DeleteById(id)
}

func (s *AuthorService) UpdateById(author *models.Author) (*models.Author, error) {
	return s.repo.UpdateById(author)
}
