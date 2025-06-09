package config

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"github.com/joho/godotenv"
)

func init() {
	// Load biến môi trường từ .env, không lỗi nếu không tìm thấy file
	_ = godotenv.Load(".env")
}

func GetDSN() string {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	database := os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", user, password, host, database)
	return dsn
}
