package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/maithuc2003/re-book-api/internal/models"

	"github.com/go-sql-driver/mysql"
)

type orderRepo struct {
	db *sql.DB
}

func NewOrderRepo(db *sql.DB) *orderRepo {
	return &orderRepo{db: db}
}

// Implement the OrderReader interface
func (r *orderRepo) Create(order *models.Order) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	// Step 1: Check current stock
	var currentStock int
	err = tx.QueryRow("SELECT stock FROM books WHERE id = ? FOR UPDATE;", order.BookID).Scan(&currentStock)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to fetch current stock: %w", err)
	}
	if currentStock < order.Quantity {
		tx.Rollback()
		return fmt.Errorf("not enough stock available")
	}
	// step 2: Insert order
	query := "INSERT INTO orders (book_id, user_id, quantity, status ,ordered_at) VALUES (?, ?, ?, ?,?)"
	result, err := tx.Exec(query, order.BookID, order.UserID, order.Quantity, order.Status, order.OrderedAt)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create order: %w", err)
	}
	// time.Sleep(10 * time.Second)
	//step 3 : update book stock
	_, err = tx.Exec(
		"UPDATE books SET stock = stock - ? WHERE id = ?",
		order.Quantity, order.BookID,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update book stock: %w", err)
	}
	// Step 4: Commit transaction
	err = tx.Commit()

	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	// Get the inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to retrieve inserted order ID: %w", err)
	}
	order.ID = int(id) // Set the ID of the order after insertion
	return nil
}

// Implement interface method
func (r *orderRepo) GetAllOrders() ([]*models.Order, error) {
	rows, err := r.db.Query("SELECT * FROM `orders`")
	if err != nil {
		return nil, fmt.Errorf("failed to query orders: %w", err)
	}
	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		order := &models.Order{}
		err := rows.Scan(&order.ID, &order.BookID, &order.UserID, &order.Quantity, &order.Status, &order.OrderedAt, &order.UpdatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (r *orderRepo) GetByOrderID(id int) (*models.Order, error) {
	row := r.db.QueryRow("SELECT `id`, `book_id`, `user_id`, `quantity`, `status`,`ordered_at`, `updated_at` FROM `orders` WHERE id = ?", id)
	order := &models.Order{}
	err := row.Scan(&order.ID, &order.BookID, &order.UserID, &order.Quantity, &order.Status, &order.OrderedAt, &order.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to fetch book: %w", err)
	}
	return order, nil
}

func (r *orderRepo) DeleteByOrderID(id int) (*models.Order, error) {
	order, err := r.GetByOrderID(id)
	if err != nil {
		return nil, err
	}
	result, err := r.db.Exec("DELETE FROM `orders` WHERE id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete order: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("no book found with id %d", id)
	}
	return order, nil
}

func (r *orderRepo) UpdateByOrderID(order *models.Order) (*models.Order, error) {
	result, err := r.db.Exec(`
		UPDATE orders 
		SET book_id = ?, user_id = ?, quantity = ?, status = ?, updated_at = ?
		WHERE id = ?`,
		order.BookID, order.UserID, order.Quantity, order.Status, order.UpdatedAt, order.ID)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1452 {
				// foreign key violation
				return nil, errors.New("foreign key constraint fails: book_id does not exist")
			}
		}
		return nil, err
	}
	// Kiểm tra có hàng nào bị ảnh hưởng không
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("no order upadted with id %d", order.ID)
	}
	return order, nil
}
