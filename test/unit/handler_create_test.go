package unit

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"todo-app/internal/handler"
	"todo-app/internal/repository"
	"todo-app/internal/service"
)

func TestTodoHandler_CreateTodo(t *testing.T) {
	// Given
	repo := repository.NewInMemoryTodoRepository()
	svc := service.NewTodoService(repo)
	h := handler.NewTodoHandler(svc)

	requestBody := map[string]string{"text": "test todo"}
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/api/todos", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// When
	h.CreateTodo(rec, req)

	// Then
	assert.Equal(t, http.StatusCreated, rec.Code)
	
	var todo map[string]interface{}
	json.NewDecoder(rec.Body).Decode(&todo)
	assert.Equal(t, "test todo", todo["text"])
	assert.Equal(t, float64(1), todo["id"]) // JSON numbers are float64
}

func TestTodoHandler_CreateTodo_EmptyText(t *testing.T) {
	// Given
	repo := repository.NewInMemoryTodoRepository()
	svc := service.NewTodoService(repo)
	h := handler.NewTodoHandler(svc)

	requestBody := map[string]string{"text": ""}
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/api/todos", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// When
	h.CreateTodo(rec, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestTodoHandler_CreateTodo_InvalidJSON(t *testing.T) {
	// Given
	repo := repository.NewInMemoryTodoRepository()
	svc := service.NewTodoService(repo)
	h := handler.NewTodoHandler(svc)

	req := httptest.NewRequest("POST", "/api/todos", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// When
	h.CreateTodo(rec, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, rec.Code)
} 