package integration

import (
	"bytes"
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

// AcceptanceTest: User can add a todo item and see it in the list
func TestAddTodoItem_UserStory(t *testing.T) {
	// Given: Empty todo list
	server := setupTestServer(t)
	defer server.Close()

	// When: User adds "süt al" todo
	todoData := map[string]string{
		"text": "süt al",
	}
	jsonData, _ := json.Marshal(todoData)

	response, err := http.Post(server.URL+"/api/todos", "application/json", bytes.NewBuffer(jsonData))

	// Then: Todo should be created successfully
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, response.StatusCode)

	// And: Todo should appear in the list
	listResponse, err := http.Get(server.URL + "/api/todos")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, listResponse.StatusCode)

	var todos []map[string]interface{}
	json.NewDecoder(listResponse.Body).Decode(&todos)

	assert.Len(t, todos, 1)
	assert.Equal(t, "süt al", todos[0]["text"])
}

// AcceptanceTest: Error handling for non-existent todos
func TestErrorHandling_NotFound(t *testing.T) {
	// Given: Empty server
	server := setupTestServer(t)
	defer server.Close()

	// When: User tries to get non-existent todo
	getResp, err := http.Get(server.URL + "/api/todos/999")

	// Then: Should return 404
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, getResp.StatusCode)

	// When: User tries to update non-existent todo
	updateData := map[string]string{"text": "updated"}
	updateJSON, _ := json.Marshal(updateData)

	client := &http.Client{}
	req, _ := http.NewRequest("PUT", server.URL+"/api/todos/999", bytes.NewBuffer(updateJSON))
	req.Header.Set("Content-Type", "application/json")
	updateResp, err := client.Do(req)

	// Then: Should return 404
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, updateResp.StatusCode)

	// When: User tries to delete non-existent todo
	req, _ = http.NewRequest("DELETE", server.URL+"/api/todos/999", nil)
	deleteResp, err := client.Do(req)

	// Then: Should return 404
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, deleteResp.StatusCode)
}

// AcceptanceTest: Empty state is shown when no todos
func TestEmptyState_UserStory(t *testing.T) {
	// Given: Empty server
	server := setupTestServer(t)
	defer server.Close()

	// When: User gets all todos
	listResponse, err := http.Get(server.URL + "/api/todos")

	// Then: Empty array should be returned
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, listResponse.StatusCode)

	var todos []map[string]interface{}
	json.NewDecoder(listResponse.Body).Decode(&todos)
	assert.Len(t, todos, 0)
}

// AcceptanceTest: Multiple todos are displayed correctly
func TestMultipleTodos_UserStory(t *testing.T) {
	// Given: Server
	server := setupTestServer(t)
	defer server.Close()

	// When: User creates multiple todos
	todos := []string{"todo 1", "todo 2", "todo 3"}
	for _, todoText := range todos {
		todoData := map[string]string{"text": todoText}
		jsonData, _ := json.Marshal(todoData)
		_, err := http.Post(server.URL+"/api/todos", "application/json", bytes.NewBuffer(jsonData))
		require.NoError(t, err)
	}

	// Then: All todos should be listed (newest first - DESC order)
	listResponse, err := http.Get(server.URL + "/api/todos")
	require.NoError(t, err)

	var fetchedTodos []map[string]interface{}
	json.NewDecoder(listResponse.Body).Decode(&fetchedTodos)

	assert.Len(t, fetchedTodos, 3)
	// Should be in reverse order (newest first)
	assert.Equal(t, "todo 3", fetchedTodos[0]["text"])
	assert.Equal(t, "todo 2", fetchedTodos[1]["text"])
	assert.Equal(t, "todo 1", fetchedTodos[2]["text"])
}

// Setup test server with real handlers
func setupTestServer(t *testing.T) *httptest.Server {
	// Create dependencies
	repo, err := repository.NewSQLiteTodoRepository(":memory:")
	require.NoError(t, err)
	svc := service.NewTodoService(repo)
	h := handler.NewTodoHandler(svc)

	// Setup routes
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/todos", h.CreateTodo)
	mux.HandleFunc("GET /api/todos", h.GetAllTodos)

	return httptest.NewServer(mux)
}
