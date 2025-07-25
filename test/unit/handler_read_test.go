package unit

import (
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

func TestTodoHandler_GetAllTodos(t *testing.T) {
	// Given
	repo, err := repository.NewSQLiteTodoRepository(":memory:")
	require.NoError(t, err)
	svc := service.NewTodoService(repo)
	h := handler.NewTodoHandler(svc)

	// Create some todos
	_, err = svc.CreateTodo("todo 1")
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
	repo, err := repository.NewSQLiteTodoRepository(":memory:")
	require.NoError(t, err)
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
