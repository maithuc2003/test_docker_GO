package author_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/maithuc2003/re-book-api/internal/models"
	"github.com/maithuc2003/re-book-api/internal/service/author"
	"github.com/maithuc2003/re-book-api/test/mockrepo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAllAuthors(t *testing.T) {
	tests := []struct {
		name           string
		mockReturn     []*models.Author
		mockError      error
		expectedResult []*models.Author
		expectErrorMsg string
	}{
		{
			name: "Success with authors",
			mockReturn: []*models.Author{
				{ID: 1, Name: "Author A"},
				{ID: 2, Name: "Author B"},
			},
			mockError:      nil,
			expectedResult: []*models.Author{{ID: 1, Name: "Author A"}, {ID: 2, Name: "Author B"}},
		},
		{
			name:           "Repository error",
			mockReturn:     nil,
			mockError:      errors.New("DB error"),
			expectedResult: nil,
			expectErrorMsg: "DB error",
		},
		{
			name:           "Empty author list",
			mockReturn:     []*models.Author{},
			mockError:      nil,
			expectedResult: nil,
			expectErrorMsg: "no authors found in the system",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockrepo := new(mockrepo.MockAuthorRepository)
			mockrepo.On("GetAllAuthors").Return(tc.mockReturn, tc.mockError)

			service := author.NewAuthorService(mockrepo)
			result, err := service.GetAllAuthors()
			if tc.expectErrorMsg != "" {
				require.Error(t, err)
				assert.Nil(t, result)
				assert.EqualError(t, err, tc.expectErrorMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
			mockrepo.AssertExpectations(t)
		})
	}
}

func TestCreateAuthor(t *testing.T) {
	tests := []struct {
		name            string
		inputAuthor     *models.Author
		existingAuthors []*models.Author
		getAllErr       error
		createErr       error
		expectErr       string
	}{
		{
			name:        "Success - author created",
			inputAuthor: &models.Author{Name: "New Author"},
			existingAuthors: []*models.Author{
				{Name: "Old Author"},
			},
			expectErr: "",
		},
		{
			name:        "Nil author input",
			inputAuthor: nil,
			expectErr:   "author is nil",
		},
		{
			name:        "Empty name",
			inputAuthor: &models.Author{Name: "   "},
			expectErr:   "author name cannot be empty",
		},
		{
			name:        "Duplicate name",
			inputAuthor: &models.Author{Name: "Jane Austen"},
			existingAuthors: []*models.Author{
				{Name: "Jane Austen"},
			},
			expectErr: "author with the same name already exists",
		},
		{
			name:        "Error when fetching existing authors",
			inputAuthor: &models.Author{Name: "Someone"},
			getAllErr:   errors.New("DB down"),
			expectErr:   "failed to fetch authors for validation: DB down",
		},
		{
			name:        "Error creating author",
			inputAuthor: &models.Author{Name: "New Guy"},
			existingAuthors: []*models.Author{
				{Name: "Old Author"},
			},
			createErr: errors.New("insert error"),
			expectErr: "failed to create author: insert error",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockrepo := new(mockrepo.MockAuthorRepository)

			// Only mock GetAllAuthors if input is non-nil and name is not empty
			if tc.inputAuthor != nil && strings.TrimSpace(tc.inputAuthor.Name) != "" {
				mockrepo.On("GetAllAuthors").Return(tc.existingAuthors, tc.getAllErr)
			}

			//Mock CreateAuthor only when we expect the service to reach that point
			if tc.expectErr == "" || strings.HasPrefix(tc.expectErr, "failed to create author") {
				mockrepo.On("CreateAuthor", tc.inputAuthor).Return(tc.createErr)
			}

			service := author.NewAuthorService(mockrepo)
			err := service.CreateAuthor(tc.inputAuthor)

			// Assert expected error or success
			if tc.expectErr != "" {
				require.Error(t, err)
				assert.EqualError(t, err, tc.expectErr)
			} else {
				assert.NoError(t, err)
			}
			// Ensure all mocked methods were called as expected
			mockrepo.AssertExpectations(t)
		})
	}
}

func TestGetByAuthor(t *testing.T) {
	tests := []struct {
		name        string
		inputID     int
		mockReturn  *models.Author
		mockError   error
		expectedRes *models.Author
		expectedErr string
	}{
		{
			name:        "Valid ID returns author",
			inputID:     1,
			mockReturn:  &models.Author{ID: 1, Name: "Author A"},
			mockError:   nil,
			expectedRes: &models.Author{ID: 1, Name: "Author A"},
		},
		{
			name:        "Invalid ID <= 0",
			inputID:     -1,
			mockReturn:  nil,
			mockError:   nil,
			expectedErr: "invalid author ID",
		},
		{
			name:        "Database error while fetching author",
			inputID:     2,
			mockReturn:  nil,
			mockError:   errors.New("database connection failed"),
			expectedErr: "failed to retrieve author: database connection failed",
		},
		{
			name:        "Author not found",
			inputID:     3,
			mockReturn:  nil,
			mockError:   nil,
			expectedErr: "author not found",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockrepo := new(mockrepo.MockAuthorRepository)

			//Only mock if ID is positive
			if tc.inputID > 0 {
				mockrepo.On("GetByAuthorID", tc.inputID).Return(tc.mockReturn, tc.mockError)
			}

			service := author.NewAuthorService(mockrepo)
			result, err := service.GetByAuthorID(tc.inputID)
			if tc.expectedErr != "" {
				require.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedRes, result)
			}
			mockrepo.AssertExpectations(t)
		})
	}
}

func TestDeleteById(t *testing.T) {
	tests := []struct {
		name        string
		inputID     int
		mockReturn  *models.Author
		mockError   error
		expectedRes *models.Author
		expectedErr string
	}{
		{
			name:        "Invalid author ID",
			inputID:     -1,
			expectedErr: "invalid author ID",
		},
		{
			name:        "Repository error while deleting",
			inputID:     1,
			mockReturn:  nil,
			mockError:   errors.New("database error"),
			expectedErr: "failed to delete author: database error",
		},
		{
			name:        "Author not found or already deleted",
			inputID:     2,
			mockReturn:  nil,
			mockError:   nil,
			expectedErr: "author not found or already deleted",
		},
		{
			name:        "Successfully deleted author",
			inputID:     3,
			mockReturn:  &models.Author{ID: 3, Name: "Author C"},
			mockError:   nil,
			expectedRes: &models.Author{ID: 3, Name: "Author C"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockrepo := new(mockrepo.MockAuthorRepository)

			// Mock GetByAuthorID nếu DeleteById cần nó
			if tc.inputID > 0 {
				mockrepo.On("DeleteById", tc.inputID).Return(tc.mockReturn, tc.mockError)
			}

			service := author.NewAuthorService(mockrepo)
			result, err := service.DeleteById(tc.inputID)
			if tc.expectedErr != "" {
				require.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedRes, result)
			}
			mockrepo.AssertExpectations(t)
		})
	}
}

func TestUpdateById(t *testing.T) {
	tests := []struct {
		name        string
		input       *models.Author
		mockGetByID *models.Author
		mockGetAll  []*models.Author
		mockUpdate  *models.Author
		mockErrors  struct {
			getByID error
			getAll  error
			update  error
		}
		expectedRes *models.Author
		expectedErr string
	}{
		{
			name:        "Nil author input",
			input:       nil,
			expectedErr: "author is nil",
		},
		{
			name:        "Invalid author ID",
			input:       &models.Author{ID: 0, Name: "John"},
			expectedErr: "invalid author ID",
		},
		{
			name:        "Empty author name",
			input:       &models.Author{ID: 1, Name: " "},
			expectedErr: "author name cannot be empty",
		},
		{
			name:        "Author not found",
			input:       &models.Author{ID: 2, Name: "John"},
			mockGetByID: nil,
			mockErrors:  struct{ getByID, getAll, update error }{getByID: nil},
			expectedErr: "author not found",
		},
		{
			name:        "Error fetching author by ID",
			input:       &models.Author{ID: 3, Name: "John"},
			mockErrors:  struct{ getByID, getAll, update error }{getByID: errors.New("db error")},
			expectedErr: "failed to fetch existing author: db error",
		},
		{
			name:        "Duplicate author name",
			input:       &models.Author{ID: 4, Name: "Same Name"},
			mockGetByID: &models.Author{ID: 4, Name: "Old Name"},
			mockGetAll: []*models.Author{
				{ID: 5, Name: "Same Name"},
			},
			expectedErr: "another author with the same name already exists",
		},
		{
			name:        "Error fetching all authors",
			input:       &models.Author{ID: 6, Name: "John"},
			mockGetByID: &models.Author{ID: 6, Name: "Johnny"},
			mockErrors:  struct{ getByID, getAll, update error }{getAll: errors.New("list error")},
			expectedErr: "failed to validate author name: list error",
		},
		{
			name:        "Error updating author",
			input:       &models.Author{ID: 7, Name: "Updated"},
			mockGetByID: &models.Author{ID: 7, Name: "Old"},
			mockGetAll:  []*models.Author{{ID: 7, Name: "Updated"}},
			mockErrors:  struct{ getByID, getAll, update error }{update: errors.New("update error")},
			expectedErr: "failed to update author : update error",
		},
		{
			name:        "Successful update",
			input:       &models.Author{ID: 8, Name: "Updated"},
			mockGetByID: &models.Author{ID: 8, Name: "Old"},
			mockGetAll: []*models.Author{
				{ID: 8, Name: "Updated"},
			},
			mockUpdate:  &models.Author{ID: 8, Name: "Updated"},
			expectedRes: &models.Author{ID: 8, Name: "Updated"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockrepo := new(mockrepo.MockAuthorRepository)

			if tc.input != nil && tc.input.ID > 0 && strings.TrimSpace(tc.input.Name) != "" {
				mockrepo.On("GetByAuthorID", tc.input.ID).Return(tc.mockGetByID, tc.mockErrors.getByID)
			}

			if tc.mockGetAll != nil || tc.mockErrors.getAll != nil {
				mockrepo.On("GetAllAuthors").Return(tc.mockGetAll, tc.mockErrors.getAll)
			}

			if tc.mockUpdate != nil && tc.mockErrors.update == nil {
				mockrepo.On("UpdateById", tc.input).Return(tc.mockUpdate, nil)
			} else if tc.mockErrors.update != nil {
				mockrepo.On("UpdateById", tc.input).Return(nil, tc.mockErrors.update)
			}

			service := author.NewAuthorService(mockrepo)
			result, err := service.UpdateById(tc.input)

			if tc.expectedErr != "" {
				require.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedRes, result)
			}
			mockrepo.AssertExpectations(t)
		})
	}
}
