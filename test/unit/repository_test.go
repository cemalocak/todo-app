package unit

import (
	"testing"

	"todo-app/internal/model"
	"todo-app/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryTodoRepository_Create(t *testing.T) {
	// Given
	repo, err := repository.NewSQLiteTodoRepository(":memory:")
	require.NoError(t, err)
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

func TestInMemoryTodoRepository_GetAll(t *testing.T) {
	// Given
	repo, err := repository.NewSQLiteTodoRepository(":memory:")
	require.NoError(t, err)
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
