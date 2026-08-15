package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/oopbest/task-app/handlers"
	"github.com/oopbest/task-app/middleware"
	"github.com/oopbest/task-app/repository"
)

func main() {
	// 1. อ่าน Port จาก Environment Variable (ถ้าไม่มีให้ใช้ค่าเริ่มต้น 8080)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 2. สร้าง In-Memory Repository และ Handlers (Dependency Injection)
	repo := repository.NewMemoryTaskRepository()
	taskHandler := handlers.NewTaskHandler(repo)

	// 3. กำหนด Router
	mux := http.NewServeMux()

	// Health Check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Task RESTful Endpoints
	mux.HandleFunc("GET /tasks", taskHandler.GetAllTasks)
	mux.HandleFunc("GET /tasks/{id}", taskHandler.GetTaskByID)
	mux.HandleFunc("POST /tasks", taskHandler.CreateTask)
	mux.HandleFunc("PUT /tasks/{id}", taskHandler.UpdateTask)
	mux.HandleFunc("DELETE /tasks/{id}", taskHandler.DeleteTask)

	// 4. สวม Middleware
	handlerWithMiddleware := middleware.Logging(middleware.JSONContentType(mux))

	// 5. สร้าง HTTP Server
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handlerWithMiddleware,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// 6. รัน Server ใน Goroutine ย่อย (Background) เพื่อไม่ให้บล็อกการรอรับสัญญาณ Shutdown
	go func() {
		fmt.Println("==================================================")
		fmt.Printf("🚀 Server is running on http://localhost:%s\n", port)
		fmt.Println("==================================================")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// 7. ดักจับสัญญาณ Shutdown (Ctrl+C หรือ SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// บรรทัดนี้จะรอ (Block) จนกว่าจะมีสัญญาณปิดส่งเข้ามา
	<-quit
	log.Println("🛑 Shutdown signal received, shutting down gracefully...")

	// ให้เวลา 5 วินาทีในการเคลียร์ Request เก่าที่ยังทำงานค้างอยู่
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited cleanly")
}
