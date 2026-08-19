package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	_ "github.com/oopbest/task-app/docs"
	"github.com/oopbest/task-app/handlers"
	"github.com/oopbest/task-app/metrics"
	"github.com/oopbest/task-app/middleware"
	"github.com/oopbest/task-app/repository"
	"github.com/oopbest/task-app/workers"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Task Management REST API (Go Zero to Hero)
// @version 1.0
// @description API สำหรับจัดการงานแบบ Multi-User พร้อม Redis Caching, Worker Pool และ Prometheus Monitoring
// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description ใส่ JWT Token ในรูปแบบ: Bearer <your_token>

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:mysecretpassword@localhost:5432/taskdb?sslmode=disable"
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	// 1. เชื่อมต่อ PostgreSQL
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer pingCancel()
	if err := db.PingContext(pingCtx); err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	log.Println("🐘 Connected to PostgreSQL database successfully!")

	// 2. เชื่อมต่อ Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	defer rdb.Close()

	redisCtx, redisCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer redisCancel()
	if err := rdb.Ping(redisCtx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("⚡ Connected to Redis cache successfully!")

	// 3. รัน Database Migrations อัตโนมัติ
	if err := repository.RunMigrations(db); err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}

	// 4. เริ่มต้น Background Worker Pool (3 Workers, Queue Size 100)
	workerPool := workers.NewWorkerPool(3, 100)
	workerPool.Start()

	// 5. สร้าง Repositories
	userRepo := repository.NewPostgresUserRepository(db)
	postgresTaskRepo := repository.NewPostgresTaskRepository(db)
	cachedTaskRepo := repository.NewCachedTaskRepository(postgresTaskRepo, rdb, 5*time.Minute)

	// 6. สร้าง Handlers
	authHandler := handlers.NewAuthHandler(userRepo)
	taskHandler := handlers.NewTaskHandler(cachedTaskRepo, workerPool)

	// 7. ตั้งค่า Gin Router & Observability Middlewares
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// 📊 สวม Structured Logger (slog) + Prometheus Metrics Middleware
	r.Use(middleware.StructuredLogger(), metrics.PrometheusMiddleware(), gin.Recovery())

	// 8. Prometheus Metrics Endpoint
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// 9. Swagger Documentation Endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 10. Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":   "ok",
			"database": "postgres",
			"cache":    "redis",
			"workers":  3,
			"swagger":  "/swagger/index.html",
			"metrics":  "/metrics",
			"time":     time.Now().Format(time.RFC3339),
		})
	})

	// 11. Public Routes (Auth)
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
	}

	// 12. Protected Routes (Tasks with JWT Auth)
	taskGroup := r.Group("/tasks", middleware.GinAuthMiddleware())
	{
		taskGroup.GET("", taskHandler.GetAllTasks)
		taskGroup.GET("/:id", taskHandler.GetTaskByID)
		taskGroup.POST("", taskHandler.CreateTask)
		taskGroup.PUT("/:id", taskHandler.UpdateTask)
		taskGroup.DELETE("/:id", taskHandler.DeleteTask)
	}

	// 13. ตั้งค่า HTTP Server พร้อม Graceful Shutdown
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		fmt.Println("==================================================")
		fmt.Printf("🚀 Gin Task API on http://localhost:%s\n", port)
		fmt.Printf("📖 Swagger UI on http://localhost:%s/swagger/index.html\n", port)
		fmt.Printf("📊 Prometheus Metrics on http://localhost:%s/metrics\n", port)
		fmt.Printf("📈 Grafana Dashboard on http://localhost:3000 (admin/admin)\n")
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

	workerPool.Stop()

	log.Println("✅ Server exited cleanly")
}
