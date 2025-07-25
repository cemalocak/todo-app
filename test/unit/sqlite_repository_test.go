package unit

import (
	"fmt"
	"os"
	"testing"
	"time"

	"todo-app/internal/model"
	"todo-app/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestSQLiteTodoRepository_DatabasePersistence(t *testing.T) {
	// Given: Database with a todo
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	originalTodo := &model.Todo{Text: "persistent todo"}
	created, err := repo.Create(originalTodo)
	require.NoError(t, err)

	// When: Create new repository instance (simulates server restart)
	repo2, cleanup2 := setupTestDBWithSameFile(t, repo.DBPath())
	defer cleanup2()

	// Then: Todo should still exist
	todos, err := repo2.GetAll()
	require.NoError(t, err)
	assert.Len(t, todos, 1)
	assert.Equal(t, created.ID, todos[0].ID)
	assert.Equal(t, "persistent todo", todos[0].Text)
}

// Test helper functions
func setupTestDB(t *testing.T) (*repository.SQLiteTodoRepository, func()) {
	// Set test environment
	os.Setenv("ENV", "test")

	// Create temporary database file with unique name
	dbFile := fmt.Sprintf("test_todos_%d.db", time.Now().UnixNano())

	// Remove if exists
	os.Remove(dbFile)

	repo, err := repository.NewSQLiteTodoRepository(dbFile)
	require.NoError(t, err)

	cleanup := func() {
		repo.Close()
		os.Remove(dbFile)
		os.Unsetenv("ENV")
	}

	return repo, cleanup
}

func setupTestDBWithSameFile(t *testing.T, dbFile string) (*repository.SQLiteTodoRepository, func()) {
	repo, err := repository.NewSQLiteTodoRepository(dbFile)
	require.NoError(t, err)

	cleanup := func() {
		repo.Close()
	}

	return repo, cleanup
}
