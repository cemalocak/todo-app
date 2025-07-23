package performance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"todo-app/internal/handler"
	"todo-app/internal/repository"
	"todo-app/internal/service"
)

// TestConcurrentRequests tests handling multiple concurrent requests
func TestConcurrentRequests(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance tests in short mode")
	}

	// Setup
	repo := repository.NewInMemoryTodoRepository()
	svc := service.NewTodoService(repo)
	h := handler.NewTodoHandler(svc)

	const numRequests = 100
	const numWorkers = 10

	// Test concurrent todo creation
	t.Run("Concurrent Create", func(t *testing.T) {
		var wg sync.WaitGroup
		responses := make(chan *httptest.ResponseRecorder, numRequests)
		
		start := time.Now()
		
		for i := 0; i < numRequests; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				
				requestBody := map[string]string{"text": fmt.Sprintf("todo %d", id)}
				jsonBody, _ := json.Marshal(requestBody)
				req := httptest.NewRequest("POST", "/api/todos", bytes.NewBuffer(jsonBody))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				
				h.CreateTodo(rec, req)
				responses <- rec
			}(i)
		}
		
		wg.Wait()
		close(responses)
		
		duration := time.Since(start)
		t.Logf("Created %d todos in %v (%.2f req/sec)", numRequests, duration, float64(numRequests)/duration.Seconds())
		
		// Verify all requests succeeded
		successCount := 0
		for rec := range responses {
			if rec.Code == http.StatusCreated {
				successCount++
			}
		}
		
		assert.Equal(t, numRequests, successCount, "All requests should succeed")
	})

	// Test concurrent reads
	t.Run("Concurrent Read", func(t *testing.T) {
		// Create some todos first
		for i := 0; i < 10; i++ {
			_, err := svc.CreateTodo(fmt.Sprintf("read test todo %d", i))
			assert.NoError(t, err)
		}

		var wg sync.WaitGroup
		responses := make(chan *httptest.ResponseRecorder, numRequests)
		
		start := time.Now()
		
		for i := 0; i < numRequests; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				
				req := httptest.NewRequest("GET", "/api/todos", nil)
				rec := httptest.NewRecorder()
				
				h.GetAllTodos(rec, req)
				responses <- rec
			}()
		}
		
		wg.Wait()
		close(responses)
		
		duration := time.Since(start)
		t.Logf("Completed %d read requests in %v (%.2f req/sec)", numRequests, duration, float64(numRequests)/duration.Seconds())
		
		// Verify all requests succeeded
		successCount := 0
		for rec := range responses {
			if rec.Code == http.StatusOK {
				successCount++
			}
		}
		
		assert.Equal(t, numRequests, successCount, "All read requests should succeed")
	})
}

// TestResponseTime tests response time under normal load
func TestResponseTime(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance tests in short mode")
	}

	repo := repository.NewInMemoryTodoRepository()
	svc := service.NewTodoService(repo)
	h := handler.NewTodoHandler(svc)

	// Test create response time
	t.Run("Create Response Time", func(t *testing.T) {
		times := make([]time.Duration, 100)
		
		for i := 0; i < 100; i++ {
			requestBody := map[string]string{"text": fmt.Sprintf("timing test %d", i)}
			jsonBody, _ := json.Marshal(requestBody)
			req := httptest.NewRequest("POST", "/api/todos", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			
			start := time.Now()
			h.CreateTodo(rec, req)
			times[i] = time.Since(start)
			
			assert.Equal(t, http.StatusCreated, rec.Code)
		}
		
		// Calculate statistics
		var total time.Duration
		var max time.Duration
		min := times[0]
		
		for _, t := range times {
			total += t
			if t > max {
				max = t
			}
			if t < min {
				min = t
			}
		}
		
		avg := total / time.Duration(len(times))
		
		t.Logf("Create Todo Response Times - Avg: %v, Min: %v, Max: %v", avg, min, max)
		
		// Assert reasonable response times (these thresholds may need adjustment)
		assert.Less(t, avg, 10*time.Millisecond, "Average response time should be reasonable")
		assert.Less(t, max, 50*time.Millisecond, "Max response time should be reasonable")
	})
}

// TestMemoryUsage tests basic memory usage patterns
func TestMemoryUsage(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance tests in short mode")
	}

	repo := repository.NewInMemoryTodoRepository()
	svc := service.NewTodoService(repo)
	h := handler.NewTodoHandler(svc)

	// Create many todos and measure growth
	t.Run("Memory Growth", func(t *testing.T) {
		const numTodos = 1000
		
		start := time.Now()
		
		for i := 0; i < numTodos; i++ {
			requestBody := map[string]string{"text": fmt.Sprintf("memory test todo %d", i)}
			jsonBody, _ := json.Marshal(requestBody)
			req := httptest.NewRequest("POST", "/api/todos", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			
			h.CreateTodo(rec, req)
			assert.Equal(t, http.StatusCreated, rec.Code)
		}
		
		duration := time.Since(start)
		t.Logf("Created %d todos in %v", numTodos, duration)
		
		// Verify we can still read all todos efficiently
		req := httptest.NewRequest("GET", "/api/todos", nil)
		rec := httptest.NewRecorder()
		
		readStart := time.Now()
		h.GetAllTodos(rec, req)
		readDuration := time.Since(readStart)
		
		assert.Equal(t, http.StatusOK, rec.Code)
		t.Logf("Read %d todos in %v", numTodos, readDuration)
		
		// Verify response contains all todos
		var todos []map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &todos)
		assert.NoError(t, err)
		assert.Len(t, todos, numTodos)
	})
} 