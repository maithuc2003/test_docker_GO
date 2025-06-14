package repositories_test

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/maithuc2003/re-book-api/internal/models"
	repositories "github.com/maithuc2003/re-book-api/internal/repositories/order"
	"github.com/stretchr/testify/assert"
)

// fakeBadResult dùng để mô phỏng lỗi khi gọi LastInsertId
type fakeBadResult struct{}

func (f *fakeBadResult) LastInsertId() (int64, error) {
	return 0, errors.New("last insert id error")
}
func (f *fakeBadResult) RowsAffected() (int64, error) {
	return 1, nil
}

func TestOrderRepo_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	repo := repositories.NewOrderRepo(db)
	fakeTime := time.Now()

	tests := []struct {
		name       string
		order      *models.Order
		prepare    func(sqlmock.Sqlmock)
		expectErr  bool
		errMessage string
		checkID    int
	}{
		{
			name:  "Success",
			order: &models.Order{BookID: 1, UserID: 2, Quantity: 3, Status: "pending", OrderedAt: fakeTime},
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT stock FROM books").WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"stock"}).AddRow(10))
				mock.ExpectExec("INSERT INTO orders").
					WithArgs(1, 2, 3, "pending", fakeTime).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE books SET stock").
					WithArgs(3, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			expectErr: false,
			checkID:   1,
		},
		{
			name:  "Insert order error",
			order: &models.Order{BookID: 1, UserID: 2, Quantity: 1, Status: "pending", OrderedAt: fakeTime},
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT stock FROM books").WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"stock"}).AddRow(10))
				mock.ExpectExec("INSERT INTO orders").
					WithArgs(1, 2, 1, "pending", fakeTime).
					WillReturnError(errors.New("insert error")) // 👈 Lỗi tại đây
				mock.ExpectRollback() // 👈 rollback khi lỗi
			},
			expectErr:  true,
			errMessage: "failed to create order",
		},
		{
			name:  "Rollback error",
			order: &models.Order{BookID: 1, UserID: 2, Quantity: 1, Status: "pending", OrderedAt: fakeTime},
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT stock FROM books").WithArgs(1).
					WillReturnError(errors.New("query error"))

				mock.ExpectRollback().WillReturnError(errors.New("rollback failed")) // 👈 lỗi rollback
			},
			expectErr:  true,
			errMessage: "rollback failed", // vì err trả về thẳng luôn
		},
		{
			name:  "Not enough stock",
			order: &models.Order{BookID: 1, UserID: 2, Quantity: 5, Status: "pending", OrderedAt: fakeTime},
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT stock FROM books").WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"stock"}).AddRow(2))
				mock.ExpectRollback()
			},
			expectErr:  true,
			errMessage: "not enough stock available",
		},
		{
			name:  "QueryRow error",
			order: &models.Order{BookID: 999, UserID: 2, Quantity: 1, Status: "pending", OrderedAt: fakeTime},
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT stock FROM books").WithArgs(999).
					WillReturnError(errors.New("query error"))
				mock.ExpectRollback()
			},
			expectErr:  true,
			errMessage: "failed to fetch current stock",
		},
		{
			name:  "Update stock error",
			order: &models.Order{BookID: 1, UserID: 2, Quantity: 2, Status: "pending", OrderedAt: fakeTime},
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT stock FROM books").WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"stock"}).AddRow(10))
				mock.ExpectExec("INSERT INTO orders").
					WithArgs(1, 2, 2, "pending", fakeTime).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE books SET stock").
					WithArgs(2, 1).
					WillReturnError(errors.New("update error"))
				mock.ExpectRollback()
			},
			expectErr:  true,
			errMessage: "failed to update book stock",
		},
		{
			name:  "Commit error",
			order: &models.Order{BookID: 1, UserID: 2, Quantity: 1, Status: "pending", OrderedAt: fakeTime},
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT stock FROM books").WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"stock"}).AddRow(10))
				mock.ExpectExec("INSERT INTO orders").
					WithArgs(1, 2, 1, "pending", fakeTime).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE books SET stock").
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit().WillReturnError(errors.New("commit error"))
			},
			expectErr:  true,
			errMessage: "failed to commit transaction",
		},
		{
			name:  "LastInsertId error",
			order: &models.Order{BookID: 1, UserID: 2, Quantity: 1, Status: "pending", OrderedAt: fakeTime},
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT stock FROM books").WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"stock"}).AddRow(10))

				mock.ExpectExec("INSERT INTO orders").
					WithArgs(1, 2, 1, "pending", fakeTime).
					WillReturnResult(&fakeBadResult{}) // 👈 dùng struct giả ở đây

				mock.ExpectExec("UPDATE books SET stock").
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()

			},
			expectErr:  true,
			errMessage: "failed to retrieve inserted order ID",
		},
		{
			name:  "Begin transaction error",
			order: &models.Order{BookID: 1, UserID: 1, Quantity: 1, Status: "pending", OrderedAt: time.Now()},
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(errors.New("begin failed"))
			},
			expectErr:  true,
			errMessage: "begin failed",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepare(mock)
			err := repo.Create(tc.order)

			if tc.expectErr {
				assert.Error(t, err)
				if tc.errMessage != "" {
					assert.Contains(t, err.Error(), tc.errMessage)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.checkID, tc.order.ID)
			}
		})
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

func TestOrderRepo_GetAllOrders(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewOrderRepo(db)
	fakeTime := time.Now()

	tests := []struct {
		name         string
		prepareMock  func(sqlmock.Sqlmock)
		expectedLen  int
		expectedErr  error
		assertErrMsg string
	}{
		{
			name: "Success",
			prepareMock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "book_id", "user_id", "quantity", "status", "ordered_at", "updated_at",
				}).AddRow(1, 101, 201, 2, "pending", fakeTime, fakeTime).
					AddRow(2, 102, 202, 1, "completed", fakeTime, fakeTime)

				m.ExpectQuery("SELECT \\* FROM `orders`").
					WillReturnRows(rows)
			},
			expectedLen: 2,
			expectedErr: nil,
		},
		{
			name: "Query error",
			prepareMock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT \\* FROM `orders`").
					WillReturnError(errors.New("query error"))
			},
			expectedLen:  0,
			expectedErr:  errors.New("query error"),
			assertErrMsg: "failed to query orders",
		},
		{
			name: "Scan error",
			prepareMock: func(m sqlmock.Sqlmock) {
				// thiếu 1 cột intentionally → lỗi scan
				rows := sqlmock.NewRows([]string{
					"id", "book_id", "user_id", "quantity", "status", "ordered_at",
				}).AddRow(1, 101, 201, 2, "pending", fakeTime)

				m.ExpectQuery("SELECT \\* FROM `orders`").
					WillReturnRows(rows)
			},
			expectedLen:  0,
			expectedErr:  errors.New("scan error"),
			assertErrMsg: "", // error không custom, chỉ assert.Error
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepareMock(mock)

			orders, err := repo.GetAllOrders()

			if tc.expectedErr != nil {
				assert.Error(t, err)
				if tc.assertErrMsg != "" {
					assert.Contains(t, err.Error(), tc.assertErrMsg)
				}
				assert.Nil(t, orders)
			} else {
				assert.NoError(t, err)
				assert.Len(t, orders, tc.expectedLen)
			}
		})
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestOrderRepo_GetByOrderIDƠ(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewOrderRepo(db)
	fakeTime := time.Now()

	tests := []struct {
		name        string
		orderID     int
		prepareMock func(sqlmock.Sqlmock)
		expectErr   bool
		errMessage  string
		expected    *models.Order
	}{
		{
			name:    "Success",
			orderID: 1,
			prepareMock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT `id`, `book_id`, `user_id`, `quantity`, `status`,`ordered_at`, `updated_at` FROM `orders` WHERE id = ?").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "book_id", "user_id", "quantity", "status", "ordered_at", "updated_at",
					}).AddRow(1, 101, 201, 3, "pending", fakeTime, fakeTime))
			},
			expectErr: false,
			expected: &models.Order{
				ID:        1,
				BookID:    101,
				UserID:    201,
				Quantity:  3,
				Status:    "pending",
				OrderedAt: fakeTime,
				UpdatedAt: fakeTime,
			},
		},
		{
			name:    "Not found",
			orderID: 2,
			prepareMock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT `id`, `book_id`, `user_id`, `quantity`, `status`,`ordered_at`, `updated_at` FROM `orders` WHERE id = ?").
					WithArgs(2).
					WillReturnError(sql.ErrNoRows)
			},
			expectErr:  true,
			errMessage: "not found",
			expected:   nil,
		},
		{
			name:    "DB error",
			orderID: 3,
			prepareMock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT `id`, `book_id`, `user_id`, `quantity`, `status`,`ordered_at`, `updated_at` FROM `orders` WHERE id = ?").
					WithArgs(3).
					WillReturnError(errors.New("some db error"))
			},
			expectErr:  true,
			errMessage: "failed to fetch book",
			expected:   nil,
		},
		{
			name:    "Scan error",
			orderID: 4,
			prepareMock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT `id`, `book_id`, `user_id`, `quantity`, `status`,`ordered_at`, `updated_at` FROM `orders` WHERE id = ?").
					WithArgs(4).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "book_id", "user_id", // thiếu cột intentionally
					}).AddRow(4, 104, 204))
			},
			expectErr:  true,
			errMessage: "", // raw error nên không cần contains
			expected:   nil,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepareMock(mock)
			result, err := repo.GetByOrderID(tc.orderID)

			if tc.expectErr {
				assert.Error(t, err)
				if tc.errMessage != "" {
					assert.Contains(t, err.Error(), tc.errMessage)
				}
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestOrderRepo_DeleteByOrderID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewOrderRepo(db)
	fakeTime := time.Now()
	tests := []struct {
		name        string
		orderID     int
		prepareMock func(sqlmock.Sqlmock)
		expectErr   bool
		errContains string
		expected    *models.Order
	}{
		{
			name:    "Success",
			orderID: 1,
			prepareMock: func(m sqlmock.Sqlmock) {
				// Mock GetByOrderID
				m.ExpectQuery("SELECT .* FROM `orders` WHERE id = ?").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "book_id", "user_id", "quantity", "status", "ordered_at", "updated_at",
					}).AddRow(1, 101, 201, 3, "pending", fakeTime, fakeTime))

				// Mock DELETE
				m.ExpectExec("DELETE FROM `orders` WHERE id = ?").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected
			},
			expectErr: false,
			expected: &models.Order{
				ID:        1,
				BookID:    101,
				UserID:    201,
				Quantity:  3,
				Status:    "pending",
				OrderedAt: fakeTime,
				UpdatedAt: fakeTime,
			},
		},
		{
			name:    "GetByOrderID returns error",
			orderID: 2,
			prepareMock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT .* FROM `orders` WHERE id = ?").
					WithArgs(2).
					WillReturnError(sql.ErrNoRows)
			},
			expectErr:   true,
			errContains: "not found",
			expected:    nil,
		},
		{
			name:    "Delete query fails",
			orderID: 3,
			prepareMock: func(m sqlmock.Sqlmock) {
				// Mock GetByOrderID success
				m.ExpectQuery("SELECT .* FROM `orders` WHERE id = ?").
					WithArgs(3).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "book_id", "user_id", "quantity", "status", "ordered_at", "updated_at",
					}).AddRow(3, 103, 203, 2, "pending", fakeTime, fakeTime))

				// Mock DELETE fail
				m.ExpectExec("DELETE FROM `orders` WHERE id = ?").
					WithArgs(3).
					WillReturnError(errors.New("delete failed"))
			},
			expectErr:   true,
			errContains: "failed to delete order",
			expected:    nil,
		},
		{
			name:    "RowsAffected returns error",
			orderID: 6,
			prepareMock: func(m sqlmock.Sqlmock) {
				// Giả lập GetByOrderID thành công
				m.ExpectQuery("SELECT .* FROM `orders` WHERE id = ?").
					WithArgs(6).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "book_id", "user_id", "quantity", "status", "ordered_at", "updated_at",
					}).AddRow(6, 106, 206, 1, "pending", fakeTime, fakeTime))

				// Giả lập DELETE trả về đối tượng .RowsAffected() lỗi
				m.ExpectExec("DELETE FROM `orders` WHERE id = ?").
					WithArgs(6).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("rows affected error")))
			},
			expectErr:   true,
			errContains: "rows affected error",
			expected:    nil,
		},

		{
			name:    "No rows affected",
			orderID: 4,
			prepareMock: func(m sqlmock.Sqlmock) {
				// Mock GetByOrderID success
				m.ExpectQuery("SELECT .* FROM `orders` WHERE id = ?").
					WithArgs(4).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "book_id", "user_id", "quantity", "status", "ordered_at", "updated_at",
					}).AddRow(4, 104, 204, 1, "completed", fakeTime, fakeTime))

				// Mock DELETE returns 0 rows affected
				m.ExpectExec("DELETE FROM `orders` WHERE id = ?").
					WithArgs(4).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectErr:   true,
			errContains: "no book found with id",
			expected:    nil,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepareMock(mock)
			result, err := repo.DeleteByOrderID(tc.orderID)

			if tc.expectErr {
				assert.Error(t, err)
				if tc.errContains != "" {
					assert.Contains(t, err.Error(), tc.errContains)
				}
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestOrderRepo_UpdateByOrderID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewOrderRepo(db)
	fakeTime := time.Now()
	sampleOrder := &models.Order{
		ID:        1,
		BookID:    101,
		UserID:    201,
		Quantity:  2,
		Status:    "completed",
		UpdatedAt: fakeTime,
	}
	tests := []struct {
		name        string
		order       *models.Order
		prepareMock func(sqlmock.Sqlmock)
		expected    *models.Order
		expectErr   bool
		errContains string
	}{
		{
			name:  "Success",
			order: sampleOrder,
			prepareMock: func(m sqlmock.Sqlmock) {
				m.ExpectExec("UPDATE orders").
					WithArgs(
						sampleOrder.BookID,
						sampleOrder.UserID,
						sampleOrder.Quantity,
						sampleOrder.Status,
						sampleOrder.UpdatedAt,
						sampleOrder.ID,
					).
					WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected
			},
			expected:    sampleOrder,
			expectErr:   false,
			errContains: "",
		},
		{
			name:  "Foreign key violation",
			order: sampleOrder,
			prepareMock: func(m sqlmock.Sqlmock) {
				mysqlErr := &mysql.MySQLError{
					Number:  1452,
					Message: "Cannot add or update a child row: a foreign key constraint fails",
				}
				m.ExpectExec("UPDATE orders").
					WithArgs(
						sampleOrder.BookID,
						sampleOrder.UserID,
						sampleOrder.Quantity,
						sampleOrder.Status,
						sampleOrder.UpdatedAt,
						sampleOrder.ID,
					).
					WillReturnError(mysqlErr)
			},
			expected:    nil,
			expectErr:   true,
			errContains: "foreign key constraint fails",
		},
		{
			name:  "Generic DB error",
			order: sampleOrder,
			prepareMock: func(m sqlmock.Sqlmock) {
				m.ExpectExec("UPDATE orders").
					WithArgs(
						sampleOrder.BookID,
						sampleOrder.UserID,
						sampleOrder.Quantity,
						sampleOrder.Status,
						sampleOrder.UpdatedAt,
						sampleOrder.ID,
					).
					WillReturnError(errors.New("db error"))
			},
			expected:    nil,
			expectErr:   true,
			errContains: "db error",
		},
		{
			name:  "RowsAffected error",
			order: sampleOrder,
			prepareMock: func(m sqlmock.Sqlmock) {
				m.ExpectExec("UPDATE orders").
					WithArgs(
						sampleOrder.BookID,
						sampleOrder.UserID,
						sampleOrder.Quantity,
						sampleOrder.Status,
						sampleOrder.UpdatedAt,
						sampleOrder.ID,
					).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("rows affected error")))
			},
			expected:    nil,
			expectErr:   true,
			errContains: "rows affected error",
		},
		{
			name:  "No row updated",
			order: sampleOrder,
			prepareMock: func(m sqlmock.Sqlmock) {
				m.ExpectExec("UPDATE orders").
					WithArgs(
						sampleOrder.BookID,
						sampleOrder.UserID,
						sampleOrder.Quantity,
						sampleOrder.Status,
						sampleOrder.UpdatedAt,
						sampleOrder.ID,
					).
					WillReturnResult(sqlmock.NewResult(0, 0)) // no rows affected
			},
			expected:    nil,
			expectErr:   true,
			errContains: "no order upadted", // typo giữ nguyên nếu code bạn dùng từ này
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepareMock(mock)

			result, err := repo.UpdateByOrderID(tc.order)
			if tc.expectErr {
				assert.Error(t, err)
				if tc.errContains != "" {
					assert.Contains(t, err.Error(), tc.errContains)
				}
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
	assert.NoError(t, mock.ExpectationsWereMet())
}
