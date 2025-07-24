package integration

import (
	"bytes"
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

// AcceptanceTest: User can add a todo item and see it in the list
func TestAddTodoItem_UserStory(t *testing.T) {
	// Given: Empty todo list
	server := setupTestServer()
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

// AcceptanceTest: User can get a specific todo by ID
func TestGetTodoByID_UserStory(t *testing.T) {
	// Given: Server with a todo
	server := setupTestServer()
	defer server.Close()

	// Create a todo first
	todoData := map[string]string{"text": "test todo"}
	jsonData, _ := json.Marshal(todoData)
	createResp, err := http.Post(server.URL+"/api/todos", "application/json", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	
	var createdTodo map[string]interface{}
	json.NewDecoder(createResp.Body).Decode(&createdTodo)
	todoID := int(createdTodo["id"].(float64))

	// When: User gets todo by ID
	getResp, err := http.Get(fmt.Sprintf("%s/api/todos/%d", server.URL, todoID))

	// Then: Todo should be returned
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, getResp.StatusCode)
	
	var fetchedTodo map[string]interface{} // key string olacak value herhangi bir veri tipi
	json.NewDecoder(getResp.Body).Decode(&fetchedTodo) // & işareti pointer olarak alır
	assert.Equal(t, "test todo", fetchedTodo["text"])
	assert.Equal(t, float64(todoID), fetchedTodo["id"])
}

// AcceptanceTest: User can update a todo
func TestUpdateTodo_UserStory(t *testing.T) {
	// Given: Server with a todo
	server := setupTestServer()
	defer server.Close()

	// Create a todo first
	todoData := map[string]string{"text": "original text"}
	jsonData, _ := json.Marshal(todoData)
	createResp, err := http.Post(server.URL+"/api/todos", "application/json", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	
	var createdTodo map[string]interface{}
	json.NewDecoder(createResp.Body).Decode(&createdTodo)
	todoID := int(createdTodo["id"].(float64))

	// When: User updates todo text
	updateData := map[string]string{"text": "updated text"}
	updateJSON, _ := json.Marshal(updateData)
	
	client := &http.Client{}
	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/api/todos/%d", server.URL, todoID), bytes.NewBuffer(updateJSON))
	req.Header.Set("Content-Type", "application/json")
	updateResp, err := client.Do(req)

	// Then: Todo should be updated
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, updateResp.StatusCode)
	
	var updatedTodo map[string]interface{}
	json.NewDecoder(updateResp.Body).Decode(&updatedTodo)
	assert.Equal(t, "updated text", updatedTodo["text"])
	assert.Equal(t, float64(todoID), updatedTodo["id"])

	// And: GET should return updated version
	getResp, err := http.Get(fmt.Sprintf("%s/api/todos/%d", server.URL, todoID))
	require.NoError(t, err)
	
	var fetchedTodo map[string]interface{}
	json.NewDecoder(getResp.Body).Decode(&fetchedTodo)
	assert.Equal(t, "updated text", fetchedTodo["text"])
}

// AcceptanceTest: User can delete a todo
func TestDeleteTodo_UserStory(t *testing.T) {
	// Given: Server with a todo
	server := setupTestServer()
	defer server.Close()

	// Create a todo first
	todoData := map[string]string{"text": "to be deleted"}
	jsonData, _ := json.Marshal(todoData)
	createResp, err := http.Post(server.URL+"/api/todos", "application/json", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	
	var createdTodo map[string]interface{}
	json.NewDecoder(createResp.Body).Decode(&createdTodo)
	todoID := int(createdTodo["id"].(float64))

	// When: User deletes todo
	client := &http.Client{}
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/api/todos/%d", server.URL, todoID), nil)
	deleteResp, err := client.Do(req)

	// Then: Todo should be deleted (204 No Content)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, deleteResp.StatusCode)

	// And: Todo should not be found
	getResp, err := http.Get(fmt.Sprintf("%s/api/todos/%d", server.URL, todoID))
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, getResp.StatusCode)

	// And: List should be empty
	listResp, err := http.Get(server.URL + "/api/todos")
	require.NoError(t, err)
	
	var todos []map[string]interface{}
	json.NewDecoder(listResp.Body).Decode(&todos)
	assert.Len(t, todos, 0)
}

// AcceptanceTest: Error handling for non-existent todos
func TestErrorHandling_NotFound(t *testing.T) {
	// Given: Empty server
	server := setupTestServer()
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

// AcceptanceTest: Error handling for invalid input
func TestErrorHandling_BadRequest(t *testing.T) {
	// Given: Server
	server := setupTestServer()
	defer server.Close()

	client := &http.Client{}

	// When: User tries invalid ID formats
	testCases := []string{
		"/api/todos/invalid",
		"/api/todos/abc",
	}

	for _, endpoint := range testCases {
		// GET with invalid ID
		getResp, err := http.Get(server.URL + endpoint)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, getResp.StatusCode, "GET %s should return 400", endpoint)

		// PUT with invalid ID
		updateData := map[string]string{"text": "test"}
		updateJSON, _ := json.Marshal(updateData)
		req, _ := http.NewRequest("PUT", server.URL+endpoint, bytes.NewBuffer(updateJSON))
		req.Header.Set("Content-Type", "application/json")
		updateResp, err := client.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, updateResp.StatusCode, "PUT %s should return 400", endpoint)

		// DELETE with invalid ID
		req, _ = http.NewRequest("DELETE", server.URL+endpoint, nil)
		deleteResp, err := client.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, deleteResp.StatusCode, "DELETE %s should return 400", endpoint)
	}
}

// AcceptanceTest: Empty state is shown when no todos
func TestEmptyState_UserStory(t *testing.T) {
	// Given: Empty server
	server := setupTestServer()
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
	server := setupTestServer()
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
func setupTestServer() *httptest.Server {
	// Create dependencies
	repo := repository.NewInMemoryTodoRepository()
	svc := service.NewTodoService(repo)
	h := handler.NewTodoHandler(svc)

	// Setup routes
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/todos", h.CreateTodo)
	mux.HandleFunc("GET /api/todos", h.GetAllTodos)
	mux.HandleFunc("GET /api/todos/", h.GetTodoByID)
	mux.HandleFunc("PUT /api/todos/", h.UpdateTodo)
	mux.HandleFunc("DELETE /api/todos/", h.DeleteTodo)

	return httptest.NewServer(mux)
} 