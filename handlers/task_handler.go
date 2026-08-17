package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oopbest/task-app/models"
	"github.com/oopbest/task-app/repository"
	"github.com/oopbest/task-app/workers"
)

type TaskHandler struct {
	repo       repository.TaskRepository
	workerPool *workers.WorkerPool
}

func NewTaskHandler(repo repository.TaskRepository, workerPool *workers.WorkerPool) *TaskHandler {
	return &TaskHandler{
		repo:       repo,
		workerPool: workerPool,
	}
}

// GetAllTasks godoc
// @Summary ดึงรายการ Task ทั้งหมด (Pagination, Search, Filter & Sort)
// @Description ดึงรายการ Task ของ User ตัวเอง พร้อมตัวเลือกค้นหา กรอง และแบ่งหน้า
// @Tags Tasks
// @Security BearerAuth
// @Produce json
// @Param page query int false "หมายเลขหน้า (Default: 1)"
// @Param limit query int false "จำนวนรายการต่อหน้า (Default: 10, Max: 100)"
// @Param search query string false "ค้นหาจากชื่อเรื่องหรือรายละเอียด"
// @Param completed query bool false "กรองตามสถานะ (true = เสร็จแล้ว, false = ยังไม่เสร็จ)"
// @Param sort query string false "เรียงตามฟิลด์: id, title, created_at (Default: created_at)"
// @Param order query string false "ลำดับ: asc หรือ desc (Default: desc)"
// @Success 200 {object} models.PaginatedTasks
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /tasks [get]
func (h *TaskHandler) GetAllTasks(c *gin.Context) {
	userID := c.GetInt("user_id")

	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	search := c.Query("search")
	sortBy := c.Query("sort")
	order := c.Query("order")

	var completed *bool
	if completedStr := c.Query("completed"); completedStr != "" {
		if val, err := strconv.ParseBool(completedStr); err == nil {
			completed = &val
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

	result, err := h.repo.GetAll(userID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetTaskByID godoc
// @Summary ดึงข้อมูล Task ตาม ID
// @Description ดึงรายละเอียดของ Task เฉพาะของ User ตัวเองตาม ID
// @Tags Tasks
// @Security BearerAuth
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} models.Task
// @Failure 400 {object} map[string]string "Invalid ID format"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Task not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /tasks/{id} [get]
func (h *TaskHandler) GetTaskByID(c *gin.Context) {
	userID := c.GetInt("user_id")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID format"})
		return
	}

	task, err := h.repo.GetByID(id, userID)
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// CreateTask godoc
// @Summary สร้าง Task ใหม่
// @Description สร้างงานใหม่และผูกกับ User ID อัตโนมัติ พร้อมส่งแจ้งเตือนเข้า Worker Pool
// @Tags Tasks
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body models.CreateTaskInput true "ข้อมูลงานใหม่"
// @Success 201 {object} models.Task
// @Failure 400 {object} map[string]string "ข้อมูลไม่ถูกต้อง"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /tasks [post]
func (h *TaskHandler) CreateTask(c *gin.Context) {
	userID := c.GetInt("user_id")

	var input models.CreateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON request body"})
		return
	}

	if err := input.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task := h.repo.Create(input, userID)

	if h.workerPool != nil {
		h.workerPool.Enqueue(workers.Job{
			Type:   "TASK_CREATED",
			UserID: userID,
			TaskID: task.ID,
			Title:  task.Title,
		})
	}

	c.JSON(http.StatusCreated, task)
}

// UpdateTask godoc
// @Summary แก้ไขข้อมูล Task
// @Description แก้ไข Title, Description หรือสถานะ Completed ของ Task
// @Tags Tasks
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param input body models.UpdateTaskInput true "ข้อมูลที่ต้องการแก้ไข"
// @Success 200 {object} models.Task
// @Failure 400 {object} map[string]string "ข้อมูลไม่ถูกต้อง"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Task not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /tasks/{id} [put]
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	userID := c.GetInt("user_id")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID format"})
		return
	}

	var input models.UpdateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON request body"})
		return
	}

	if err := input.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.repo.Update(id, input, userID)
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if h.workerPool != nil {
		h.workerPool.Enqueue(workers.Job{
			Type:   "TASK_UPDATED",
			UserID: userID,
			TaskID: task.ID,
			Title:  task.Title,
		})
	}

	c.JSON(http.StatusOK, task)
}

// DeleteTask godoc
// @Summary ลบ Task
// @Description ลบ Task ของตัวเองออกจากระบบ
// @Tags Tasks
// @Security BearerAuth
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} map[string]string "Task deleted successfully"
// @Failure 400 {object} map[string]string "Invalid ID format"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Task not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /tasks/{id} [delete]
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	userID := c.GetInt("user_id")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID format"})
		return
	}

	if err := h.repo.Delete(id, userID); err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if h.workerPool != nil {
		h.workerPool.Enqueue(workers.Job{
			Type:   "TASK_DELETED",
			UserID: userID,
			TaskID: id,
			Title:  fmt.Sprintf("Task #%d", id),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task deleted successfully",
	})
}
