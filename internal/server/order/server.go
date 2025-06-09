package order

import (
	"database/sql"
	"net/http"

	orderHandler "github.com/maithuc2003/re-book-api/internal/handler/order"
	orderRepo "github.com/maithuc2003/re-book-api/internal/repositories/order"
	orderService "github.com/maithuc2003/re-book-api/internal/service/order"
)

func SetupOrderServer(mux *http.ServeMux, db *sql.DB) {
	// Khởi tạo các tầng
	repo := orderRepo.NewOrderRepo(db)
	service := orderService.NewOrderService(repo)
	handler := orderHandler.NewOrderHandler(service)

	mux.HandleFunc("/order/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler.CreateOrder(w, r) // Gọi hàm CreateOrder từ handler
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handler.GetAllOrders(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/order", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handler.GetByOrderID(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/order/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			handler.DeleteByOrderID(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/order/update", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			handler.UpdateByOrderID(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

}
