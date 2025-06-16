package book

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/maithuc2003/re-book-api/internal/models"
	"github.com/maithuc2003/re-book-api/internal/service/book"

	"github.com/go-sql-driver/mysql"
)

type BookHandler struct {
	serviceBook book.BookServiceInterface
}

func NewBookHandler(serviceBook book.BookServiceInterface) *BookHandler {
	return &BookHandler{serviceBook: serviceBook}
}

// Thuộc tính Fontend gửi backend gửi cái gì (intetnet) tcp,http
func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Parse the request body to get the book details
	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest) // JSON gửi lên từ client không hợp
		return
	}

	// Set createdAt hiện tại
	book.CreatedAt = time.Now()
	err := h.serviceBook.CreateBook(&book)

	if err != nil {
		// Log chi tiết lỗi ở server để biết nguyên nhân
		log.Printf("Created book error : %v", err)
		// Nếu là lỗi business logic (validate), trả lỗi chi tiết cho client
		switch err.Error() {
		case "book is nil":
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		case "book title is required":
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		case "book author ID is required":
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		case "book quantity cannot be negative":
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1452 {
			http.Error(w, "Failed to create book: the book_id does not exist.", http.StatusBadRequest)
			return
		}

		http.Error(w, "Failed to create book due to internal server error.", http.StatusInternalServerError)
		//Server gặp lỗi khi tạo sách
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)

}

func (h *BookHandler) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	books, err := h.serviceBook.GetAllBooks()
	if err != nil {
		log.Printf("GetAllBooks error : %v", err)
		if err.Error() == "no books found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get books", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(books)
}

func (h *BookHandler) GetByBookID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// 1. Lấy tham số `id` từ URL query
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing 'id' parameter", http.StatusBadRequest)
		return
	}
	// 2. Chuyển id từ string sang int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid 'id' parameter", http.StatusBadRequest)
		return
	}
	// 3. Gọi service để lấy sách
	book, err := h.serviceBook.GetByBookID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) DeleteById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. Lấy tham số `id` từ URL query
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing 'id' parameter", http.StatusBadRequest)
		return
	}
	// 2. Chuyển id từ string sang int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid 'id' parameter", http.StatusBadRequest)
		return
	}
	// 3.Gọi service để xóa sách
	book, err := h.serviceBook.DeleteById(id)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "invalid book ID"):
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		case strings.Contains(err.Error(), "existing orders"):
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		default:
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) UpdateById(w http.ResponseWriter, r *http.Request) {
	// Bảo vệ: chỉ cho phép PUT
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. Lấy tham số `id` từ URL query
	idStr := r.URL.Query().Get("id")
	// fmt.Println(idStr)
	if idStr == "" {
		http.Error(w, "Missing 'id' parameter", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid 'id' parameter", http.StatusBadRequest)
		return
	}

	// 2. Decode body (JSON) vào struct Book
	var updateBook models.Book
	if err := json.NewDecoder(r.Body).Decode(&updateBook); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
	// 3. Gán lại id cho book để chắc chắn đúng
	updateBook.ID = id // Gán ID từ URL vào struct
	updateBook.UpdatedAt = time.Now()
	// 3.Gọi service để cập nhất sách
	book, err := h.serviceBook.UpdateById(&updateBook)
	if err != nil {
		switch err.Error() {
		case "book is nil", "invalid book ID", "book title is required", "book author ID is required", "book quantity cannot be negative":
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, err.Error(), http.StatusNotFound)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)

}
