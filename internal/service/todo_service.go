package service

import (
	"fmt"
	"strings"
	
	"todo-app/internal/model"
	"todo-app/internal/repository"
)

// TodoService handles business logic for todos
type TodoService struct {
	repo repository.TodoRepository
}

// NewTodoService creates a new todo service
func NewTodoService(repo repository.TodoRepository) *TodoService {
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

// GetTodoByID returns a todo by its ID
func (s *TodoService) GetTodoByID(id int) (*model.Todo, error) {
	return s.repo.GetByID(id)
}

// GetAllTodos returns all todo items
func (s *TodoService) GetAllTodos() ([]*model.Todo, error) {
	return s.repo.GetAll()
}

// UpdateTodo updates an existing todo item
func (s *TodoService) UpdateTodo(id int, text string) (*model.Todo, error) {
	// Validate input
	if strings.TrimSpace(text) == "" {
		return nil, fmt.Errorf("text cannot be empty")
	}
	
	// First, check if todo exists
	existingTodo, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update the text
	existingTodo.Text = text
	return s.repo.Update(existingTodo)
}

// DeleteTodo deletes a todo item by ID
func (s *TodoService) DeleteTodo(id int) error {
	return s.repo.Delete(id)
} 