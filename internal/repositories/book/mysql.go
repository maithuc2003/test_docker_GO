package book

import (
	"database/sql"
	"fmt"
	"github.com/maithuc2003/re-book-api/internal/models"

	"github.com/go-sql-driver/mysql"
)

type bookRepo struct {
	db *sql.DB
}

func NewBookRepo(db *sql.DB) *bookRepo {
	return &bookRepo{db: db}
}

// Implement the BookReader interface
func (r *bookRepo) Create(book *models.Book) error {
	query := "INSERT INTO `books`(`id`, `title`, `author_id`, `stock`, `created_at`) VALUES (?,?,?,?,?)"
	result, err := r.db.Exec(query, book.ID, book.Title, book.AuthorID, book.Stock, book.CreatedAt)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	book.ID = int(id)
	return nil
}

// Implement interface method
func (r *bookRepo) GetAllBooks() ([]*models.Book, error) {
	// rows, err := r.db.Query("SELECT id, title, author_id, stock, created_at, updated_at FROM books")
	if err != nil {
		return nil, fmt.Errorf("failed to query books: %w", err)
	}
	defer rows.Close()

	var books []*models.Book
	for rows.Next() {
		book := &models.Book{}
		// err := rows.Scan(&book.ID, &book.Title, &book.AuthorID, &book.Stock, &book.CreatedAt, &book.UpdatedAt)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

func (r *bookRepo) GetByBookID(id int) (*models.Book, error) {
	row := r.db.QueryRow("SELECT id, title, author_id, stock, created_at, updated_at FROM books WHERE id = ?", id)
	book := &models.Book{}
	err := row.Scan(&book.ID, &book.Title, &book.AuthorID, &book.Stock, &book.CreatedAt, &book.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("book with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to fetch book: %w", err)
	}
	return book, nil
}

func (r *bookRepo) DeleteById(id int) (*models.Book, error) {
	book, err := r.GetByBookID(id)
	if err != nil {
		return nil, err
	}
	result, err := r.db.Exec("DELETE FROM `books` WHERE id = ?", id)
	if err != nil {
		// Kiểm tra nếu lỗi là lỗi khóa ngoại (foreign key)
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1451 {
			return nil, fmt.Errorf("cannot delete book: existing orders depend on it")
		}
		return nil, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("no book found with id %d", id)
	}
	return book, nil
}

func (r *bookRepo) UpdateById(book *models.Book) (*models.Book, error) {
	// Kiểm tra author_id có tồn tại không
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM authors WHERE id = ?)", book.AuthorID).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("author_id %d does not exist", book.AuthorID)
	}
	result, err := r.db.Exec(`
			UPDATE books
			SET title = ?, author_id = ?, stock = ? , updated_at = ?
			WHERE id = ?`,
		book.Title, book.AuthorID, book.Stock, book.UpdatedAt, book.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update book: %w", err)
	}
	// Kiểm tra có hàng nào bị ảnh hưởng không
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("no book updated with id %d", book.ID)
	}
	return book, nil
}

