package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq" // Postgres Driver
	"github.com/oopbest/task-app/handlers"
	"github.com/oopbest/task-app/middleware"
	"github.com/oopbest/task-app/repository"
)

func main() {
	// 1. อ่านค่า Config จาก Environment Variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Connection String ไปยัง PostgreSQL ใน Docker
		dbURL = "postgres://postgres:mysecretpassword@localhost:5432/taskdb?sslmode=disable"
	}

	// 2. เชื่อมต่อฐานข้อมูล PostgreSQL
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	defer db.Close()

	// ตั้งค่า Connection Pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// ทดสอบการเชื่อมต่อ
	pingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer pingCancel()
	if err := db.PingContext(pingCtx); err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	log.Println("🐘 Connected to PostgreSQL database successfully!")

	// 3. สร้าง PostgreSQL Repository
	repo, err := repository.NewPostgresTaskRepository(db)
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}

	// 4. สร้าง Handlers (ส่ง Postgres Repo เข้าไป)
	taskHandler := handlers.NewTaskHandler(repo)

	// 5. กำหนด Router
	mux := http.NewServeMux()

	// Health Check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status":   "ok",
			"database": "postgres",
			"time":     time.Now().Format(time.RFC3339),
		})
	})

	// Task RESTful Endpoints
	mux.HandleFunc("GET /tasks", taskHandler.GetAllTasks)
	mux.HandleFunc("GET /tasks/{id}", taskHandler.GetTaskByID)
	mux.HandleFunc("POST /tasks", taskHandler.CreateTask)
	mux.HandleFunc("PUT /tasks/{id}", taskHandler.UpdateTask)
	mux.HandleFunc("DELETE /tasks/{id}", taskHandler.DeleteTask)

	// 6. สวม Middleware
	handlerWithMiddleware := middleware.Logging(middleware.JSONContentType(mux))

	// 7. ตั้งค่า HTTP Server
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handlerWithMiddleware,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// 8. รัน Server ใน Background Goroutine
	go func() {
		fmt.Println("==================================================")
		fmt.Printf("🚀 Task API (PostgreSQL) is running on http://localhost:%s\n", port)
		fmt.Println("==================================================")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// 9. ดักจับสัญญาณ Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Println("🛑 Shutdown signal received, shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited cleanly")
}
