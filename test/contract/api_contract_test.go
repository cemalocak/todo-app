package contract

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"todo-app/internal/handler"
	"todo-app/internal/repository"
	"todo-app/internal/service"
)

// TestAPI_ResponseStructure tests the API response structure compliance
func TestAPI_ResponseStructure(t *testing.T) {
	tests := []struct {
		name            string
		method          string
		path            string
		body            interface{}
		expectedStatus  int
		expectedFields  []string
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
			repo := repository.NewInMemoryTodoRepository()
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
	repo := repository.NewInMemoryTodoRepository()
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

// TestAPI_ErrorResponses tests consistent error response format
func TestAPI_ErrorResponses(t *testing.T) {
	repo := repository.NewInMemoryTodoRepository()
	svc := service.NewTodoService(repo)
	h := handler.NewTodoHandler(svc)

	tests := []struct {
		name           string
		endpoint       func(w http.ResponseWriter, r *http.Request)
		request        *http.Request
		expectedStatus int
	}{
		{
			name:           "Get non-existent todo",
			endpoint:       h.GetTodoByID,
			request:        httptest.NewRequest("GET", "/api/todos/999", nil),
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Update non-existent todo",
			endpoint:       h.UpdateTodo,
			request:        httptest.NewRequest("PUT", "/api/todos/999", bytes.NewBuffer([]byte(`{"text":"updated"}`))),
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Delete non-existent todo",
			endpoint:       h.DeleteTodo,
			request:        httptest.NewRequest("DELETE", "/api/todos/999", nil),
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			tt.endpoint(rec, tt.request)
			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}

// TestAPI_CORS tests CORS headers (if implemented)
func TestAPI_CORSHeaders(t *testing.T) {
	repo := repository.NewInMemoryTodoRepository()
	svc := service.NewTodoService(repo)
	h := handler.NewTodoHandler(svc)

	req := httptest.NewRequest("GET", "/api/todos", nil)
	rec := httptest.NewRecorder()

	h.GetAllTodos(rec, req)

	// Check if CORS headers are set (implementation dependent)
	// This test documents expected CORS behavior
	assert.Equal(t, http.StatusOK, rec.Code)
	
	// Note: Add CORS header assertions here when CORS is implemented
	// assert.Equal(t, "*", rec.Header().Get("Access-Control-Allow-Origin"))
} 