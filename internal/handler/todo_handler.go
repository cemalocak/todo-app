package handler

import (
	"encoding/json"
	"net/http"
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

// TruncateTodos handles removing all todos (for testing only)
func (h *TodoHandler) TruncateTodos(w http.ResponseWriter, r *http.Request) {
	if err := h.service.TruncateTodos(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
