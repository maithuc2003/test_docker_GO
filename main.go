package main

import (
	"fmt"
	"github.com/maithuc2003/re-book-api/internal/db"
	"log"
	"net/http"
	"os"
	server_author "github.com/maithuc2003/re-book-api/internal/server/author"
	server_book "github.com/maithuc2003/re-book-api/internal/server/book"
	server_order "github.com/maithuc2003/re-book-api/internal/server/order"
)

func main() {
	conn, err := db.NewMySQLConnection() // nhận biến conn và err
	if err != nil {
		// xử lý lỗi, ví dụ in ra và thoát
		fmt.Println("Failed to connect:", err)
		return
	}
	// OK do
	defer conn.Close() // gọi đóng kết nối khi main kết thúc

	// Route api
	mux := http.NewServeMux()
	server_book.SetupServerBook(mux, conn.DB)
	server_order.SetupOrderServer(mux, conn.DB)
	server_author.SetupServerAuthor(mux, conn.DB)

	// Port
	log.Println("Server started at", os.Getenv("PORT"))
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), mux); err != nil {
		log.Fatal(err)
	}
}
