package author

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/maithuc2003/re-book-api/internal/models"
	"github.com/maithuc2003/re-book-api/internal/service/author"
)

type AuthorHandler struct {
	serviceAuthor author.AuthorServiceInterface
}

func NewAuthorHandler(serviceAuthor author.AuthorServiceInterface) *AuthorHandler {
	return &AuthorHandler{serviceAuthor: serviceAuthor}
}

func (h *AuthorHandler) GetAllAuthors(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	authors, err := h.serviceAuthor.GetAllAuthors()
	if err != nil {
		log.Printf("GetAllAuthors error : %v", err)
		if err.Error() == "no authors found in the system" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get authors", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(authors)
}

func (h *AuthorHandler) GetByAuthorID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. Get the 'id' parameter from query
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing 'id' parameter", http.StatusBadRequest)
		return
	}

	// 2. Convert 'id' to integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid 'id' parameter", http.StatusBadRequest)
		return
	}

	// 3. Call service to fetch author
	author, err := h.serviceAuthor.GetByAuthorID(id)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "invalid author ID"):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case strings.Contains(err.Error(), "author not found"):
			http.Error(w, err.Error(), http.StatusNotFound)
		case strings.Contains(err.Error(), "failed to retrieve author"):
			http.Error(w, err.Error(), http.StatusInternalServerError)
		default:
			http.Error(w, "Unexpected error", http.StatusInternalServerError)
		}
		return
	}

	// 4. Return author in JSON format
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(author)
}

// Thuộc tính Fontend gửi backend gửi cái gì (intetnet) tcp,http
func (h *AuthorHandler) CreateAuthor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Parse the request body to get the book details
	var author models.Author
	if err := json.NewDecoder(r.Body).Decode(&author); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest) // JSON gửi lên từ client không hợp
		return
	}

	// Set createdAt hiện tại
	author.CreatedAt = time.Now()
	err := h.serviceAuthor.CreateAuthor(&author)

	if err != nil {
		// Log chi tiết lỗi ở server để biết nguyên nhân
		log.Printf("Created author error : %v", err)
		// error user (invalid input)
		if strings.Contains(err.Error(), "name cannot be empty") ||
			strings.Contains(err.Error(), "already exists") ||
			strings.Contains(err.Error(), "author is nil") {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		//error from Mysql
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1452 {
			http.Error(w, "Failed to create author: the author_id does not exist.", http.StatusBadRequest)
			return
		}

		http.Error(w, "Failed to create author due to internal server error.", http.StatusInternalServerError)
		//Server gặp lỗi khi tạo sách
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(author)

}

func (h *AuthorHandler) DeleteById(w http.ResponseWriter, r *http.Request) {
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
	author, err := h.serviceAuthor.DeleteById(id)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "invalid author ID"):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case strings.Contains(err.Error(), "author not found"):
			http.Error(w, err.Error(), http.StatusNotFound)
		case strings.Contains(err.Error(), "existing author with books"):
			http.Error(w, err.Error(), http.StatusBadRequest) 
		case strings.Contains(err.Error(), "failed to delete author"):
			http.Error(w, err.Error(), http.StatusInternalServerError)
		default:
			http.Error(w, "Unexpected error", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(author)
}

func (h *AuthorHandler) UpdateById(w http.ResponseWriter, r *http.Request) {
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

	// 2. Decode body (JSON) vào struct Author
	var updateAuthor models.Author
	if err := json.NewDecoder(r.Body).Decode(&updateAuthor); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
	// 3. Gán lại id cho book để chắc chắn đúng
	updateAuthor.ID = id // Gán ID từ URL vào struct
	updateAuthor.UpdatedAt = time.Now()
	// 3.Gọi service để cập nhất sách
	author, err := h.serviceAuthor.UpdateById(&updateAuthor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(author)

}
