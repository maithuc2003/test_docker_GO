package author

import (
	"database/sql"
	"fmt"
	"github.com/maithuc2003/re-book-api/internal/models"

	"github.com/go-sql-driver/mysql"
)

type authorRepo struct {
	db *sql.DB
}

func NewAuthorRepo(db *sql.DB) AuthorRepositoriesInterface {
	return &authorRepo{db: db}
}

func (r *authorRepo) GetAllAuthors() ([]*models.Author, error) {
	rows, err := r.db.Query("SELECT `id`, `name`, `nationality`, `created_at`, `updated_at` FROM `authors`")
	if err != nil {
		return nil, fmt.Errorf("failed to query author: %w", err)
	}
	defer rows.Close()

	var authors []*models.Author
	for rows.Next() {
		author := &models.Author{}
		err := rows.Scan(&author.ID, &author.Name, &author.Nationality, &author.CreatedAt, &author.UpdatedAt)
		if err != nil {
			return nil, err
		}
		authors = append(authors, author)
	}
	return authors, nil
}

func (r *authorRepo) GetByAuthorID(id int) (*models.Author, error) {
	row := r.db.QueryRow("SELECT `id`, `name`, `nationality`, `created_at`, `updated_at` FROM `authors` WHERE id = ?", id)
	author := &models.Author{}
	err := row.Scan(&author.ID, &author.Name, &author.Nationality, &author.CreatedAt, &author.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("author with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to fetch author: %w", err)
	}
	return author, nil
}

// Implement the BookReader interface
func (r *authorRepo) CreateAuthor(author *models.Author) error {
	query := "INSERT INTO `authors`(`id`, `name`, `nationality`, `created_at`) VALUES (?,?,?,?)"
	result, err := r.db.Exec(query, author.ID, author.Name, author.Nationality, author.CreatedAt)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	author.ID = int(id)
	return nil
}

func (r *authorRepo) DeleteById(id int) (*models.Author, error) {
	author, err := r.GetByAuthorID(id)
	if err != nil {
		return nil, err
	}
	result, err := r.db.Exec("DELETE FROM `authors` WHERE id = ?", id)
	if err != nil {
		// Kiểm tra nếu lỗi là lỗi khóa ngoại (foreign key)
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1451 {
			return nil, fmt.Errorf("cannot delete author: existing orders depend on it")
		}
		return nil, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("no author found with id %d", id)
	}
	return author, nil
}

func (r *authorRepo) UpdateById(author *models.Author) (*models.Author, error) {
	// Kiểm tra author_id có tồn tại không
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM authors WHERE id = ?)", author.ID).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("author_id %d does not exist", author.ID)
	}
	result, err := r.db.Exec(`
			UPDATE authors
			SET name = ?, nationality = ? , updated_at = ?
			WHERE id = ?`,
		author.Name, author.Nationality, author.UpdatedAt, author.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update author: %w", err)
	}
	// Kiểm tra có hàng nào bị ảnh hưởng không
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("no book updated with id %d", author.ID)
	}
	return author, nil
}
