package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/oopbest/task-app/handlers"
	"github.com/oopbest/task-app/middleware"
	"github.com/oopbest/task-app/repository"
)

func main() {
	// 1. สร้าง In-Memory Repository
	repo := repository.NewMemoryTaskRepository()

	// 2. สร้าง TaskHandler โดยส่ง Repository เข้าไป (Dependency Injection)
	taskHandler := handlers.NewTaskHandler(repo)

	// 3. กำหนด Router ด้วย http.NewServeMux
	mux := http.NewServeMux()

	// Health Check Endpoint
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Task RESTful Endpoints (Go 1.22+ Syntax)
	mux.HandleFunc("GET /tasks", taskHandler.GetAllTasks)
	mux.HandleFunc("GET /tasks/{id}", taskHandler.GetTaskByID)
	mux.HandleFunc("POST /tasks", taskHandler.CreateTask)
	mux.HandleFunc("PUT /tasks/{id}", taskHandler.UpdateTask)
	mux.HandleFunc("DELETE /tasks/{id}", taskHandler.DeleteTask)

	// 4. สวม Middleware (JSON Header + Logging)
	handlerWithMiddleware := middleware.Logging(middleware.JSONContentType(mux))

	// 5. ตั้งค่า HTTP Server
	port := ":8080"
	server := &http.Server{
		Addr:         port,
		Handler:      handlerWithMiddleware,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Println("==================================================")
	fmt.Printf("🚀 Task Management API is running on http://localhost%s\n", port)
	fmt.Println("==================================================")
	fmt.Println("📌 Available Endpoints:")
	fmt.Println("  - GET    /health      : Health check")
	fmt.Println("  - GET    /tasks       : Get all tasks")
	fmt.Println("  - GET    /tasks/{id}  : Get task by ID")
	fmt.Println("  - POST   /tasks       : Create a new task")
	fmt.Println("  - PUT    /tasks/{id}  : Update a task")
	fmt.Println("  - DELETE /tasks/{id}  : Delete a task")
	fmt.Println("==================================================")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}
