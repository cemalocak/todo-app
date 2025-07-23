package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"todo-app/internal/handler"
	"todo-app/internal/repository"
	"todo-app/internal/service"
)

func main() {
	// Get configuration from environment variables with environment-specific defaults
	env := getEnv("ENV", "development")
	port := getEnv("PORT", "8080")

	// Set database path based on environment
	var defaultDBPath string
	switch env {
	case "production":
		defaultDBPath = "/data/todos.db"
	case "test":
		defaultDBPath = ":memory:" // In-memory database for tests
	default: // development
		defaultDBPath = "todos_dev.db"
	}

	dbPath := getEnv("DB_PATH", defaultDBPath)

	// Create dependencies
	repo, err := repository.NewSQLiteTodoRepository(dbPath)
	if err != nil {
		log.Fatalf("Failed to create SQLite repository: %v", err)
	}
	defer repo.Close() // ← Program bitince database'i kapat

	svc := service.NewTodoService(repo)
	h := handler.NewTodoHandler(svc)

	// Setup routes
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("POST /api/todos", h.CreateTodo)
	mux.HandleFunc("GET /api/todos", h.GetAllTodos)
	mux.HandleFunc("GET /api/todos/", h.GetTodoByID)
	mux.HandleFunc("PUT /api/todos/", h.UpdateTodo)
	mux.HandleFunc("DELETE /api/todos/", h.DeleteTodo)

	// Serve static files (frontend)
	mux.Handle("/", http.FileServer(http.Dir("web/")))

	// Start server
	serverPort := ":" + port
	fmt.Printf("🚀 Server starting on http://localhost%s\n", serverPort)
	fmt.Printf("📝 API: http://localhost%s/api/todos\n", serverPort)
	fmt.Printf("🌐 Frontend: http://localhost%s\n", serverPort)
	fmt.Printf("💾 Database: %s\n", dbPath)

	log.Fatal(http.ListenAndServe(serverPort, mux))
}

// getEnv gets environment variable with default fallback
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
