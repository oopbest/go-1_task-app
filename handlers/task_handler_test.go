package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/oopbest/task-app/models"
	"github.com/oopbest/task-app/repository"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupTestRouter(repo repository.TaskRepository) *gin.Engine {
	r := gin.New()
	handler := NewTaskHandler(repo, nil)

	// Mock Auth middleware setting user_id = 1
	r.Use(func(c *gin.Context) {
		c.Set("user_id", 1)
		c.Next()
	})

	r.GET("/tasks", handler.GetAllTasks)
	r.POST("/tasks", handler.CreateTask)
	return r
}

func TestTaskHandler_GetAllTasks(t *testing.T) {
	repo := repository.NewMemoryTaskRepository()
	r := setupTestRouter(repo)

	req := httptest.NewRequest(http.MethodGet, "/tasks?page=1&limit=10", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response models.PaginatedTasks
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response JSON: %v", err)
	}

	if len(response.Data) < 2 {
		t.Errorf("expected at least 2 seed tasks, got %d", len(response.Data))
	}
	if response.Page != 1 || response.Limit != 10 {
		t.Errorf("expected page 1, limit 10; got page %d, limit %d", response.Page, response.Limit)
	}
}

func TestTaskHandler_CreateTask(t *testing.T) {
	repo := repository.NewMemoryTaskRepository()
	r := setupTestRouter(repo)

	payload := []byte(`{"title":"New Task via Test","description":"Testing HTTP handler"}`)
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var created models.Task
	if err := json.Unmarshal(rr.Body.Bytes(), &created); err != nil {
		t.Fatalf("failed to decode response JSON: %v", err)
	}

	if created.Title != "New Task via Test" {
		t.Errorf("expected title 'New Task via Test', got %q", created.Title)
	}
	if created.UserID != 1 {
		t.Errorf("expected UserID 1, got %d", created.UserID)
	}
}
