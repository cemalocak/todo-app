package service

import (
	"fmt"
	"strings"

	"todo-app/internal/model"
	"todo-app/internal/repository"
)

// TodoService handles business logic for todos
type TodoService struct {
	repo *repository.SQLiteTodoRepository
}

// NewTodoService creates a new todo service
func NewTodoService(repo *repository.SQLiteTodoRepository) *TodoService {
	return &TodoService{
		repo: repo,
	}
}

// CreateTodo creates a new todo item
func (s *TodoService) CreateTodo(text string) (*model.Todo, error) {
	// Validate input
	if strings.TrimSpace(text) == "" {
		return nil, fmt.Errorf("text cannot be empty")
	}

	todo := &model.Todo{
		Text: text,
	}
	return s.repo.Create(todo)
}

// GetAllTodos returns all todo items
func (s *TodoService) GetAllTodos() ([]*model.Todo, error) {
	return s.repo.GetAll()
}

// TruncateTodos removes all todos (for testing only)
func (s *TodoService) TruncateTodos() error {
	return s.repo.Truncate()
}
