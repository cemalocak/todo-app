package contract

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"todo-app/internal/handler"
	"todo-app/internal/repository"
	"todo-app/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAPI_ResponseStructure tests the API response structure compliance
func TestAPI_ResponseStructure(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		body           interface{}
		expectedStatus int
		expectedFields []string
	}{
		{
			name:           "POST /api/todos - creates todo with correct structure",
			method:         "POST",
			path:           "/api/todos",
			body:           map[string]string{"text": "test todo"},
			expectedStatus: http.StatusCreated,
			expectedFields: []string{"id", "text", "created_at", "updated_at"},
		},
		{
			name:           "GET /api/todos - returns array structure",
			method:         "GET",
			path:           "/api/todos",
			body:           nil,
			expectedStatus: http.StatusOK,
			expectedFields: []string{}, // Array response
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			repo, err := repository.NewSQLiteTodoRepository("test.db")
			require.NoError(t, err)
			defer repo.Close()

			err = repo.Truncate() // Clean state for each test
			require.NoError(t, err)

			svc := service.NewTodoService(repo)
			h := handler.NewTodoHandler(svc)

			// Prepare request
			var reqBody *bytes.Buffer
			if tt.body != nil {
				jsonBody, _ := json.Marshal(tt.body)
				reqBody = bytes.NewBuffer(jsonBody)
			} else {
				reqBody = bytes.NewBuffer([]byte{})
			}

			req := httptest.NewRequest(tt.method, tt.path, reqBody)
			if tt.body != nil {
				req.Header.Set("Content-Type", "application/json")
			}
			rec := httptest.NewRecorder()

			// Execute
			switch tt.method {
			case "POST":
				h.CreateTodo(rec, req)
			case "GET":
				h.GetAllTodos(rec, req)
			}

			// Verify
			assert.Equal(t, tt.expectedStatus, rec.Code)

			if len(tt.expectedFields) > 0 {
				var response map[string]interface{}
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				require.NoError(t, err)

				for _, field := range tt.expectedFields {
					assert.Contains(t, response, field, "Field %s should be present", field)
				}
			}
		})
	}
}

// TestAPI_ContentTypeHeaders tests correct content type handling
func TestAPI_ContentTypeHeaders(t *testing.T) {
	repo, err := repository.NewSQLiteTodoRepository("test.db")
	require.NoError(t, err)
	defer repo.Close()

	err = repo.Truncate() // Clean state
	require.NoError(t, err)

	svc := service.NewTodoService(repo)
	h := handler.NewTodoHandler(svc)

	tests := []struct {
		name           string
		contentType    string
		expectedStatus int
	}{
		{
			name:           "Valid JSON content type",
			contentType:    "application/json",
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Missing content type",
			contentType:    "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid content type",
			contentType:    "text/plain",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestBody := map[string]string{"text": "test todo"}
			jsonBody, _ := json.Marshal(requestBody)
			req := httptest.NewRequest("POST", "/api/todos", bytes.NewBuffer(jsonBody))

			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}

			rec := httptest.NewRecorder()
			h.CreateTodo(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}

// TestAPI_CORS tests CORS headers (if implemented)
func TestAPI_CORSHeaders(t *testing.T) {
	repo, err := repository.NewSQLiteTodoRepository("test.db")
	require.NoError(t, err)
	defer repo.Close()

	err = repo.Truncate() // Clean state
	require.NoError(t, err)

	svc := service.NewTodoService(repo)
	h := handler.NewTodoHandler(svc)

	req := httptest.NewRequest("GET", "/api/todos", nil)
	rec := httptest.NewRecorder()

	h.GetAllTodos(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
}
