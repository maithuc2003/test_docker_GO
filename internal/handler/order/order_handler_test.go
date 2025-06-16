package order_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/maithuc2003/re-book-api/internal/handler/order"
	"github.com/maithuc2003/re-book-api/internal/models"
	"github.com/maithuc2003/re-book-api/test/mockservice"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllOrders(t *testing.T) {
	tests := []struct {
		name             string
		httpMethod       string
		mockReturn       []*models.Order
		mockError        error
		expectedStatus   int
		expectedResult   []*models.Order
		expectedErrorMsg string
	}{
		{
			name:       "Success - return orders",
			httpMethod: http.MethodGet,
			mockReturn: []*models.Order{
				{
					ID:       1,
					BookID:   101,
					UserID:   202,
					Quantity: 2,
					Status:   "Pending",
				},
			},

			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedResult: []*models.Order{
				{
					ID:       1,
					BookID:   101,
					UserID:   202,
					Quantity: 2,
					Status:   "Pending",
				},
			},
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
			name:             "Error from service",
			httpMethod:       http.MethodGet,
			mockReturn:       nil,
			mockError:        errors.New("DB error"),
			expectedStatus:   http.StatusInternalServerError,
			expectedErrorMsg: "Failed to get order",
		},
		{
			name:             "Method Not Allowed",
			httpMethod:       http.MethodPost, // hoặc PUT, DELETE,...
			mockReturn:       nil,
			mockError:        nil,
			expectedStatus:   http.StatusMethodNotAllowed,
			expectedErrorMsg: "Method not allowed",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mock_service := new(mockservice.MockOrderService)
			handler := order.NewOrderHandler(mock_service)

			if tc.httpMethod == http.MethodGet {
				mock_service.On("GetAllOrders").Return(tc.mockReturn, tc.mockError)
			}
			req := httptest.NewRequest(tc.httpMethod, "/orders", nil)
			w := httptest.NewRecorder()

			handler.GetAllOrders(w, req)
			assert.Equal(t, tc.expectedStatus, w.Code)

			if tc.expectedStatus == http.StatusOK && tc.mockError == nil {
				var result []*models.Order
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

func TestCreateOrder(t *testing.T) {
	tests := []struct {
		name           string
		httpMethod     string
		requestBody    interface{}
		mockError      error
		expectedStatus int
		expectErrorMsg string
	}{
		{
			name:       "Success - valid order",
			httpMethod: http.MethodPost,
			requestBody: &models.Order{
				BookID:   1,
				UserID:   2,
				Quantity: 3,
				Status:   "Pending",
			},
			mockError:      nil,
			expectedStatus: http.StatusCreated,
		},
		{
			name:       "Foreign key constraint violation",
			httpMethod: http.MethodPost,
			requestBody: &models.Order{
				BookID:   1,
				UserID:   2,
				Quantity: 3,
				Status:   "Pending",
			},
			mockError:      &mysql.MySQLError{Number: 1452, Message: "Cannot add or update a child row"},
			expectedStatus: http.StatusBadRequest,
			expectErrorMsg: "foreign key constraint violation",
		},
		{
			name:           "Invalid method",
			httpMethod:     http.MethodGet,
			requestBody:    nil,
			expectedStatus: http.StatusMethodNotAllowed,
			expectErrorMsg: "Method not allowed",
		},
		{
			name:           "Order is nil",
			httpMethod:     http.MethodPost,
			requestBody:    &models.Order{}, // vẫn cần truyền JSON hợp lệ
			mockError:      errors.New("order is nil"),
			expectedStatus: http.StatusBadRequest,
			expectErrorMsg: "order is nil",
		},
		{
			name:       "Invalid book ID",
			httpMethod: http.MethodPost,
			requestBody: &models.Order{
				BookID:   0,
				UserID:   2,
				Quantity: 3,
				Status:   "Pending",
			},
			mockError:      errors.New("invalid book ID"),
			expectedStatus: http.StatusBadRequest,
			expectErrorMsg: "invalid book ID",
		},
		{
			name:       "Invalid user ID",
			httpMethod: http.MethodPost,
			requestBody: &models.Order{
				BookID:   1,
				UserID:   0,
				Quantity: 3,
				Status:   "Pending",
			},
			mockError:      errors.New("invalid user ID"),
			expectedStatus: http.StatusBadRequest,
			expectErrorMsg: "invalid user ID",
		},
		{
			name:       "Quantity must be greater than zero",
			httpMethod: http.MethodPost,
			requestBody: &models.Order{
				BookID:   1,
				UserID:   2,
				Quantity: 0,
				Status:   "Pending",
			},
			mockError:      errors.New("quantity must be greater than zero"),
			expectedStatus: http.StatusBadRequest,
			expectErrorMsg: "quantity must be greater than zero",
		},
		{
			name:       "Status is required",
			httpMethod: http.MethodPost,
			requestBody: &models.Order{
				BookID:   1,
				UserID:   2,
				Quantity: 1,
				Status:   "",
			},
			mockError:      errors.New("status is required"),
			expectedStatus: http.StatusBadRequest,
			expectErrorMsg: "status is required",
		},
		{
			name:       "No rows in result set",
			httpMethod: http.MethodPost,
			requestBody: &models.Order{
				BookID:   1,
				UserID:   2,
				Quantity: 3,
				Status:   "Pending",
			},
			mockError:      errors.New("sql: no rows in result set"),
			expectedStatus: http.StatusNotFound,
			expectErrorMsg: "Product not found or no stock information",
		},
		{
			name:       "Not enough stock available",
			httpMethod: http.MethodPost,
			requestBody: &models.Order{
				BookID:   1,
				UserID:   2,
				Quantity: 3,
				Status:   "Pending",
			},
			mockError:      errors.New("not enough stock"),
			expectedStatus: http.StatusBadRequest,
			expectErrorMsg: "Not enough stock available",
		},
		{
			name:           "Invalid JSON body",
			httpMethod:     http.MethodPost,
			requestBody:    "invalid-json",
			expectedStatus: http.StatusBadRequest,
			expectErrorMsg: "Invalid request body",
		},
		{
			name:       "Service error",
			httpMethod: http.MethodPost,
			requestBody: &models.Order{
				BookID:   1,
				UserID:   2,
				Quantity: 3,
				Status:   "Pending",
			},
			mockError:      errors.New("insert error"),
			expectedStatus: http.StatusInternalServerError,
			expectErrorMsg: "Internal server error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mock_service := new(mockservice.MockOrderService)
			handler := order.NewOrderHandler(mock_service)

			var bodyBytes []byte
			var err error

			if s, ok := tc.requestBody.(string); ok && s == "invalid-json" {
				bodyBytes = []byte(`{invalid json`)
			} else {
				bodyBytes, err = json.Marshal(tc.requestBody)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(tc.httpMethod, "/order/add", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			if tc.httpMethod != http.MethodPost {
				handler.CreateOrder(w, req)
				assert.Equal(t, tc.expectedStatus, w.Code)
				assert.Contains(t, w.Body.String(), tc.expectErrorMsg)
				return
			}

			if tc.expectErrorMsg != "Invalid request body" {
				mock_service.On("CreateOrder", mock.AnythingOfType("*models.Order")).Return(tc.mockError)
			}

			handler.CreateOrder(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)

			if tc.expectErrorMsg != "" {
				assert.Contains(t, w.Body.String(), tc.expectErrorMsg)
			}

			mock_service.AssertExpectations(t)
		})
	}
}

func TestDeleteById(t *testing.T) {
	tests := []struct {
		name             string
		httpMethod       string
		queryParam       string
		mockReturn       *models.Order
		mockError        error
		expectedStatus   int
		expectedErrorMsg string
	}{
		{
			name:             "Success - Deleted",
			httpMethod:       http.MethodDelete,
			queryParam:       "id=1",
			mockReturn:       &models.Order{ID: 1, BookID: 100, UserID: 200, Quantity: 1, Status: "Deleted"},
			mockError:        nil,
			expectedStatus:   http.StatusOK,
			expectedErrorMsg: `"id":1`,
		},
		{
			name:             "Method not allowed",
			httpMethod:       http.MethodGet,
			queryParam:       "id=1",
			expectedStatus:   http.StatusMethodNotAllowed,
			expectedErrorMsg: "Method not allowed",
		},
		{
			name:             "Missing ID parameter",
			httpMethod:       http.MethodDelete,
			queryParam:       "",
			expectedStatus:   http.StatusBadRequest,
			expectedErrorMsg: "Missing 'id' parameter",
		},
		{
			name:             "Invalid ID parameter",
			httpMethod:       http.MethodDelete,
			queryParam:       "id=abc",
			expectedStatus:   http.StatusBadRequest,
			expectedErrorMsg: "Invalid 'id' parameter",
		},
		{
			name:             "Order not found",
			httpMethod:       http.MethodDelete,
			queryParam:       "id=999",
			mockReturn:       nil,
			mockError:        errors.New("Order not found"),
			expectedStatus:   http.StatusNotFound,
			expectedErrorMsg: "Order not found",
		},
		{
			name:             "Invalid order ID - negative",
			httpMethod:       http.MethodDelete,
			queryParam:       "id=-1",
			mockReturn:       nil,
			mockError:        errors.New("invalid order ID"),
			expectedStatus:   http.StatusBadRequest,
			expectedErrorMsg: "invalid order ID",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mock_service := new(mockservice.MockOrderService)
			handler := order.NewOrderHandler(mock_service)

			url := "/order/delete"
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
						mock_service.On("DeleteByOrderID", id).Return(tc.mockReturn, tc.mockError)
					}
				}
			}

			handler.DeleteByOrderID(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tc.expectedErrorMsg)
			mock_service.AssertExpectations(t)
		})
	}
}

func TestUpdateByID(t *testing.T) {
	tests := []struct {
		name             string
		queryParam       string
		requestBody      string
		mockReturn       *models.Order
		mockError        error
		expectedStatus   int
		expectedErrorMsg string
		httpMethod       string // đổi tên từ 'method' -> 'httpMethod'
	}{
		{
			name:             "Success - updated",
			queryParam:       "id=1",
			requestBody:      `{"book_id":101,"user_id":201,"quantity":2,"status":"Confirmed"}`,
			mockReturn:       &models.Order{ID: 1, BookID: 101, UserID: 201, Quantity: 2, Status: "Confirmed"},
			mockError:        nil,
			expectedStatus:   http.StatusOK,
			expectedErrorMsg: `"id":1`,
			httpMethod:       http.MethodPut,
		},
		{
			name:             "Invalid method",
			queryParam:       "id=1",
			requestBody:      `{"book_id":101}`,
			expectedStatus:   http.StatusMethodNotAllowed,
			expectedErrorMsg: "Method not allowed",
			httpMethod:       http.MethodGet,
		},
		{
			name:             "Missing ID param",
			queryParam:       "",
			requestBody:      `{"book_id":101}`,
			expectedStatus:   http.StatusBadRequest,
			expectedErrorMsg: "Missing 'id' parameter",
			httpMethod:       http.MethodPut,
		},
		{
			name:             "Invalid ID param",
			queryParam:       "id=abc",
			requestBody:      `{"book_id":101}`,
			expectedStatus:   http.StatusBadRequest,
			expectedErrorMsg: "Invalid 'id' parameter",
			httpMethod:       http.MethodPut,
		},
		{
			name:             "Invalid JSON body",
			queryParam:       "id=1",
			requestBody:      `invalid-json`,
			expectedStatus:   http.StatusBadRequest,
			expectedErrorMsg: "Invalid JSON body",
			httpMethod:       http.MethodPut,
		},
		{
			name:             "Book ID does not exist",
			queryParam:       "id=1",
			requestBody:      `{"book_id":999,"user_id":1,"quantity":1,"status":"Pending"}`,
			mockReturn:       nil,
			mockError:        errors.New("foreign key constraint fails: book_id does not exist"),
			expectedStatus:   http.StatusBadRequest,
			expectedErrorMsg: "Invalid book_id: book does not exist",
			httpMethod:       http.MethodPut,
		},
		{
			name:             "Internal service error",
			queryParam:       "id=1",
			requestBody:      `{"book_id":101,"user_id":1,"quantity":1,"status":"Pending"}`,
			mockReturn:       nil,
			mockError:        errors.New("some db error"),
			expectedStatus:   http.StatusInternalServerError,
			expectedErrorMsg: "Failed to update order",
			httpMethod:       http.MethodPut,
		},
		{
			name:             "Invalid order ID (<= 0)",
			queryParam:       "id=0",
			requestBody:      `{"book_id":101,"user_id":201,"quantity":2,"status":"Confirmed"}`,
			mockError:        errors.New("invalid order ID"),
			expectedStatus:   http.StatusBadRequest,
			expectedErrorMsg: "Validation error: invalid order ID",
			httpMethod:       http.MethodPut,
		},
		{
			name:             "Invalid book ID (<= 0)",
			queryParam:       "id=1",
			requestBody:      `{"book_id":0,"user_id":201,"quantity":2,"status":"Confirmed"}`,
			mockError:        errors.New("invalid book ID"),
			expectedStatus:   http.StatusBadRequest,
			expectedErrorMsg: "Validation error: invalid book ID",
			httpMethod:       http.MethodPut,
		},
		{
			name:             "Invalid user ID (<= 0)",
			queryParam:       "id=1",
			requestBody:      `{"book_id":101,"user_id":0,"quantity":2,"status":"Confirmed"}`,
			mockError:        errors.New("invalid user ID"),
			expectedStatus:   http.StatusBadRequest,
			expectedErrorMsg: "Validation error: invalid user ID",
			httpMethod:       http.MethodPut,
		},
		{
			name:             "Quantity must be greater than zero",
			queryParam:       "id=1",
			requestBody:      `{"book_id":101,"user_id":201,"quantity":0,"status":"Confirmed"}`,
			mockError:        errors.New("quantity must be greater than zero"),
			expectedStatus:   http.StatusBadRequest,
			expectedErrorMsg: "Validation error: quantity must be greater than zero",
			httpMethod:       http.MethodPut,
		},
		{
			name:             "Status is required",
			queryParam:       "id=1",
			requestBody:      `{"book_id":101,"user_id":201,"quantity":2,"status":""}`,
			mockError:        errors.New("status is required"),
			expectedStatus:   http.StatusBadRequest,
			expectedErrorMsg: "Validation error: status is required",
			httpMethod:       http.MethodPut,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mock_service := new(mockservice.MockOrderService)
			handler := order.NewOrderHandler(mock_service)

			url := "/order/update"
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
						mock_service.On("UpdateByOrderID", mock.MatchedBy(func(a *models.Order) bool {
							return a.ID == id
						})).Return(tc.mockReturn, tc.mockError)
					}
				}
			}
			handler.UpdateByOrderID(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tc.expectedErrorMsg)
			mock_service.AssertExpectations(t)
		})
	}
}

func TestGetOrderByID(t *testing.T) {
	tests := []struct {
		name           string
		httpMethod     string
		queryParam     string
		mockReturn     *models.Order
		mockError      error
		expectedStatus int
		expectedResult *models.Order
		expectErrorMsg string
	}{
		{
			name:           "Success - valid ID",
			httpMethod:     http.MethodGet,
			queryParam:     "id=1",
			mockReturn:     &models.Order{ID: 1, BookID: 101, UserID: 202, Quantity: 2, Status: "Confirmed"},
			expectedStatus: http.StatusOK,
			expectedResult: &models.Order{ID: 1, BookID: 101, UserID: 202, Quantity: 2, Status: "Confirmed"},
		},
		{
			name:           "Invalid method",
			httpMethod:     http.MethodPost,
			queryParam:     "id=1",
			expectedStatus: http.StatusMethodNotAllowed,
			expectErrorMsg: "Method not allowed",
		},
		{
			name:           "Missing id parameter",
			httpMethod:     http.MethodGet,
			queryParam:     "",
			expectedStatus: http.StatusBadRequest,
			expectErrorMsg: "Missing 'id' parameter",
		},
		{
			name:           "Invalid id (not a number)",
			httpMethod:     http.MethodGet,
			queryParam:     "id=abc",
			expectedStatus: http.StatusBadRequest,
			expectErrorMsg: "Invalid 'id' parameter",
		},
		{
			name:           "Service error - order not found",
			httpMethod:     http.MethodGet,
			queryParam:     "id=99",
			mockError:      errors.New("Order not found"),
			expectedStatus: http.StatusNotFound,
			expectErrorMsg: "Order not found",
		},
		{
			name:           "Invalid order ID (zero or negative)",
			httpMethod:     http.MethodGet,
			queryParam:     "id=-5",
			mockError:      errors.New("invalid order ID"),
			expectedStatus: http.StatusBadRequest,
			expectErrorMsg: "invalid order ID",
		},
		{
			name:           "Other service error (e.g. unknown)",
			httpMethod:     http.MethodGet,
			queryParam:     "id=100",
			mockError:      errors.New("unexpected error"),
			expectedStatus: http.StatusNotFound,
			expectErrorMsg: "unexpected error",
		},
		{
			name:           "Existing orders error",
			httpMethod:     http.MethodGet,
			queryParam:     "id=10",
			mockError:      errors.New("existing orders"),
			expectedStatus: http.StatusBadRequest,
			expectErrorMsg: "existing orders",
		},
	}

	for _, tc := range tests {
		mock_service := new(mockservice.MockOrderService)
		handler := order.NewOrderHandler(mock_service)

		url := "/order"
		if tc.queryParam != "" {
			url += "?" + tc.queryParam
		}
		req := httptest.NewRequest(tc.httpMethod, url, nil)
		w := httptest.NewRecorder()

		if tc.httpMethod == http.MethodGet && tc.queryParam != "" && (tc.mockError != nil || tc.mockReturn != nil) {
			mock_service.On("GetByOrderID", mock.AnythingOfType("int")).Return(tc.mockReturn, tc.mockError)
		}

		handler.GetByOrderID(w, req)
		assert.Equal(t, tc.expectedStatus, w.Code)
		if tc.expectedStatus == http.StatusOK && tc.mockReturn != nil {
			var result models.Order
			err := json.NewDecoder(w.Body).Decode(&result)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedResult, &result)
		} else if tc.expectErrorMsg != "" {
			assert.Contains(t, w.Body.String(), tc.expectErrorMsg)
		}
		mock_service.AssertExpectations(t)
	}
}
