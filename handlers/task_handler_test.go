package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/oopbest/task-app/models"
	"github.com/oopbest/task-app/repository"
)

func TestTaskHandler_GetAllTasks(t *testing.T) {
	repo := repository.NewMemoryTaskRepository()
	handler := NewTaskHandler(repo)

	// จำลอง GET /tasks request
	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rr := httptest.NewRecorder()

	handler.GetAllTasks(rr, req)

	// ตรวจสอบ HTTP Status Code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// ตรวจสอบ Response JSON
	var tasks []models.Task
	if err := json.Unmarshal(rr.Body.Bytes(), &tasks); err != nil {
		t.Fatalf("failed to decode response JSON: %v", err)
	}

	if len(tasks) < 2 {
		t.Errorf("expected at least 2 seed tasks, got %d", len(tasks))
	}
}

func TestTaskHandler_CreateTask(t *testing.T) {
	repo := repository.NewMemoryTaskRepository()
	handler := NewTaskHandler(repo)

	payload := []byte(`{"title":"New Task via Test","description":"Testing HTTP handler"}`)
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
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
}
