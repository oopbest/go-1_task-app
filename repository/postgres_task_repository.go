package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strings"

	_ "github.com/lib/pq"
	"github.com/oopbest/task-app/models"
)

type PostgresTaskRepository struct {
	db *sql.DB
}

func NewPostgresTaskRepository(db *sql.DB) *PostgresTaskRepository {
	return &PostgresTaskRepository{db: db}
}

// GetAll ดึงรายการ Task ตามเงื่อนไข Search, Filter, Sort และ Pagination
func (r *PostgresTaskRepository) GetAll(userID int, filter models.TaskFilter) (models.PaginatedTasks, error) {
	filter.Sanitize()

	// 1. สร้างเงื่อนไข WHERE แบบ Dynamic
	conditions := []string{"user_id = $1"}
	args := []any{userID}
	argIdx := 2

	// ค้นหาข้อความใน title หรือ description (ILIKE = Case-insensitive search ใน Postgres)
	if filter.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(title ILIKE $%d OR description ILIKE $%d)", argIdx, argIdx))
		args = append(args, "%"+filter.Search+"%")
		argIdx++
	}

	// กรองตามสถานะ completed (ถ้ามีการส่งค่ามา)
	if filter.Completed != nil {
		conditions = append(conditions, fmt.Sprintf("completed = $%d", argIdx))
		args = append(args, *filter.Completed)
		argIdx++
	}

	whereClause := strings.Join(conditions, " AND ")

	// 2. นับจำนวนรายการทั้งหมด (Total Items)
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM tasks WHERE %s", whereClause)
	var totalItems int
	if err := r.db.QueryRow(countQuery, args...).Scan(&totalItems); err != nil {
		return models.PaginatedTasks{}, err
	}

	// 3. ดึงรายการข้อมูลตาม Pagination (LIMIT & OFFSET)
	offset := (filter.Page - 1) * filter.Limit
	dataQuery := fmt.Sprintf(`
		SELECT id, title, description, completed, user_id, created_at
		FROM tasks
		WHERE %s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d`,
		whereClause, filter.SortBy, filter.Order, argIdx, argIdx+1,
	)

	dataArgs := append(args, filter.Limit, offset)

	rows, err := r.db.Query(dataQuery, dataArgs...)
	if err != nil {
		return models.PaginatedTasks{}, err
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

	if err := rows.Err(); err != nil {
		return models.PaginatedTasks{}, err
	}

	// 4. คำนวณจำนวนหน้าทั้งหมด (Total Pages)
	totalPages := 0
	if totalItems > 0 {
		totalPages = int(math.Ceil(float64(totalItems) / float64(filter.Limit)))
	}

	return models.PaginatedTasks{
		Data:       tasks,
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}, nil
}

func (r *PostgresTaskRepository) GetByID(id int, userID int) (models.Task, error) {
	query := `SELECT id, title, description, completed, user_id, created_at FROM tasks WHERE id = $1 AND user_id = $2`

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

func (r *PostgresTaskRepository) Update(id int, input models.UpdateTaskInput, userID int) (models.Task, error) {
	existing, err := r.GetByID(id, userID)
	if err != nil {
		return models.Task{}, err
	}

	if input.Title != nil {
		existing.Title = *input.Title
	}
	if input.Description != nil {
		existing.Description = *input.Description
	}
	if input.Completed != nil {
		existing.Completed = *input.Completed
	}

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
