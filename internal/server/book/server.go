package book

import (
	"database/sql"
	"net/http"
	bookHandler "github.com/maithuc2003/re-book-api/internal/handler/book"
	bookRepo "github.com/maithuc2003/re-book-api/internal/repositories/book"
	bookService "github.com/maithuc2003/re-book-api/internal/service/book"
)

func SetupServerBook(mux *http.ServeMux, db *sql.DB) {
	repo := bookRepo.NewBookRepo(db)
	service := bookService.NewBookService(repo)
	handler := bookHandler.NewBookHandler(service)

	mux.HandleFunc("/book/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler.CreateBook(w, r) // Gọi hàm CreateBook từ handle
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handler.GetAllBooks(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/book", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handler.GetBookByID(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/book/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			handler.DeleteById(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/book/update", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			handler.UpdateById(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	
}
