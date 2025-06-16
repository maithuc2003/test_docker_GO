package author

import (
	"errors"
	"fmt"
	"strings"

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
	if author == nil {
		return errors.New("author is nil")
	}
	if strings.TrimSpace(author.Name) == "" {
		return errors.New("author name cannot be empty")
	}
	existingAuthors, err := s.repo.GetAllAuthors()

	if err != nil {
		return fmt.Errorf("failed to fetch authors for validation: %v", err)
	}
	for _, existing := range existingAuthors {
		if strings.EqualFold(existing.Name, author.Name) {
			return errors.New("author with the same name already exists")
		}
	}
	err = s.repo.CreateAuthor(author)
	if err != nil {
		return fmt.Errorf("failed to create author: %v", err)
	}
	return nil
}
func (s *AuthorService) GetAllAuthors() ([]*models.Author, error) {
	authors, err := s.repo.GetAllAuthors()
	if err != nil {
		return nil, err
	}
	if len(authors) == 0 {
		return nil, errors.New("no authors found in the system")
	}
	return s.repo.GetAllAuthors()
}

func (s *AuthorService) GetByAuthorID(id int) (*models.Author, error) {
	if id <= 0 {
		return nil, errors.New("invalid author ID")
	}
	author, err := s.repo.GetByAuthorID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve author: %v", err)
	}
	if author == nil {
		return nil, errors.New("author not found")
	}
	return author, nil
}

func (s *AuthorService) DeleteById(id int) (*models.Author, error) {
	if id <= 0 {
		return nil, errors.New("invalid author ID")
	}

	deletedAuthor, err := s.repo.DeleteById(id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete author: %v", err)
	}
	if deletedAuthor == nil {
		return nil, errors.New("author not found or already deleted")
	}
	return deletedAuthor, nil
}

func (s *AuthorService) UpdateById(author *models.Author) (*models.Author, error) {
	if author == nil {
		return nil, errors.New("author is nil")
	}
	//Validate the author ID
	if author.ID <= 0 {
		return nil, errors.New("invalid author ID")
	}
	//Ensure the author's name is not empty or just whitespace
	if strings.TrimSpace(author.Name) == "" {
		return nil, errors.New("author name cannot be empty")
	}
	// Check if the author with the given ID actually exists
	existring, err := s.repo.GetByAuthorID(author.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch existing author: %v", err)
	}
	if existring == nil {
		return nil, errors.New("author not found")
	}
	//Ensure the new same does not conflict with any other author's name
	authors, err := s.repo.GetAllAuthors()
	if err != nil {
		return nil, fmt.Errorf("failed to validate author name: %v", err)
	}

	for _, a := range authors {
		//Allow the current author to keep their name, but prevent duplicate
		if a.ID != author.ID && strings.EqualFold(a.Name, author.Name) {
			return nil, errors.New("another author with the same name already exists")
		}
	}

	// Attempt to update the author in the repository
	updateAuthor, err := s.repo.UpdateById(author)
	if err != nil {
		return nil, fmt.Errorf("failed to update author : %v", err)
	}
	return updateAuthor, nil
}
