package repository

import (
	"fmt"
	"time"
	"todo-app/internal/model"
)

// TodoRepository defines the interface for todo data operations
type TodoRepository interface {
	Create(todo *model.Todo) (*model.Todo, error)
	GetByID(id int) (*model.Todo, error)
	GetAll() ([]*model.Todo, error)
	Update(todo *model.Todo) (*model.Todo, error)
	Delete(id int) error
}

// InMemoryTodoRepository implements TodoRepository using in-memory storage
type InMemoryTodoRepository struct {
	todos  []*model.Todo
	nextID int
}

// NewInMemoryTodoRepository creates a new in-memory todo repository
func NewInMemoryTodoRepository() *InMemoryTodoRepository {
	return &InMemoryTodoRepository{
		todos:  make([]*model.Todo, 0), // Boş bir slice oluştur
		nextID: 1,
	}
}

// Create adds a new todo to the repository
func (r *InMemoryTodoRepository) Create(todo *model.Todo) (*model.Todo, error) {
	now := time.Now()
	todo.ID = r.nextID
	todo.CreatedAt = now
	todo.UpdatedAt = now
	r.nextID++ // sınraki katıt için ID'yi artır
	r.todos = append(r.todos, todo)
	return todo, nil
}

// GetByID returns a todo by its ID
func (r *InMemoryTodoRepository) GetByID(id int) (*model.Todo, error) {
	for _, todo := range r.todos {
		if todo.ID == id {
			return todo, nil
		}
	}
	return nil, fmt.Errorf("todo with id %d not found", id)
}

// GetAll returns all todos from the repository, ordered by created_at DESC (newest first)
func (r *InMemoryTodoRepository) GetAll() ([]*model.Todo, error) {
	// Create a copy of the slice in reverse order (newest first)
	result := make([]*model.Todo, len(r.todos))
	for i, todo := range r.todos {
		result[len(r.todos)-1-i] = todo
	}
	return result, nil
}

// Update modifies an existing todo in the repository
func (r *InMemoryTodoRepository) Update(todo *model.Todo) (*model.Todo, error) {
	for i, existing := range r.todos {
		if existing.ID == todo.ID {
			todo.CreatedAt = existing.CreatedAt // Preserve creation time
			todo.UpdatedAt = time.Now()
			r.todos[i] = todo
			return todo, nil
		}
	}
	return nil, fmt.Errorf("todo with id %d not found", todo.ID)
}

// Delete removes a todo from the repository
func (r *InMemoryTodoRepository) Delete(id int) error {
	for i, todo := range r.todos {
		if todo.ID == id {
			// Remove element at index i
			r.todos = append(r.todos[:i], r.todos[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("todo with id %d not found", id)
} 