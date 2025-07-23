package unit

import (
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

func TestTodoHandler_GetTodoByID(t *testing.T) {
	// Given
	repo := repository.NewInMemoryTodoRepository()
	svc := service.NewTodoService(repo)
	h := handler.NewTodoHandler(svc)
	
	// Create a todo first
	created, err := svc.CreateTodo("test todo")
	require.NoError(t, err)

	req := httptest.NewRequest("GET", fmt.Sprintf("/api/todos/%d", created.ID), nil)
	rec := httptest.NewRecorder()

	// When
	h.GetTodoByID(rec, req)

	// Then
	assert.Equal(t, http.StatusOK, rec.Code)
	
	var todo map[string]interface{}
	json.NewDecoder(rec.Body).Decode(&todo)
	assert.Equal(t, "test todo", todo["text"])
	assert.Equal(t, float64(created.ID), todo["id"])
}

func TestTodoHandler_GetTodoByID_NotFound(t *testing.T) {
	// Given
	repo := repository.NewInMemoryTodoRepository()
	svc := service.NewTodoService(repo)
	h := handler.NewTodoHandler(svc)

	req := httptest.NewRequest("GET", "/api/todos/999", nil)
	rec := httptest.NewRecorder()

	// When
	h.GetTodoByID(rec, req)

	// Then
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestTodoHandler_GetTodoByID_InvalidID(t *testing.T) {
	// Given
	repo := repository.NewInMemoryTodoRepository()
	svc := service.NewTodoService(repo)
	h := handler.NewTodoHandler(svc)

	req := httptest.NewRequest("GET", "/api/todos/invalid", nil)
	rec := httptest.NewRecorder()

	// When
	h.GetTodoByID(rec, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestTodoHandler_GetAllTodos(t *testing.T) {
	// Given
	repo := repository.NewInMemoryTodoRepository()
	svc := service.NewTodoService(repo)
	h := handler.NewTodoHandler(svc)
	
	// Create some todos
	_, err := svc.CreateTodo("todo 1")
	require.NoError(t, err)
	_, err = svc.CreateTodo("todo 2")
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/api/todos", nil)
	rec := httptest.NewRecorder()

	// When
	h.GetAllTodos(rec, req)

	// Then
	assert.Equal(t, http.StatusOK, rec.Code)
	
	var todos []map[string]interface{}
	json.NewDecoder(rec.Body).Decode(&todos)
	assert.Len(t, todos, 2)
	
	// Check that both todos exist (order may vary)
	todoTexts := make([]string, len(todos))
	for i, todo := range todos {
		todoTexts[i] = todo["text"].(string)
	}
	assert.Contains(t, todoTexts, "todo 1")
	assert.Contains(t, todoTexts, "todo 2")
}

func TestTodoHandler_GetAllTodos_Empty(t *testing.T) {
	// Given
	repo := repository.NewInMemoryTodoRepository()
	svc := service.NewTodoService(repo)
	h := handler.NewTodoHandler(svc)

	req := httptest.NewRequest("GET", "/api/todos", nil)
	rec := httptest.NewRecorder()

	// When
	h.GetAllTodos(rec, req)

	// Then
	assert.Equal(t, http.StatusOK, rec.Code)
	
	var todos []map[string]interface{}
	json.NewDecoder(rec.Body).Decode(&todos)
	assert.Len(t, todos, 0)
} 