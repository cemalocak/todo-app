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
	// Get configuration from environment variables
	dbPath := getEnv("DB_PATH", "todos.db")
	port := getEnv("PORT", "8080")
	
	// Create dependencies
	repo, err := repository.NewSQLiteTodoRepository(dbPath)
	if err != nil {
		log.Fatalf("Failed to create SQLite repository: %v", err)
	}
	defer repo.Close()	// ← Program bitince database'i kapat
	
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