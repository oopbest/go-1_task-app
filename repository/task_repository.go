package repository

import (
	"errors"
	"sync"
	"time"

	"github.com/oopbest/task-app/models"
)

// ErrTaskNotFound กำหนด Error เฉพาะกรณีหา Task ไม่พบ
var ErrTaskNotFound = errors.New("task not found")

// TaskRepository กำหนด Contract ว่าระบบจัดการ Task ต้องทำอะไรได้บ้าง
type TaskRepository interface {
	GetAll() []models.Task
	GetByID(id int) (models.Task, error)
	Create(input models.CreateTaskInput) models.Task
	Update(id int, input models.UpdateTaskInput) (models.Task, error)
	Delete(id int) error
}

// MemoryTaskRepository เก็บข้อมูลใน Memory โดยใช้ Go map
type MemoryTaskRepository struct {
	mu     sync.RWMutex // ป้องกันหลาย Goroutine แก้ข้อมูลพร้อมกัน
	tasks  map[int]models.Task
	nextID int
}

// NewMemoryTaskRepository Constructor สำหรับสร้าง Repository พร้อม Seed Data เริ่มต้น
func NewMemoryTaskRepository() *MemoryTaskRepository {
	repo := &MemoryTaskRepository{
		tasks:  make(map[int]models.Task),
		nextID: 1,
	}

	repo.Create(models.CreateTaskInput{
		Title:       "Task 1",
		Description: "Description 1",
	})

	repo.Create(models.CreateTaskInput{
		Title:       "Task 2",
		Description: "Description 2",
	})

	return repo
}

// Create สร้าง Task ใหม่
func (r *MemoryTaskRepository) Create(input models.CreateTaskInput) models.Task {
	r.mu.Lock()
	defer r.mu.Unlock()
	task := models.Task{
		ID:          r.nextID,
		Title:       input.Title,
		Description: input.Description,
		Completed:   false,
		CreatedAt:   time.Now(),
	}
	r.tasks[task.ID] = task
	r.nextID++
	return task
}

// GetAll ดึงรายการ Task ทั้งหมด
func (r *MemoryTaskRepository) GetAll() []models.Task {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]models.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		result = append(result, task)
	}
	return result
}

// GetByID ดึงข้อมูล Task ตาม ID
func (r *MemoryTaskRepository) GetByID(id int) (models.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	task, exists := r.tasks[id]
	if !exists {
		return models.Task{}, ErrTaskNotFound
	}
	return task, nil
}

// Update แก้ไข Task ตาม ID
func (r *MemoryTaskRepository) Update(id int, input models.UpdateTaskInput) (models.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	task, exists := r.tasks[id]
	if !exists {
		return models.Task{}, ErrTaskNotFound
	}

	if input.Title != nil {
		task.Title = *input.Title
	}

	if input.Description != nil {
		task.Description = *input.Description
	}

	if input.Completed != nil {
		task.Completed = *input.Completed
	}
	r.tasks[id] = task
	return task, nil
}

// Delete ลบ Task ตาม ID ที่ระบุ
func (r *MemoryTaskRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[id]; !exists {
		return ErrTaskNotFound
	}
	delete(r.tasks, id)
	return nil
}
