package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"todo-app/internal/model"
	"todo-app/internal/repository"
)

func TestInMemoryTodoRepository_Create(t *testing.T) {
	// Given
	repo := repository.NewInMemoryTodoRepository()
	todo := &model.Todo{Text: "test todo"}

	// When
	result, err := repo.Create(todo)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "test todo", result.Text)
	assert.False(t, result.CreatedAt.IsZero())
	assert.False(t, result.UpdatedAt.IsZero())
}

func TestInMemoryTodoRepository_GetByID(t *testing.T) {
	// Given
	repo := repository.NewInMemoryTodoRepository()
	todo := &model.Todo{Text: "test todo"}
	created, err := repo.Create(todo)
	require.NoError(t, err)

	// When
	result, err := repo.GetByID(created.ID)

	// Then
	require.NoError(t, err)
	assert.Equal(t, created.ID, result.ID)
	assert.Equal(t, "test todo", result.Text)
	assert.Equal(t, created.CreatedAt, result.CreatedAt)
	assert.Equal(t, created.UpdatedAt, result.UpdatedAt)
}

func TestInMemoryTodoRepository_GetByID_NotFound(t *testing.T) {
	// Given
	repo := repository.NewInMemoryTodoRepository()

	// When
	result, err := repo.GetByID(999)

	// Then
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")
}

func TestInMemoryTodoRepository_GetAll(t *testing.T) {
	// Given
	repo := repository.NewInMemoryTodoRepository()
	todo1 := &model.Todo{Text: "todo 1"}
	todo2 := &model.Todo{Text: "todo 2"}
	repo.Create(todo1)
	repo.Create(todo2)

	// When
	todos, err := repo.GetAll()

	// Then
	assert.NoError(t, err)
	assert.Len(t, todos, 2)
	// Should be in DESC order (newest first)
	assert.Equal(t, "todo 2", todos[0].Text)
	assert.Equal(t, "todo 1", todos[1].Text)
}

func TestInMemoryTodoRepository_Update(t *testing.T) {
	// Given
	repo := repository.NewInMemoryTodoRepository()
	originalTodo := &model.Todo{Text: "original text"}
	created, err := repo.Create(originalTodo)
	require.NoError(t, err)

	// When
	created.Text = "updated text"
	result, err := repo.Update(created)

	// Then
	require.NoError(t, err)
	assert.Equal(t, created.ID, result.ID)
	assert.Equal(t, "updated text", result.Text)
	assert.Equal(t, created.CreatedAt, result.CreatedAt)
	assert.True(t, result.UpdatedAt.After(created.CreatedAt))

	// And: Repository should reflect the change
	fetched, err := repo.GetByID(created.ID)
	require.NoError(t, err)
	assert.Equal(t, "updated text", fetched.Text)
}

func TestInMemoryTodoRepository_Update_NotFound(t *testing.T) {
	// Given
	repo := repository.NewInMemoryTodoRepository()

	// When
	nonExistentTodo := &model.Todo{ID: 999, Text: "updated"}
	result, err := repo.Update(nonExistentTodo)

	// Then
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")
}

func TestInMemoryTodoRepository_Delete(t *testing.T) {
	// Given
	repo := repository.NewInMemoryTodoRepository()
	todo := &model.Todo{Text: "to be deleted"}
	created, err := repo.Create(todo)
	require.NoError(t, err)

	// When
	err = repo.Delete(created.ID)

	// Then
	require.NoError(t, err)

	// And: Todo should not be found
	_, err = repo.GetByID(created.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")

	// And: GetAll should return empty list
	todos, err := repo.GetAll()
	require.NoError(t, err)
	assert.Len(t, todos, 0)
}

func TestInMemoryTodoRepository_Delete_NotFound(t *testing.T) {
	// Given
	repo := repository.NewInMemoryTodoRepository()

	// When
	err := repo.Delete(999)

	// Then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
} 