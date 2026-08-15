package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

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

// GetAllTasks จัดการ GET /tasks
func (h *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks := h.repo.GetAll()
	respondJSON(w, http.StatusOK, tasks)
}

// GetTaskByID จัดการ GET /tasks/{id}
func (h *TaskHandler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	// ดึงค่า id จาก URL path (เช่น /tasks/1 -> "1")
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid task ID format")
		return
	}

	task, err := h.repo.GetByID(id)
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

// CreateTask จัดการ POST /tasks
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var input models.CreateTaskInput

	// แปลง JSON Request Body เข้าสู่ struct
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid JSON request body")
		return
	}

	// Validate ข้อมูล
	if err := input.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	task := h.repo.Create(input)
	respondJSON(w, http.StatusCreated, task) // 201 Created
}

// UpdateTask จัดการ PUT /tasks/{id}
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
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

	task, err := h.repo.Update(id, input)
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

// DeleteTask จัดการ DELETE /tasks/{id}
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid task ID format")
		return
	}

	if err := h.repo.Delete(id); err != nil {
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
