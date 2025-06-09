package book

import (
	"database/sql"
	"net/http"
	authorHandler "github.com/maithuc2003/re-book-api/internal/handler/author"
	authorRepo "github.com/maithuc2003/re-book-api/internal/repositories/author"
	authorService "github.com/maithuc2003/re-book-api/internal/service/author"
)

func SetupServerAuthor(mux *http.ServeMux, db *sql.DB) {
	repo := authorRepo.NewAuthorRepo(db)
	service := authorService.NewBookService(repo)
	handler := authorHandler.NewAuthorHandler(service)
	mux.HandleFunc("/authors", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handler.GetAllAuthors(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/author", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handler.GetByAuthorID(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/author/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler.CreateAuthor(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/author/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			handler.DeleteById(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/author/update", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			handler.UpdateById(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}
