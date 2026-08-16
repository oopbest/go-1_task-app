package repository

import (
	"database/sql"
	"errors"

	"github.com/oopbest/task-app/models"
)

var ErrUserAlreadyExists = errors.New("user with this email already exists")

// UserRepository Interface สำหรับจัดการผู้ใช้
type UserRepository interface {
	Create(email, passwordHash string) (models.User, error)
	GetByEmail(email string) (models.User, error)
	GetByID(id int) (models.User, error)
}

// PostgresUserRepository จัดการตาราง users บน PostgreSQL
type PostgresUserRepository struct {
	db *sql.DB
}

// NewPostgresUserRepository Constructor
func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

// Create บันทึก User ใหม่ลงฐานข้อมูล
func (r *PostgresUserRepository) Create(email, passwordHash string) (models.User, error) {
	query := `
	INSERT INTO users (email, password)
	VALUES ($1, $2)
	RETURNING id, email, created_at`

	var user models.User
	err := r.db.QueryRow(query, email, passwordHash).Scan(
		&user.ID,
		&user.Email,
		&user.CreatedAt,
	)

	if err != nil {
		return models.User{}, ErrUserAlreadyExists
	}

	return user, nil
}

// GetByEmail ค้นหา User ด้วย Email สำหรับขั้นตอน Login
func (r *PostgresUserRepository) GetByEmail(email string) (models.User, error) {
	query := `SELECT id, email, password, created_at FROM users WHERE email = $1`

	var user models.User
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, ErrTaskNotFound
		}
		return models.User{}, err
	}

	return user, nil
}

// GetByID ดึงข้อมูล User ตาม ID
func (r *PostgresUserRepository) GetByID(id int) (models.User, error) {
	query := `SELECT id, email, created_at FROM users WHERE id = $1`

	var user models.User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, ErrTaskNotFound
		}
		return models.User{}, err
	}

	return user, nil
}
