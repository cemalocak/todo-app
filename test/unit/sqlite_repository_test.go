package unit

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"todo-app/internal/model"
	"todo-app/internal/repository"
)

func TestSQLiteTodoRepository_Create(t *testing.T) {
	// Given: Clean test database
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	todo := &model.Todo{Text: "test todo"}

	// When: Create todo
	result, err := repo.Create(todo)

	// Then: Todo should be created successfully
	require.NoError(t, err)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "test todo", result.Text)
	assert.False(t, result.CreatedAt.IsZero())
	assert.False(t, result.UpdatedAt.IsZero())
}

func TestSQLiteTodoRepository_GetByID(t *testing.T) {
	// Given: Database with a todo
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	originalTodo := &model.Todo{Text: "test todo"}
	created, err := repo.Create(originalTodo)
	require.NoError(t, err)

	// When: Get todo by ID
	result, err := repo.GetByID(created.ID)

	// Then: Todo should be found
	require.NoError(t, err)
	assert.Equal(t, created.ID, result.ID)
	assert.Equal(t, "test todo", result.Text)
	// SQLite precision: compare within 1 second tolerance
	assert.WithinDuration(t, created.CreatedAt, result.CreatedAt, time.Second)
	assert.WithinDuration(t, created.UpdatedAt, result.UpdatedAt, time.Second)
}

func TestSQLiteTodoRepository_GetByID_NotFound(t *testing.T) {
	// Given: Empty database
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	// When: Try to get non-existent todo
	result, err := repo.GetByID(999)

	// Then: Should return error
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")
}

func TestSQLiteTodoRepository_GetAll(t *testing.T) {
	// Given: Test database with todos
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	todo1 := &model.Todo{Text: "todo 1"}
	todo2 := &model.Todo{Text: "todo 2"}
	repo.Create(todo1)
	repo.Create(todo2)

	// When: Get all todos
	todos, err := repo.GetAll()

	// Then: All todos should be returned
	require.NoError(t, err)
	assert.Len(t, todos, 2)
	
	// And: Should be ordered by created_at DESC (newest first)
	assert.Equal(t, "todo 2", todos[0].Text) // Second todo (newer) should be first
	assert.Equal(t, "todo 1", todos[1].Text) // First todo (older) should be second
	assert.True(t, todos[0].CreatedAt.After(todos[1].CreatedAt) || 
				todos[0].CreatedAt.Equal(todos[1].CreatedAt))
}

func TestSQLiteTodoRepository_Update(t *testing.T) {
	// Given: Database with a todo
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	originalTodo := &model.Todo{Text: "original text"}
	created, err := repo.Create(originalTodo)
	require.NoError(t, err)

	// When: Update todo text
	created.Text = "updated text"
	result, err := repo.Update(created)

	// Then: Todo should be updated
	require.NoError(t, err)
	assert.Equal(t, created.ID, result.ID)
	assert.Equal(t, "updated text", result.Text)
	assert.WithinDuration(t, created.CreatedAt, result.CreatedAt, time.Second)
	assert.True(t, result.UpdatedAt.After(created.UpdatedAt) || 
				result.UpdatedAt.Equal(created.UpdatedAt))

	// And: Database should reflect the change
	fetched, err := repo.GetByID(created.ID)
	require.NoError(t, err)
	assert.Equal(t, "updated text", fetched.Text)
}

func TestSQLiteTodoRepository_Update_NotFound(t *testing.T) {
	// Given: Empty database
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	// When: Try to update non-existent todo
	nonExistentTodo := &model.Todo{ID: 999, Text: "updated"}
	result, err := repo.Update(nonExistentTodo)

	// Then: Should return error
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")
}

func TestSQLiteTodoRepository_Delete(t *testing.T) {
	// Given: Database with a todo
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	todo := &model.Todo{Text: "to be deleted"}
	created, err := repo.Create(todo)
	require.NoError(t, err)

	// When: Delete todo
	err = repo.Delete(created.ID)

	// Then: Todo should be deleted
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

func TestSQLiteTodoRepository_Delete_NotFound(t *testing.T) {
	// Given: Empty database
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	// When: Try to delete non-existent todo
	err := repo.Delete(999)

	// Then: Should return error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestSQLiteTodoRepository_DatabasePersistence(t *testing.T) {
	// Given: Database with a todo
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	originalTodo := &model.Todo{Text: "persistent todo"}
	created, err := repo.Create(originalTodo)
	require.NoError(t, err)

	// When: Create new repository instance (simulates server restart)
	repo2, cleanup2 := setupTestDBWithSameFile(t, repo.(*repository.SQLiteTodoRepository).DBPath())
	defer cleanup2()

	// Then: Todo should still exist
	todos, err := repo2.GetAll()
	require.NoError(t, err)
	assert.Len(t, todos, 1)
	assert.Equal(t, created.ID, todos[0].ID)
	assert.Equal(t, "persistent todo", todos[0].Text)
}

// Test helper functions
func setupTestDB(t *testing.T) (repository.TodoRepository, func()) {
	// Create temporary database file
	dbFile := "test_todos.db"
	
	// Remove if exists
	os.Remove(dbFile)
	
	repo, err := repository.NewSQLiteTodoRepository(dbFile)
	require.NoError(t, err)

	cleanup := func() {
		repo.Close()
		os.Remove(dbFile)
	}

	return repo, cleanup
}

func setupTestDBWithSameFile(t *testing.T, dbFile string) (repository.TodoRepository, func()) {
	repo, err := repository.NewSQLiteTodoRepository(dbFile)
	require.NoError(t, err)

	cleanup := func() {
		repo.Close()
	}

	return repo, cleanup
} 