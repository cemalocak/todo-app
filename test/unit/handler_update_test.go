package unit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"todo-app/internal/handler"
	"todo-app/internal/repository"
	"todo-app/internal/service"
)

func TestTodoHandler_UpdateTodo(t *testing.T) {
	// Given
	repo := repository.NewInMemoryTodoRepository()
	svc := service.NewTodoService(repo)
	h := handler.NewTodoHandler(svc)
	
	// Create a todo first
	created, err := svc.CreateTodo("original todo")
	require.NoError(t, err)

	requestBody := map[string]string{"text": "updated todo"}
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("PUT", fmt.Sprintf("/api/todos/%d", created.ID), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// When
	h.UpdateTodo(rec, req)

	// Then
	assert.Equal(t, http.StatusOK, rec.Code)
	
	var todo map[string]interface{}
	json.NewDecoder(rec.Body).Decode(&todo)
	assert.Equal(t, "updated todo", todo["text"])
	assert.Equal(t, float64(created.ID), todo["id"])
}

func TestTodoHandler_UpdateTodo_NotFound(t *testing.T) {
	// Given
	repo := repository.NewInMemoryTodoRepository()
	svc := service.NewTodoService(repo)
	h := handler.NewTodoHandler(svc)

	requestBody := map[string]string{"text": "updated todo"}
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("PUT", "/api/todos/999", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// When
	h.UpdateTodo(rec, req)

	// Then
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestTodoHandler_UpdateTodo_InvalidID(t *testing.T) {
	// Given
	repo := repository.NewInMemoryTodoRepository()
	svc := service.NewTodoService(repo)
	h := handler.NewTodoHandler(svc)

	requestBody := map[string]string{"text": "updated todo"}
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("PUT", "/api/todos/invalid", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// When
	h.UpdateTodo(rec, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestTodoHandler_UpdateTodo_EmptyText(t *testing.T) {
	// Given
	repo := repository.NewInMemoryTodoRepository()
	svc := service.NewTodoService(repo)
	h := handler.NewTodoHandler(svc)
	
	// Create a todo first
	created, err := svc.CreateTodo("original todo")
	require.NoError(t, err)

	requestBody := map[string]string{"text": ""}
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("PUT", fmt.Sprintf("/api/todos/%d", created.ID), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// When
	h.UpdateTodo(rec, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, rec.Code)
} 