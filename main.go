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

	_ "github.com/lib/pq"
	"github.com/oopbest/task-app/handlers"
	"github.com/oopbest/task-app/middleware"
	"github.com/oopbest/task-app/repository"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:mysecretpassword@localhost:5432/taskdb?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// ทดสอบการเชื่อมต่อ Database
	pingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer pingCancel()
	if err := db.PingContext(pingCtx); err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	log.Println("🐘 Connected to PostgreSQL database successfully!")
	// 1. รัน Database Migrations อัตโนมัติ
	if err := repository.RunMigrations(db); err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}
	// 2. สร้าง Repositories
	userRepo := repository.NewPostgresUserRepository(db)
	taskRepo := repository.NewPostgresTaskRepository(db)
	// 3. สร้าง Handlers
	authHandler := handlers.NewAuthHandler(userRepo)
	taskHandler := handlers.NewTaskHandler(taskRepo)

	// 4. กำหนด Router
	mux := http.NewServeMux()

	// Public Endpoints
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status":   "ok",
			"database": "postgres",
			"time":     time.Now().Format(time.RFC3339),
		})
	})
	mux.HandleFunc("POST /auth/register", authHandler.Register)
	mux.HandleFunc("POST /auth/login", authHandler.Login)

	// Protected Endpoints (ต้องผ่าน AuthMiddleware)
	authMW := middleware.AuthMiddleware
	mux.Handle("GET /tasks", authMW(http.HandlerFunc(taskHandler.GetAllTasks)))
	mux.Handle("GET /tasks/{id}", authMW(http.HandlerFunc(taskHandler.GetTaskByID)))
	mux.Handle("POST /tasks", authMW(http.HandlerFunc(taskHandler.CreateTask)))
	mux.Handle("PUT /tasks/{id}", authMW(http.HandlerFunc(taskHandler.UpdateTask)))
	mux.Handle("DELETE /tasks/{id}", authMW(http.HandlerFunc(taskHandler.DeleteTask)))

	// 5. สวม Global Middleware (JSON + Logging)
	handlerWithMiddleware := middleware.Logging(middleware.JSONContentType(mux))

	// 6. ตั้งค่า HTTP Server
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handlerWithMiddleware,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		fmt.Println("==================================================")
		fmt.Printf("🚀 Task API (JWT Auth + Postgres) on http://localhost:%s\n", port)
		fmt.Println("==================================================")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

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
