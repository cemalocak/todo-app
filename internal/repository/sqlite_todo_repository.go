package repository

import (
	"database/sql"
	"embed"
	"fmt"
	"time"

	"todo-app/internal/model"

	_ "modernc.org/sqlite"
)

//go:embed database/schema.sql
var schemaFS embed.FS

// SQLiteTodoRepository implements TodoRepository using SQLite
type SQLiteTodoRepository struct {
	db     *sql.DB
	dbPath string
}

// NewSQLiteTodoRepository creates a new SQLite todo repository
func NewSQLiteTodoRepository(dbPath string) (*SQLiteTodoRepository, error) {
	db, err := sql.Open("sqlite", dbPath) // databse bağlantısını aç
	if err != nil {
		return nil, err
	}

	repo := &SQLiteTodoRepository{
		db:     db,
		dbPath: dbPath,
	}

	// Run migrations
	if err := repo.migrate(); err != nil {
		return nil, err
	}

	return repo, nil
}

// migrate runs database migrations
func (r *SQLiteTodoRepository) migrate() error {
	schema, err := schemaFS.ReadFile("database/schema.sql")
	if err != nil {
		return err
	}

	_, err = r.db.Exec(string(schema))
	return err
}

// Create adds a new todo to the database
func (r *SQLiteTodoRepository) Create(todo *model.Todo) (*model.Todo, error) {
	now := time.Now()

	query := `
		INSERT INTO todos (text, created_at, updated_at) 
		VALUES (?, ?, ?)
	`

	result, err := r.db.Exec(query, todo.Text, now, now)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId() // ← Son eklenen ID'yi al
	if err != nil {
		return nil, err
	}

	todo.ID = int(id)
	todo.CreatedAt = now
	todo.UpdatedAt = now

	return todo, nil
}

// GetByID returns a todo by its ID
func (r *SQLiteTodoRepository) GetByID(id int) (*model.Todo, error) {
	query := `
		SELECT id, text, created_at, updated_at 
		FROM todos 
		WHERE id = ?
	`

	todo := &model.Todo{}
	err := r.db.QueryRow(query, id).Scan(&todo.ID, &todo.Text, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("todo with id %d not found", id)
		}
		return nil, err
	}

	return todo, nil
}

// GetAll returns all todos from the database, ordered by created_at DESC
func (r *SQLiteTodoRepository) GetAll() ([]*model.Todo, error) {
	query := `
		SELECT id, text, created_at, updated_at 
		FROM todos 
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := make([]*model.Todo, 0) // Initialize as empty slice, not nil
	for rows.Next() {
		todo := &model.Todo{}
		err := rows.Scan(&todo.ID, &todo.Text, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, rows.Err()
}

// Update modifies an existing todo in the database
func (r *SQLiteTodoRepository) Update(todo *model.Todo) (*model.Todo, error) {
	now := time.Now()

	query := `
		UPDATE todos 
		SET text = ?, updated_at = ? 
		WHERE id = ?
	`

	result, err := r.db.Exec(query, todo.Text, now, todo.ID)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("todo with id %d not found", todo.ID)
	}

	todo.UpdatedAt = now
	return todo, nil
}

// Delete removes a todo from the database
func (r *SQLiteTodoRepository) Delete(id int) error {
	query := `DELETE FROM todos WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("todo with id %d not found", id)
	}

	return nil
}

// DBPath returns the database file path
func (r *SQLiteTodoRepository) DBPath() string {
	return r.dbPath
}

// Close closes the database connection
func (r *SQLiteTodoRepository) Close() error {
	return r.db.Close()
}
