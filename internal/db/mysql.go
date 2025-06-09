package db

import (
	"database/sql"
	"fmt"
	"github.com/maithuc2003/re-book-api/config"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLConnection struct {
	DB *sql.DB
}

// NewMySQLConnection tạo và trả về kết nối DB mới
func NewMySQLConnection() (*MySQLConnection, error) {
	dsn := config.GetDSN()

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("Failed connect:", err) // nếu có lỗi thì in ra
		return nil, err
	}
	if err := db.Ping(); err != nil {
		fmt.Println("Error connect:", err) // nếu có lỗi thì in ra
		return nil, err
	}
	return &MySQLConnection{DB: db}, nil
}

func (conn *MySQLConnection) Close() error {
	// Đóng kết nối DB nếu nó không phải là nil
	if conn.DB != nil {
		return conn.DB.Close()
	}
	return nil
}
