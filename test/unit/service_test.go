package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"todo-app/internal/repository"
	"todo-app/internal/service"
)

func TestTodoService_CreateTodo(t *testing.T) {
	// Given
	repo := repository.NewInMemoryTodoRepository()
	svc := service.NewTodoService(repo)

	// When
	todo, err := svc.CreateTodo("test todo")

	// Then
	assert.NoError(t, err)
	assert.Equal(t, "test todo", todo.Text)
	assert.Equal(t, 1, todo.ID)
}

func TestTodoService_GetTodoByID(t *testing.T) {
	// Given
	repo := repository.NewInMemoryTodoRepository()
	svc := service.NewTodoService(repo)
	
	created, err := svc.CreateTodo("test todo")
	require.NoError(t, err)

	// When
	result, err := svc.GetTodoByID(created.ID)

	// Then
	require.NoError(t, err)
	assert.Equal(t, created.ID, result.ID)
	assert.Equal(t, "test todo", result.Text)
}

func TestTodoService_GetTodoByID_NotFound(t *testing.T) {
	// Given
	repo := repository.NewInMemoryTodoRepository()
	svc := service.NewTodoService(repo)

	// When
	result, err := svc.GetTodoByID(999)

	// Then
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")
}

func TestTodoService_GetAllTodos(t *testing.T) {
	// Given
	repo := repository.NewInMemoryTodoRepository()
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

func TestTodoService_UpdateTodo(t *testing.T) {
	// Given
	repo := repository.NewInMemoryTodoRepository()
	svc := service.NewTodoService(repo)
	
	created, err := svc.CreateTodo("original text")
	require.NoError(t, err)

	// When
	updated, err := svc.UpdateTodo(created.ID, "updated text")

	// Then
	require.NoError(t, err)
	assert.Equal(t, created.ID, updated.ID)
	assert.Equal(t, "updated text", updated.Text)

	// And: GetByID should return updated version
	fetched, err := svc.GetTodoByID(created.ID)
	require.NoError(t, err)
	assert.Equal(t, "updated text", fetched.Text)
}

func TestTodoService_UpdateTodo_NotFound(t *testing.T) {
	// Given
	repo := repository.NewInMemoryTodoRepository()
	svc := service.NewTodoService(repo)

	// When
	result, err := svc.UpdateTodo(999, "updated text")

	// Then
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")
}

func TestTodoService_DeleteTodo(t *testing.T) {
	// Given
	repo := repository.NewInMemoryTodoRepository()
	svc := service.NewTodoService(repo)
	
	created, err := svc.CreateTodo("to be deleted")
	require.NoError(t, err)

	// When
	err = svc.DeleteTodo(created.ID)

	// Then
	require.NoError(t, err)

	// And: Todo should not be found
	_, err = svc.GetTodoByID(created.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestTodoService_DeleteTodo_NotFound(t *testing.T) {
	// Given
	repo := repository.NewInMemoryTodoRepository()
	svc := service.NewTodoService(repo)

	// When
	err := svc.DeleteTodo(999)

	// Then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
} 