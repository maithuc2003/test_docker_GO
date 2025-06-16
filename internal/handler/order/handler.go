package order

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/maithuc2003/re-book-api/internal/models"
	"github.com/maithuc2003/re-book-api/internal/service/order"

	"github.com/go-sql-driver/mysql"
)

type OrderHandler struct {
	serviceOrder order.OrderServiceInterface
}

func NewOrderHandler(serviceOrder order.OrderServiceInterface) *OrderHandler {
	return &OrderHandler{serviceOrder: serviceOrder}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Parse request body
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	order.OrderedAt = time.Now()

	err := h.serviceOrder.CreateOrder(&order)
	if err != nil {
		// Log lỗi server
		log.Printf("CreateOrder error: %v", err)
		// Validate các lỗi đầu vào từ service
		switch err.Error() {
		case "order is nil", "invalid book ID", "invalid user ID", "quantity must be greater than zero", "status is required":
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Kiểm tra lỗi MySQL foreign key
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1452 {
			http.Error(w, "Failed to create order: foreign key constraint violation.", http.StatusBadRequest)
			return
		}

		// Xử lý lỗi theo nội dung chuỗi error message
		w.Header().Set("Content-Type", "application/json")

		var statusCode int
		var errorMsg string

		switch {
		case strings.Contains(err.Error(), "no rows in result set"):
			statusCode = http.StatusNotFound
			errorMsg = "Product not found or no stock information"
		case strings.Contains(err.Error(), "not enough stock"):
			statusCode = http.StatusBadRequest
			errorMsg = "Not enough stock available"
		default:
			statusCode = http.StatusInternalServerError
			errorMsg = "Internal server error"
		}

		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(map[string]string{
			"error": errorMsg,
		})
		return
	}

	// Nếu thành công trả về JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Order created successfully",
		"order":   order,
	})
}

func (h *OrderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	orders, err := h.serviceOrder.GetAllOrders()
	if err != nil {
		log.Printf("GetAllOrder errr: %v", err)
		if err.Error() == "no books found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get order", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}

func (h *OrderHandler) GetByOrderID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// 1. Lấy tham số id từ url query
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing 'id' parameter", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid 'id' parameter", http.StatusBadRequest)
		return
	}
	// 3. Gọi service để lấy order
	order, err := h.serviceOrder.GetByOrderID(id)
	if err != nil {
		switch err.Error() {
		case "invalid order ID":
			http.Error(w, err.Error(), http.StatusBadRequest)
		case "existing orders": // nếu có xử lý cụ thể, giữ lại
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, err.Error(), http.StatusNotFound)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(order)

}

func (h *OrderHandler) DeleteByOrderID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing 'id' parameter", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid 'id' parameter", http.StatusBadRequest)
		return
	}
	order, err := h.serviceOrder.DeleteByOrderID(id)
	if err != nil {
		if err.Error() == "invalid order ID" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(order)

}

func (h *OrderHandler) UpdateByOrderID(w http.ResponseWriter, r *http.Request) {
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
	// 2. Deconde body (JSON) vào struct BOOK
	var updateOrder models.Order
	if err := json.NewDecoder(r.Body).Decode(&updateOrder); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
	// 3. Gán lại id cho order để đúng
	updateOrder.ID = id
	updateOrder.UpdatedAt = time.Now()
	// 3. Gọi service để cập nhập order
	order, err := h.serviceOrder.UpdateByOrderID(&updateOrder)
	if err != nil {
		switch err.Error() {
		case "order is nil",
			"invalid order ID",
			"invalid book ID",
			"invalid user ID",
			"quantity must be greater than zero",
			"status is required":
			http.Error(w, "Validation error: "+err.Error(), http.StatusBadRequest)
			return
		case "foreign key constraint fails: book_id does not exist":
			http.Error(w, "Invalid book_id: book does not exist", http.StatusBadRequest)
			return
		default:
			http.Error(w, "Failed to update order: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(order)
}
