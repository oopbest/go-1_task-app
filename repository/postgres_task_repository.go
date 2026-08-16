package repository

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq" // ลงทะเบียน Postgres driver
	"github.com/oopbest/task-app/models"
)

// PostgresTaskRepository จัดการข้อมูล Task บนฐานข้อมูล PostgreSQL จริง
type PostgresTaskRepository struct {
	db *sql.DB
}

// NewPostgresTaskRepository Constructor สำหรับสร้าง Repository และสร้าง Table อัตโนมัติถ้ายังไม่มี
func NewPostgresTaskRepository(db *sql.DB) (*PostgresTaskRepository, error) {
	repo := &PostgresTaskRepository{db: db}

	// สร้าง Table tasks อัตโนมัติ (Auto Migration แบบง่าย)
	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		completed BOOLEAN NOT NULL DEFAULT FALSE,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	ALTER TABLE tasks ADD COLUMN IF NOT EXISTS user_id INT REFERENCES users(id) ON DELETE CASCADE;
	`

	if _, err := db.Exec(query); err != nil {
		return nil, fmt.Errorf("failed to create tasks table: %w", err)
	}

	return repo, nil
}

// GetAll ดึงรายการ Task ทั้งหมดจาก PostgreSQL
func (r *PostgresTaskRepository) GetAll(userID int) []models.Task {
	query := `SELECT id, title, description, completed, COALESCE(user_id, 0), created_at FROM tasks WHERE user_id = $1 ORDER BY id ASC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return []models.Task{}
	}
	defer rows.Close()

	tasks := make([]models.Task, 0)
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Completed, &task.UserID, &task.CreatedAt); err != nil {
			continue
		}
		tasks = append(tasks, task)
	}

	// เช็ค error ที่อาจเกิดขึ้นระหว่างการวนลูปอ่าน rows (Best Practice)
	if err := rows.Err(); err != nil {
		return []models.Task{}
	}

	return tasks
}

// GetByID ดึงข้อมูล Task ตาม ID
func (r *PostgresTaskRepository) GetByID(id int, userID int) (models.Task, error) {
	query := `SELECT id, title, description, completed, COALESCE(user_id, 0), created_at FROM tasks WHERE id = $1 AND user_id = $2`

	var task models.Task
	err := r.db.QueryRow(query, id, userID).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.UserID,
		&task.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Task{}, ErrTaskNotFound
		}
		return models.Task{}, err
	}

	return task, nil
}

// Create สร้าง Task ใหม่และรับ ID + CreatedAt กลับมาจาก DB
func (r *PostgresTaskRepository) Create(input models.CreateTaskInput, userID int) models.Task {
	query := `
	INSERT INTO tasks (title, description, user_id)
	VALUES ($1, $2, $3)
	RETURNING id, title, description, completed, user_id, created_at`

	var task models.Task
	_ = r.db.QueryRow(query, input.Title, input.Description, userID).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.UserID,
		&task.CreatedAt,
	)

	return task
}

// Update แก้ไข Task ตาม ID
func (r *PostgresTaskRepository) Update(id int, input models.UpdateTaskInput, userID int) (models.Task, error) {
	// 1. ดึงข้อมูลเดิมออกมาก่อน
	existing, err := r.GetByID(id, userID)
	if err != nil {
		return models.Task{}, err
	}

	// 2. อัปเดตเฉพาะฟิลด์ที่ส่งค่ามา
	if input.Title != nil {
		existing.Title = *input.Title
	}
	if input.Description != nil {
		existing.Description = *input.Description
	}
	if input.Completed != nil {
		existing.Completed = *input.Completed
	}

	// 3. บันทึกลง PostgreSQL
	query := `
	UPDATE tasks
	SET title = $1, description = $2, completed = $3
	WHERE id = $4 AND user_id = $5
	RETURNING id, title, description, completed, user_id, created_at`

	var updated models.Task
	err = r.db.QueryRow(query, existing.Title, existing.Description, existing.Completed, id, userID).Scan(
		&updated.ID,
		&updated.Title,
		&updated.Description,
		&updated.Completed,
		&updated.UserID,
		&updated.CreatedAt,
	)

	if err != nil {
		return models.Task{}, err
	}

	return updated, nil
}

// Delete ลบ Task ตาม ID
func (r *PostgresTaskRepository) Delete(id int, userID int) error {
	query := `DELETE FROM tasks WHERE id = $1 AND user_id = $2`

	result, err := r.db.Exec(query, id, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return ErrTaskNotFound
	}

	return nil
}
