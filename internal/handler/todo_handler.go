package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"todo-app/internal/service"
)

// TodoHandler handles HTTP requests for todos
type TodoHandler struct {
	service *service.TodoService
}

// NewTodoHandler creates a new todo handler
func NewTodoHandler(service *service.TodoService) *TodoHandler {
	return &TodoHandler{
		service: service,
	}
}

// CreateTodo handles POST /api/todos
func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	// Validate Content-Type
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
		return
	}

	// json'u struct yapısına çevir
	var request struct {
		Text string `json:"text"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate empty text
	if strings.TrimSpace(request.Text) == "" {
		http.Error(w, "Text cannot be empty", http.StatusBadRequest)
		return
	}

	todo, err := h.service.CreateTodo(request.Text)
	if err != nil {
		if strings.Contains(err.Error(), "empty") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to create todo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	json.NewEncoder(w).Encode(todo)   // struct'ı json'a çevir
}

// GetTodoByID handles GET /api/todos/{id}
func (h *TodoHandler) GetTodoByID(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/todos/")
	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	todo, err := h.service.GetTodoByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get todo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

// GetAllTodos handles GET /api/todos
func (h *TodoHandler) GetAllTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := h.service.GetAllTodos()
	if err != nil {
		http.Error(w, "Failed to get todos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

// UpdateTodo handles PUT /api/todos/{id}
func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/todos/")
	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	var request struct {
		Text string `json:"text"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate empty text
	if strings.TrimSpace(request.Text) == "" {
		http.Error(w, "Text cannot be empty", http.StatusBadRequest)
		return
	}

	todo, err := h.service.UpdateTodo(id, request.Text)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update todo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

// DeleteTodo handles DELETE /api/todos/{id}
func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/todos/")
	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteTodo(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete todo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content for successful deletion
}

// TruncateTodos handles removing all todos (for testing only)
func (h *TodoHandler) TruncateTodos(w http.ResponseWriter, r *http.Request) {
	if err := h.service.TruncateTodos(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
