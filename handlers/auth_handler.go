package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/oopbest/task-app/models"
	"github.com/oopbest/task-app/repository"
	"github.com/oopbest/task-app/utils"
)

// AuthHandler จัดการ HTTP Request ด้าน Authentication
type AuthHandler struct {
	userRepo repository.UserRepository
}

// NewAuthHandler Constructor สำหรับสร้าง AuthHandler
func NewAuthHandler(userRepo repository.UserRepository) *AuthHandler {
	return &AuthHandler{userRepo: userRepo}
}

// Register จัดการ POST /auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input models.RegisterInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid JSON request body")
		return
	}

	if err := input.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// 1. เข้ารหัสผ่านด้วย bcrypt
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// 2. บันทึกลงฐานข้อมูล
	user, err := h.userRepo.Create(input.Email, hashedPassword)
	if err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			respondError(w, http.StatusConflict, "User with this email already exists")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// 3. ออก JWT Token
	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	respondJSON(w, http.StatusCreated, models.AuthResponse{
		Token: token,
		User:  user,
	})
}

// Login จัดการ POST /auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input models.LoginInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid JSON request body")
		return
	}

	if err := input.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// 1. ค้นหา User จาก Email
	user, err := h.userRepo.GetByEmail(input.Email)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// 2. ตรวจสอบรหัสผ่าน
	if !utils.CheckPasswordHash(input.Password, user.Password) {
		respondError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// 3. ออก JWT Token
	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	respondJSON(w, http.StatusOK, models.AuthResponse{
		Token: token,
		User:  user,
	})
}
