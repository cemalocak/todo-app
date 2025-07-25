package unit

import (
	"testing"

	"todo-app/internal/repository"
	"todo-app/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTodoService_CreateTodo(t *testing.T) {
	// Given
	repo, err := repository.NewSQLiteTodoRepository(":memory:")
	require.NoError(t, err)
	svc := service.NewTodoService(repo)

	// When
	todo, err := svc.CreateTodo("test todo")

	// Then
	assert.NoError(t, err)
	assert.Equal(t, "test todo", todo.Text)
	assert.Equal(t, 1, todo.ID)
}

func TestTodoService_GetAllTodos(t *testing.T) {
	// Given
	repo, err := repository.NewSQLiteTodoRepository(":memory:")
	require.NoError(t, err)
	svc := service.NewTodoService(repo)
	svc.CreateTodo("todo 1")
	svc.CreateTodo("todo 2")

	// When
	todos, err := svc.GetAllTodos()

	// Then
	assert.NoError(t, err)
	assert.Len(t, todos, 2)
	// Should be in DESC order (newest first)
	assert.Equal(t, "todo 2", todos[0].Text)
	assert.Equal(t, "todo 1", todos[1].Text)
}
