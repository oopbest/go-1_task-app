package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/oopbest/task-app/middleware"
	"github.com/oopbest/task-app/models"
	"github.com/oopbest/task-app/repository"
)

// TaskHandler จัดการ HTTP Request ทั้งหมดที่เกี่ยวกับ Task
type TaskHandler struct {
	repo repository.TaskRepository
}

// NewTaskHandler ฟังก์ชัน Constructor สำหรับสร้าง TaskHandler
func NewTaskHandler(repo repository.TaskRepository) *TaskHandler {
	return &TaskHandler{repo: repo}
}

// Helper: ตอบกลับข้อมูล JSON และ Status Code
func respondJSON(w http.ResponseWriter, status int, payload any) {
	w.WriteHeader(status)
	if payload != nil {
		_ = json.NewEncoder(w).Encode(payload)
	}
}

// Helper: ตอบกลับข้อความ Error ในรูปแบบ JSON
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

// GetAllTasks ดึงรายการ Task ตามเงื่อนไข Query Parameters (Search, Filter, Sort, Pagination)
func (h *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	// 1. อ่าน Query Parameters จาก URL
	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	limit, _ := strconv.Atoi(query.Get("limit"))
	search := query.Get("search")
	sortBy := query.Get("sort")
	order := query.Get("order")
	var completed *bool
	if completedStr := query.Get("completed"); completedStr != "" {
		if c, err := strconv.ParseBool(completedStr); err == nil {
			completed = &c
		}
	}
	filter := models.TaskFilter{
		Page:      page,
		Limit:     limit,
		Search:    search,
		Completed: completed,
		SortBy:    sortBy,
		Order:     order,
	}
	// 2. เรียก Repository
	result, err := h.repo.GetAll(userID, filter)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch tasks")
		return
	}
	respondJSON(w, http.StatusOK, result)
}

// GetTaskByID ดึง Task ตาม ID (เฉพาะของ User ตัวเอง)
func (h *TaskHandler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid task ID format")
		return
	}
	task, err := h.repo.GetByID(id, userID)
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			respondError(w, http.StatusNotFound, "Task not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	respondJSON(w, http.StatusOK, task)
}

// CreateTask สร้าง Task โดยผูกกับ UserID
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	var input models.CreateTaskInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid JSON request body")
		return
	}
	if err := input.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	task := h.repo.Create(input, userID)
	respondJSON(w, http.StatusCreated, task)
}

// UpdateTask แก้ไข Task (เฉพาะของตัวเอง)
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid task ID format")
		return
	}
	var input models.UpdateTaskInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid JSON request body")
		return
	}
	if err := input.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	task, err := h.repo.Update(id, input, userID)
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			respondError(w, http.StatusNotFound, "Task not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	respondJSON(w, http.StatusOK, task)
}

// DeleteTask ลบ Task (เฉพาะของตัวเอง)
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid task ID format")
		return
	}
	if err := h.repo.Delete(id, userID); err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			respondError(w, http.StatusNotFound, "Task not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{
		"message": "Task deleted successfully",
	})
}
