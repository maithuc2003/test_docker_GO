package book_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/maithuc2003/re-book-api/internal/handler/book"
	"github.com/maithuc2003/re-book-api/internal/models"
	"github.com/maithuc2003/re-book-api/test/mockservice"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllBooks(t *testing.T) {
	tests := []struct {
		name             string
		httpMethod       string
		mockReturn       []*models.Book
		mockError        error
		expectedStatus   int
		expectedResult   []*models.Book
		expectedErrorMsg string
	}{
		{
			name:       "Success - return list of books",
			httpMethod: http.MethodGet,
			mockReturn: []*models.Book{
				{ID: 1, Title: "Go Programming", AuthorID: 2, CreatedAt: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC), UpdatedAt: time.Date(2024, time.January, 2, 0, 0, 0, 0, time.UTC)},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedResult: []*models.Book{
				{ID: 1, Title: "Go Programming", AuthorID: 2, CreatedAt: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC), UpdatedAt: time.Date(2024, time.January, 2, 0, 0, 0, 0, time.UTC)},
			},
		},
		{
			name:             "Error from service - internal server error",
			httpMethod:       http.MethodGet,
			mockReturn:       nil,
			mockError:        errors.New("database error"),
			expectedStatus:   http.StatusInternalServerError,
			expectedErrorMsg: "Failed to get books",
		},
		{
			name:             "No books found - 404 error",
			httpMethod:       http.MethodGet,
			mockReturn:       nil,
			mockError:        errors.New("no books found"),
			expectedStatus:   http.StatusNotFound,
			expectedErrorMsg: "no books found",
		},
		{
			name:             "Invalid HTTP method - Method Not Allowed",
			httpMethod:       http.MethodPost,
			expectedStatus:   http.StatusMethodNotAllowed,
			expectedErrorMsg: "Method not allowed",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mock_service := new(mockservice.MockBookService)
			handler := book.NewBookHandler(mock_service)

			if tc.httpMethod == http.MethodGet {
				mock_service.On("GetAllBooks").Return(tc.mockReturn, tc.mockError)
			}
			req := httptest.NewRequest(tc.httpMethod, "/books", nil)
			w := httptest.NewRecorder()

			handler.GetAllBooks(w, req)
			assert.Equal(t, tc.expectedStatus, w.Code)

			if tc.expectedStatus == http.StatusOK && tc.mockError == nil {
				var result []*models.Book
				err := json.NewDecoder(w.Body).Decode(&result)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			} else if tc.expectedErrorMsg != "" {
				assert.Contains(t, w.Body.String(), tc.expectedErrorMsg)
			}
			mock_service.AssertExpectations(t)
		})
	}
}

func TestCreateBook(t *testing.T) {
	Tests := []struct {
		name           string
		httpMethod     string
		requestBody    interface{}
		mockError      error
		expectedStatus int
		expectErrorMsg string
	}{
		{
			name:       "Success - valid book creation",
			httpMethod: http.MethodPost,
			requestBody: &models.Book{
				Title:    "Clean Code",
				Stock:    100,
				AuthorID: 1,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid JSON body",
			httpMethod:     http.MethodPost,
			requestBody:    "invalid-json",
			expectedStatus: http.StatusBadRequest,
			expectErrorMsg: "Invalid request body",
		},
		{
			name:           "Invalid method - GET not allowed",
			httpMethod:     http.MethodGet,
			expectedStatus: http.StatusMethodNotAllowed,
			expectErrorMsg: "Method not allowed",
		},
		{
			name:       "MySQL foreign key error",
			httpMethod: http.MethodPost,
			requestBody: &models.Book{
				Title:    "Domain Driven Design",
				Stock:    50,
				AuthorID: 999,
			},
			mockError: &mysql.MySQLError{
				Number:  1452,
				Message: "Cannot add or update a child row: a foreign key constraint fails",
			},
			expectedStatus: http.StatusBadRequest,
			expectErrorMsg: "book_id does not exist",
		},
		{
			name:       "Internal server error",
			httpMethod: http.MethodPost,
			requestBody: &models.Book{
				Title:    "System Failure",
				Stock:    20,
				AuthorID: 1,
			},
			mockError:      errors.New("db error"),
			expectedStatus: http.StatusInternalServerError,
			expectErrorMsg: "internal server error",
		},
		{
			name:           "Book is nil",
			httpMethod:     http.MethodPost,
			requestBody:    &models.Book{}, // Gửi rỗng để mô phỏng book bị nil trong service
			mockError:      errors.New("book is nil"),
			expectedStatus: http.StatusBadRequest,
			expectErrorMsg: "book is nil",
		},
		{
			name:       "Missing title",
			httpMethod: http.MethodPost,
			requestBody: &models.Book{
				Stock:    10,
				AuthorID: 1,
			},
			mockError:      errors.New("book title is required"),
			expectedStatus: http.StatusBadRequest,
			expectErrorMsg: "book title is required",
		},
		{
			name:       "Missing author ID",
			httpMethod: http.MethodPost,
			requestBody: &models.Book{
				Title: "No Author",
				Stock: 5,
			},
			mockError:      errors.New("book author ID is required"),
			expectedStatus: http.StatusBadRequest,
			expectErrorMsg: "book author ID is required",
		},
		{
			name:       "Negative stock",
			httpMethod: http.MethodPost,
			requestBody: &models.Book{
				Title:    "Negative Stock",
				Stock:    -5,
				AuthorID: 1,
			},
			mockError:      errors.New("book quantity cannot be negative"),
			expectedStatus: http.StatusBadRequest,
			expectErrorMsg: "book quantity cannot be negative",
		},
	}
	for _, tc := range Tests {
		t.Run(tc.name, func(t *testing.T) {
			mock_service := new(mockservice.MockBookService)
			handler := book.NewBookHandler(mock_service)

			var bodyBytes []byte
			var err error

			if s, ok := tc.requestBody.(string); ok && s == "invalid-json" {
				bodyBytes = []byte(`invalid-json`)
			} else {
				bodyBytes, err = json.Marshal(tc.requestBody)
				assert.NoError(t, err)
			}
			req := httptest.NewRequest(tc.httpMethod, "/book/add", bytes.NewReader(bodyBytes))
			w := httptest.NewRecorder()

			if tc.httpMethod != http.MethodPost {
				handler.CreateBook(w, req)
				assert.Equal(t, tc.expectedStatus, w.Code)
				assert.Contains(t, w.Body.String(), tc.expectErrorMsg)
				return
			}

			if tc.expectErrorMsg != "Invalid request body" {
				mock_service.On("CreateBook", mock.AnythingOfType("*models.Book")).Return(tc.mockError)
			}
			handler.CreateBook(w, req)
			assert.Equal(t, tc.expectedStatus, w.Code)
			if tc.expectErrorMsg != "" {
				assert.Contains(t, w.Body.String(), tc.expectErrorMsg)
			} else {
				var result models.Book
				err := json.NewDecoder(w.Body).Decode(&result)
				assert.NoError(t, err)
				expected := tc.requestBody.(*models.Book)
				assert.Equal(t, expected.Title, result.Title)
				assert.Equal(t, expected.Stock, result.Stock)
				assert.Equal(t, expected.AuthorID, result.AuthorID)
			}
			mock_service.AssertExpectations(t)
		})
	}

}

func TestDeleteById(t *testing.T) {
	tests := []struct {
		name           string
		httpMethod     string // ✅ thêm field
		queryParam     string
		mockReturn     *models.Book
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Missing ID",
			httpMethod:     http.MethodDelete,
			queryParam:     "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Missing 'id' parameter",
		},
		{
			name:           "Invalid Book ID in service",
			httpMethod:     http.MethodDelete,
			queryParam:     "id=0", // id <= 0 sẽ gây ra lỗi
			mockReturn:     nil,
			mockError:      errors.New("invalid book ID"),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid book ID",
		},

		{
			name:           "Invalid ID",
			httpMethod:     http.MethodDelete,
			queryParam:     "id=abc",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid 'id' parameter",
		},
		{
			name:           "Book not found",
			httpMethod:     http.MethodDelete,
			queryParam:     "id=2",
			mockReturn:     nil,
			mockError:      errors.New("book not found"),
			expectedStatus: http.StatusNotFound,
			expectedBody:   "book not found",
		},
		{
			name:           "Book has existing orders",
			httpMethod:     http.MethodDelete,
			queryParam:     "id=3",
			mockReturn:     nil,
			mockError:      errors.New("book has existing orders"),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "book has existing orders",
		},
		{
			name:       "Success delete book",
			httpMethod: http.MethodDelete,
			queryParam: "id=1",
			mockReturn: &models.Book{
				ID:       1,
				Title:    "Clean Code",
				Stock:    10,
				AuthorID: 1,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `"title":"Clean Code"`,
		},
		{
			name:           "Invalid HTTP method",
			httpMethod:     http.MethodGet, // ✅ sai method để test
			queryParam:     "id=1",
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "Method not allowed",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mock_service := new(mockservice.MockBookService)
			handler := book.NewBookHandler(mock_service)

			url := "/book/delete"
			if tc.queryParam != "" {
				url += "?" + tc.queryParam
			}

			req := httptest.NewRequest(tc.httpMethod, url, nil)
			w := httptest.NewRecorder()

			// Chỉ mock khi đúng method DELETE và có dữ liệu hợp lệ
			if tc.httpMethod == http.MethodDelete && (tc.mockReturn != nil || tc.mockError != nil) {
				if idStr := req.URL.Query().Get("id"); idStr != "" {
					id, err := strconv.Atoi(idStr)
					if err == nil {
						mock_service.On("DeleteById", id).Return(tc.mockReturn, tc.mockError)
					}
				}
			}

			handler.DeleteById(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tc.expectedBody)
			mock_service.AssertExpectations(t)
		})
	}
}

func TestUpdateByID(t *testing.T) {
	tests := []struct {
		name           string
		queryParam     string
		requestBody    string
		mockReturn     *models.Book
		mockError      error
		expectedStatus int
		expectedBody   string
		httpMethod     string // đổi tên từ 'method' -> 'httpMethod'
	}{
		{
			name:           "Missing ID",
			queryParam:     "",
			requestBody:    `{"title":"Clean Code","stock":10,"author_id":1}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Missing 'id' parameter",
			httpMethod:     http.MethodPut,
		},

		{
			name:           "Invalid ID",
			queryParam:     "id=abc",
			requestBody:    `{"title":"Clean Code","stock":10,"author_id":1}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid 'id' parameter",
			httpMethod:     http.MethodPut,
		},
		{
			name:           "Invalid JSON body",
			queryParam:     "id=1",
			requestBody:    `invalid_json`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid JSON body",
			httpMethod:     http.MethodPut,
		},
		{
			name:           "Book not found",
			queryParam:     "id=2",
			requestBody:    `{"title":"New Title","stock":5,"author_id":1}`,
			mockReturn:     nil,
			mockError:      errors.New("book not found"),
			expectedStatus: http.StatusNotFound,
			expectedBody:   "book not found",
			httpMethod:     http.MethodPut,
		},
		{
			name:        "Successful update",
			queryParam:  "id=1",
			requestBody: `{"title":"Updated Book","stock":20,"author_id":1}`,
			mockReturn: &models.Book{
				ID:       1,
				Title:    "Updated Book",
				Stock:    20,
				AuthorID: 1,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `"title":"Updated Book"`,
			httpMethod:     http.MethodPut,
		},
		{
			name:           "Invalid Method",
			queryParam:     "id=1",
			requestBody:    `{"title":"Should Fail","stock":5,"author_id":1}`,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "Method not allowed",
			httpMethod:     http.MethodGet,
		},
		{
			name:           "Book is nil",
			queryParam:     "id=1",
			requestBody:    `{"title":"Updated Book","stock":5,"author_id":1}`, // vẫn cần requestBody để decode thành book
			mockReturn:     nil,
			mockError:      errors.New("book is nil"),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "book is nil",
			httpMethod:     http.MethodPut,
		},
		{
			name:           "Invalid book ID",
			queryParam:     "id=-1",
			requestBody:    `{"title":"Updated Book","stock":5,"author_id":1}`,
			mockReturn:     nil,
			mockError:      errors.New("invalid book ID"),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid book ID",
			httpMethod:     http.MethodPut,
		},
		{
			name:           "Missing title",
			queryParam:     "id=1",
			requestBody:    `{"title":" ","stock":5,"author_id":1}`,
			mockReturn:     nil,
			mockError:      errors.New("book title is required"),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "book title is required",
			httpMethod:     http.MethodPut,
		},
		{
			name:           "Invalid author ID",
			queryParam:     "id=1",
			requestBody:    `{"title":"Book","stock":5,"author_id":0}`,
			mockReturn:     nil,
			mockError:      errors.New("book author ID is required"),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "book author ID is required",
			httpMethod:     http.MethodPut,
		},
		{
			name:           "Negative stock",
			queryParam:     "id=1",
			requestBody:    `{"title":"Book","stock":-5,"author_id":1}`,
			mockReturn:     nil,
			mockError:      errors.New("book quantity cannot be negative"),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "book quantity cannot be negative",
			httpMethod:     http.MethodPut,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mock_service := new(mockservice.MockBookService)
			handler := book.NewBookHandler(mock_service)

			url := "/book/update"
			if tc.queryParam != "" {
				url += "?" + tc.queryParam
			}

			req := httptest.NewRequest(tc.httpMethod, url, strings.NewReader(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			// Nếu cần mock
			if tc.mockReturn != nil || tc.mockError != nil {
				if idStr := req.URL.Query().Get("id"); idStr != "" {
					if id, err := strconv.Atoi(idStr); err == nil {
						mock_service.On("UpdateById", mock.MatchedBy(func(a *models.Book) bool {
							return a.ID == id
						})).Return(tc.mockReturn, tc.mockError)
					}
				}
			}
			handler.UpdateById(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tc.expectedBody)
			mock_service.AssertExpectations(t)
		})
	}
}

func TestGetBookByID(t *testing.T) {
	tests := []struct {
		name           string
		httpMethod     string
		queryParam     string
		mockReturn     *models.Book
		mockError      error
		expectedStatus int
		expectedResult *models.Book
		expectErrorMsg string
	}{
		{
			name:           "Missing ID",
			httpMethod:     http.MethodGet,
			queryParam:     "",
			expectedStatus: http.StatusBadRequest,
			expectErrorMsg: "Missing 'id' parameter",
		},
		{
			name:           "Invalid ID format",
			httpMethod:     http.MethodGet,
			queryParam:     "id=abc",
			expectedStatus: http.StatusBadRequest,
			expectErrorMsg: "Invalid 'id' parameter",
		},
		{
			name:           "Book not found",
			httpMethod:     http.MethodGet,
			queryParam:     "id=1",
			mockReturn:     nil,
			mockError:      errors.New("book not found"),
			expectedStatus: http.StatusNotFound,
			expectErrorMsg: "book not found",
		},
		{
			name:       "Successful Get",
			httpMethod: http.MethodGet,
			queryParam: "id=2",
			mockReturn: &models.Book{
				ID:       2,
				Title:    "Golang 101",
				Stock:    5,
				AuthorID: 1,
			},
			expectedStatus: http.StatusOK,
			expectedResult: &models.Book{
				ID:       2,
				Title:    "Golang 101",
				Stock:    5,
				AuthorID: 1,
			},
		},
		{
			name:           "Wrong HTTP method",
			httpMethod:     http.MethodPost,
			queryParam:     "id=1",
			expectedStatus: http.StatusMethodNotAllowed,
			expectErrorMsg: "Method not allowed",
		},
	}

	for _, tc := range tests {
		mock_service := new(mockservice.MockBookService)
		handler := book.NewBookHandler(mock_service)

		url := "/book"
		if tc.queryParam != "" {
			url += "?" + tc.queryParam
		}
		req := httptest.NewRequest(tc.httpMethod, url, nil)
		w := httptest.NewRecorder()

		if tc.httpMethod == http.MethodGet && tc.queryParam != "" && (tc.mockError != nil || tc.mockReturn != nil) {
			mock_service.On("GetByBookID", mock.AnythingOfType("int")).Return(tc.mockReturn, tc.mockError)
		}

		handler.GetByBookID(w, req)
		assert.Equal(t, tc.expectedStatus, w.Code)
		if tc.expectedStatus == http.StatusOK && tc.mockReturn != nil {
			var result models.Book
			err := json.NewDecoder(w.Body).Decode(&result)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedResult, &result)
		} else if tc.expectErrorMsg != "" {
			assert.Contains(t, w.Body.String(), tc.expectErrorMsg)
		}
		mock_service.AssertExpectations(t)
	}
}
