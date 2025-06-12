package author_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/maithuc2003/re-book-api/internal/handler/author"
	"github.com/maithuc2003/re-book-api/internal/models"
	mock "github.com/maithuc2003/re-book-api/test/mockservice"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
)

func TestGetAllAuthors(t *testing.T) {
	tests := []struct {
		name           string
		httpMethod     string
		mockReturn     []*models.Author
		mockError      error
		expectedStatus int
		expectedResult []*models.Author
		expectErrorMsg string
	}{
		{
			name:           "success with authors",
			httpMethod:     http.MethodGet,
			mockReturn:     []*models.Author{{ID: 1, Name: "Author A"}, {ID: 2, Name: "Author B"}},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedResult: []*models.Author{{ID: 1, Name: "Author A"}, {ID: 2, Name: "Author B"}},
		},
		{
			name:           "method not allowed",
			httpMethod:     http.MethodPost,
			expectedStatus: http.StatusMethodNotAllowed,
			expectErrorMsg: "Method not allowed",
		},
		{
			name:           "Empty author list",
			httpMethod:     http.MethodGet,
			mockReturn:     []*models.Author{},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedResult: []*models.Author{},
		},
		{
			name:           "Database error from service",
			httpMethod:     http.MethodGet,
			mockReturn:     nil,
			mockError:      errors.New("DB error"),
			expectedStatus: http.StatusInternalServerError,
			expectErrorMsg: "Failed to get authors",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Tạo một mock service để thay thế service thật
			mock_service := new(mock.MockAuthorService)
			// Tạo một handler, truyền mock service vào
			handler := author.NewAuthorHandler(mock_service)
			// Định nghĩa hành vi giả của mock:
			// Khi gọi GetAllAuthor thì trả về kết quả mock và lỗi mock tương ứng
			if tc.httpMethod == http.MethodGet {
				mock_service.On("GetAllAuthors").Return(tc.mockReturn, tc.mockError)
			}
			// Tạo HTTP request giả (GET /authors) và response recorder
			req := httptest.NewRequest(tc.httpMethod, "/authors", nil)
			w := httptest.NewRecorder()
			// Gọi handler để xử lý request
			handler.GetAllAuthors(w, req)
			assert.Equal(t, tc.expectedStatus, w.Code)
			if tc.expectedStatus == http.StatusOK && tc.mockError == nil {
				var result []*models.Author
				err := json.NewDecoder(w.Body).Decode(&result)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			} else if tc.expectErrorMsg != "" {
				assert.Contains(t, w.Body.String(), tc.expectErrorMsg)
			}
			// Đảm bảo rằng tất cả hàm mock đã được gọi như mong đợi
			mock_service.AssertExpectations(t)
		})
	}
}

func TestCreateAuthor(t *testing.T) {
	tests := []struct {
		name           string
		httpMethod     string
		requestBody    interface{}
		mockError      error
		expectedStatus int
		expectErrorMsg string
	}{
		{
			name:           "Success",
			httpMethod:     http.MethodPost,
			requestBody:    &models.Author{Name: "Haruki Murakami", Nationality: "Japanese"},
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
			name:           "MySQL foreign key error",
			httpMethod:     http.MethodPost,
			requestBody:    &models.Author{Name: "Jane Doe"},
			mockError:      &mysql.MySQLError{Number: 1452, Message: "foreign key constraint fails"},
			expectedStatus: http.StatusBadRequest,
			expectErrorMsg: "author_id does not exist",
		},
		{
			name:           "Wrong HTTP method",
			httpMethod:     http.MethodGet,
			requestBody:    nil,
			expectedStatus: http.StatusMethodNotAllowed,
			expectErrorMsg: "method not allowed",
		},
		{
			name:           "Internal server error",
			httpMethod:     http.MethodPost,
			requestBody:    &models.Author{Name: "Jane Doe"},
			mockError:      errors.New("DB connection failed"),
			expectedStatus: http.StatusInternalServerError,
			expectErrorMsg: "internal server error",
		},
	}
	
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockService := new(mock.MockAuthorService)
			handler := author.NewAuthorHandler(mockService)

			var bodyBytes []byte
			var err error

			if s, ok := tc.requestBody.(string); ok && s == "invalid-json" {
				bodyBytes = []byte(`{invalid-json}`)
			} else {
				bodyBytes, err = json.Marshal(tc.requestBody)
				assert.NoError(t, err)
			}
			req := httptest.NewRequest(tc.httpMethod, "/author/add", bytes.NewReader(bodyBytes))
			w := httptest.NewRecorder()

			if tc.httpMethod == http.MethodPost && tc.expectErrorMsg != "Invalid request body" {
				mockService.On("CreateAuthor", testifymock.AnythingOfType("*models.Author")).Return(tc.mockError)
			}

			handler.CreateAuthor(w, req)
			assert.Equal(t, tc.expectedStatus, w.Code)
			if tc.expectErrorMsg != "" {
				assert.Contains(t, w.Body.String(), tc.expectErrorMsg)
			} else {
				var result models.Author
				err := json.NewDecoder(w.Body).Decode(&result)
				assert.NoError(t, err)
				expected := tc.requestBody.(*models.Author)
				assert.Equal(t, expected.Name, result.Name)
			}
			mockService.AssertExpectations(t)
		})
	}
}

func TestGetByAuthorID(t *testing.T) {
	tests := []struct {
		name           string
		httpMethod     string
		queryParam     string
		mockReturn     *models.Author
		mockError      error
		expectedStatus int
		expectedResult *models.Author
		expectErrorMsg string
	}{
		{
			name:           "Success get author by ID",
			httpMethod:     http.MethodGet,
			queryParam:     "id=1",
			mockReturn:     &models.Author{ID: 1, Name: "Author One"},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedResult: &models.Author{ID: 1, Name: "Author One"},
		},
		{
			name:           "Missing ID param",
			httpMethod:     http.MethodGet,
			queryParam:     "",
			expectedStatus: http.StatusBadRequest,
			expectErrorMsg: "Missing 'id' parameter",
		},
		{
			name:           "Invalid ID param",
			httpMethod:     http.MethodGet,
			queryParam:     "id=abc",
			expectedStatus: http.StatusBadRequest,
			expectErrorMsg: "Invalid 'id' parameter",
		},
		{
			name:           "Author not found",
			httpMethod:     http.MethodGet,
			queryParam:     "id=99",
			mockReturn:     nil,
			mockError:      errors.New("author not found"),
			expectedStatus: http.StatusNotFound,
			expectErrorMsg: "author not found",
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
		t.Run(tc.name, func(t *testing.T) {
			mockService := new(mock.MockAuthorService)
			handler := author.NewAuthorHandler(mockService)

			url := "/author"
			if tc.queryParam != "" {
				url += "?" + tc.queryParam
			}
			req := httptest.NewRequest(tc.httpMethod, url, nil)
			w := httptest.NewRecorder()

			if tc.httpMethod == http.MethodGet && tc.queryParam != "" && tc.mockError != nil || tc.mockReturn != nil {
				mockService.On("GetByAuthorID", testifymock.AnythingOfType("int")).Return(tc.mockReturn, tc.mockError)
			}

			handler.GetByAuthorID(w, req)
			assert.Equal(t, tc.expectedStatus, w.Code)
			if tc.expectedStatus == http.StatusOK && tc.mockReturn != nil {
				var result models.Author
				err := json.NewDecoder(w.Body).Decode(&result)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, &result)
			} else if tc.expectErrorMsg != "" {
				assert.Contains(t, w.Body.String(), tc.expectErrorMsg)
			}
			mockService.AssertExpectations(t)
		})
	}
}

func TestDeleteById(t *testing.T) {
	tests := []struct {
		name           string
		queryParam     string
		mockReturn     *models.Author
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Success delete author",
			queryParam:     "id=1",
			mockReturn:     &models.Author{ID: 1, Name: "Author One"},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `"id":1`, // hoặc `"Author One"` nếu bạn kiểm tra chuỗi cụ thể
		},
		{
			name:           "Missing ID param",
			queryParam:     "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Missing 'id' parameter",
		},
		{
			name:           "Invalid ID param",
			queryParam:     "id=abc",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid 'id' parameter",
		},
		{
			name:           "Author not found",
			queryParam:     "id=99",
			mockReturn:     nil,
			mockError:      errors.New("author not found"),
			expectedStatus: http.StatusNotFound,
			expectedBody:   "author not found",
		},
		{
			name:           "Existing author constraint error",
			queryParam:     "id=2",
			mockReturn:     nil,
			mockError:      errors.New("existing author with books"),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "existing author with books",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockService := new(mock.MockAuthorService)
			handler := author.NewAuthorHandler(mockService)

			url := "/author/delete"
			if tc.queryParam != "" {
				url += "?" + tc.queryParam
			}
			req := httptest.NewRequest(http.MethodDelete, url, nil)
			w := httptest.NewRecorder()

			// Chỉ mock nếu ID hợp lệ
			if tc.mockReturn != nil || tc.mockError != nil {
				if idStr := req.URL.Query().Get("id"); idStr != "" {
					id, err := strconv.Atoi(idStr)
					if err == nil {
						mockService.On("DeleteById", id).Return(tc.mockReturn, tc.mockError)
					}
				}
			}

			handler.DeleteById(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tc.expectedBody)
			mockService.AssertExpectations(t)
		})
	}
}

func TestUpdateById(t *testing.T) {
	tests := []struct {
		name           string
		queryParam     string
		requestBody    string
		mockReturn     *models.Author
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "Success update author",
			queryParam:  "id=1",
			requestBody: `{"name":"Updated Author", "nationality":"Vietnamese"}`,
			mockReturn: &models.Author{
				ID:          1,
				Name:        "Updated Author",
				Nationality: "Vietnamese",
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `"name":"Updated Author"`,
		},
		{
			name:           "Missing ID param",
			queryParam:     "",
			requestBody:    `{"name":"Test"}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Missing 'id' parameter",
		},
		{
			name:           "Invalid ID param",
			queryParam:     "id=abc",
			requestBody:    `{"name":"Test"}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid 'id' parameter",
		},
		{
			name:           "Invalid JSON body",
			queryParam:     "id=1",
			requestBody:    `invalid-json`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid JSON body",
		},
		{
			name:           "Author not found",
			queryParam:     "id=99",
			requestBody:    `{"name":"Unknown"}`,
			mockReturn:     nil,
			mockError:      errors.New("author not found"),
			expectedStatus: http.StatusNotFound,
			expectedBody:   "author not found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockService := new(mock.MockAuthorService)
			handler := author.NewAuthorHandler(mockService)

			url := "/author/update"
			if tc.queryParam != "" {
				url += "?" + tc.queryParam
			}

			req := httptest.NewRequest(http.MethodPut, url, bytes.NewBufferString(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Nếu cần mock
			if tc.mockReturn != nil || tc.mockError != nil {
				if idStr := req.URL.Query().Get("id"); idStr != "" {
					if id, err := strconv.Atoi(idStr); err == nil {
						mockService.On("UpdateById", testifymock.MatchedBy(func(a *models.Author) bool {
							return a.ID == id
						})).Return(tc.mockReturn, tc.mockError)
					}
				}
			}

			handler.UpdateById(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tc.expectedBody)
			mockService.AssertExpectations(t)
		})
	}
}
