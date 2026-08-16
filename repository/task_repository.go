package repository

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/oopbest/task-app/models"
)

var ErrTaskNotFound = errors.New("task not found")

// TaskRepository กำหนด Contract สำหรับจัดการ Task
type TaskRepository interface {
	GetAll(userID int, filter models.TaskFilter) (models.PaginatedTasks, error)
	GetByID(id int, userID int) (models.Task, error)
	Create(input models.CreateTaskInput, userID int) models.Task
	Update(id int, input models.UpdateTaskInput, userID int) (models.Task, error)
	Delete(id int, userID int) error
}

type MemoryTaskRepository struct {
	mu     sync.RWMutex
	tasks  map[int]models.Task
	nextID int
}

func NewMemoryTaskRepository() *MemoryTaskRepository {
	repo := &MemoryTaskRepository{
		tasks:  make(map[int]models.Task),
		nextID: 1,
	}

	repo.Create(models.CreateTaskInput{Title: "Task 1", Description: "Description 1"}, 1)
	repo.Create(models.CreateTaskInput{Title: "Task 2", Description: "Description 2"}, 1)

	return repo
}

func (r *MemoryTaskRepository) Create(input models.CreateTaskInput, userID int) models.Task {
	r.mu.Lock()
	defer r.mu.Unlock()
	task := models.Task{
		ID:          r.nextID,
		Title:       input.Title,
		Description: input.Description,
		Completed:   false,
		UserID:      userID,
		CreatedAt:   time.Now(),
	}
	r.tasks[task.ID] = task
	r.nextID++
	return task
}

func (r *MemoryTaskRepository) GetAll(userID int, filter models.TaskFilter) (models.PaginatedTasks, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	filter.Sanitize()

	matched := make([]models.Task, 0)
	for _, task := range r.tasks {
		if task.UserID != userID {
			continue
		}
		if filter.Search != "" && !strings.Contains(strings.ToLower(task.Title), strings.ToLower(filter.Search)) {
			continue
		}
		if filter.Completed != nil && task.Completed != *filter.Completed {
			continue
		}
		matched = append(matched, task)
	}

	return models.PaginatedTasks{
		Data:       matched,
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalItems: len(matched),
		TotalPages: 1,
	}, nil
}

func (r *MemoryTaskRepository) GetByID(id int, userID int) (models.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	task, exists := r.tasks[id]
	if !exists || task.UserID != userID {
		return models.Task{}, ErrTaskNotFound
	}
	return task, nil
}

func (r *MemoryTaskRepository) Update(id int, input models.UpdateTaskInput, userID int) (models.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	task, exists := r.tasks[id]
	if !exists || task.UserID != userID {
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

func (r *MemoryTaskRepository) Delete(id int, userID int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	task, exists := r.tasks[id]
	if !exists || task.UserID != userID {
		return ErrTaskNotFound
	}
	delete(r.tasks, id)
	return nil
}
