package unit

import (
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

func TestTodoHandler_DeleteTodo(t *testing.T) {
	// Given
	repo := repository.NewInMemoryTodoRepository()
	svc := service.NewTodoService(repo)
	h := handler.NewTodoHandler(svc)
	
	// Create a todo first
	created, err := svc.CreateTodo("todo to delete")
	require.NoError(t, err)

	req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/todos/%d", created.ID), nil)
	rec := httptest.NewRecorder()

	// When
	h.DeleteTodo(rec, req)

	// Then
	assert.Equal(t, http.StatusNoContent, rec.Code)
	
	// Verify todo is deleted
	_, err = svc.GetTodoByID(created.ID)
	assert.Error(t, err)
}

func TestTodoHandler_DeleteTodo_NotFound(t *testing.T) {
	// Given
	repo := repository.NewInMemoryTodoRepository()
	svc := service.NewTodoService(repo)
	h := handler.NewTodoHandler(svc)

	req := httptest.NewRequest("DELETE", "/api/todos/999", nil)
	rec := httptest.NewRecorder()

	// When
	h.DeleteTodo(rec, req)

	// Then
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestTodoHandler_DeleteTodo_InvalidID(t *testing.T) {
	// Given
	repo := repository.NewInMemoryTodoRepository()
	svc := service.NewTodoService(repo)
	h := handler.NewTodoHandler(svc)

	req := httptest.NewRequest("DELETE", "/api/todos/invalid", nil)
	rec := httptest.NewRecorder()

	// When
	h.DeleteTodo(rec, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, rec.Code)
} 