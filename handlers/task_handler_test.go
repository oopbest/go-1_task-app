package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/oopbest/task-app/middleware"
	"github.com/oopbest/task-app/models"
	"github.com/oopbest/task-app/repository"
)

func TestTaskHandler_GetAllTasks(t *testing.T) {
	repo := repository.NewMemoryTaskRepository()
	handler := NewTaskHandler(repo)

	// จำลอง Request พร้อม Context ที่มี UserID = 1
	req := httptest.NewRequest(http.MethodGet, "/tasks?page=1&limit=10", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDContextKey, 1)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.GetAllTasks(rr, req)

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
	handler := NewTaskHandler(repo)

	payload := []byte(`{"title":"New Task via Test","description":"Testing HTTP handler"}`)
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	ctx := context.WithValue(req.Context(), middleware.UserIDContextKey, 1)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.CreateTask(rr, req)

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
