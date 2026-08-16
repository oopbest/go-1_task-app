package repository

import (
	"testing"

	"github.com/oopbest/task-app/models"
)

func TestMemoryTaskRepository_CRUD(t *testing.T) {
	repo := NewMemoryTaskRepository()
	testUserID := 1

	// 1. ทดสอบ Create
	created := repo.Create(models.CreateTaskInput{
		Title:       "Test Task",
		Description: "Testing Description",
	}, testUserID)

	if created.ID == 0 {
		t.Fatalf("expected valid task ID, got 0")
	}
	if created.Title != "Test Task" {
		t.Errorf("expected title %q, got %q", "Test Task", created.Title)
	}

	// 2. ทดสอบ GetByID
	fetched, err := repo.GetByID(created.ID, testUserID)
	if err != nil {
		t.Fatalf("unexpected error getting task by ID: %v", err)
	}
	if fetched.ID != created.ID {
		t.Errorf("expected ID %d, got %d", created.ID, fetched.ID)
	}

	// 3. ทดสอบ Update
	newTitle := "Updated Task Title"
	isDone := true
	updated, err := repo.Update(created.ID, models.UpdateTaskInput{
		Title:     &newTitle,
		Completed: &isDone,
	}, testUserID)
	if err != nil {
		t.Fatalf("unexpected error updating task: %v", err)
	}
	if updated.Title != newTitle || updated.Completed != true {
		t.Errorf("expected title %q and completed %v, got title %q and completed %v", newTitle, true, updated.Title, updated.Completed)
	}

	// 4. ทดสอบ Delete
	err = repo.Delete(created.ID, testUserID)
	if err != nil {
		t.Fatalf("unexpected error deleting task: %v", err)
	}

	// 5. ตรวจสอบว่า GetByID หลังจากลบแล้วต้องได้ ErrTaskNotFound
	_, err = repo.GetByID(created.ID, testUserID)
	if err != ErrTaskNotFound {
		t.Errorf("expected ErrTaskNotFound, got %v", err)
	}
}
