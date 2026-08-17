package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oopbest/task-app/models"
	"github.com/oopbest/task-app/repository"
	"github.com/oopbest/task-app/utils"
)

type AuthHandler struct {
	userRepo repository.UserRepository
}

func NewAuthHandler(userRepo repository.UserRepository) *AuthHandler {
	return &AuthHandler{userRepo: userRepo}
}

// Register godoc
// @Summary สมัครสมาชิกใหม่ (Register)
// @Description สร้างบัญชีผู้ใช้ใหม่ด้วย Email และ Password
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body models.RegisterInput true "ข้อมูลการสมัครสมาชิก"
// @Success 201 {object} models.AuthResponse
// @Failure 400 {object} map[string]string "ข้อมูลไม่ถูกต้อง"
// @Failure 409 {object} map[string]string "Email นี้ถูกใช้งานแล้ว"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var input models.RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON request body"})
		return
	}

	if err := input.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user, err := h.userRepo.Create(input.Email, hashedPassword)
	if err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, models.AuthResponse{
		Token: token,
		User:  user,
	})
}

// Login godoc
// @Summary เข้าสู่ระบบ (Login)
// @Description เข้าสู่ระบบด้วย Email และ Password เพื่อรับ JWT Token
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body models.LoginInput true "ข้อมูลเข้าสู่ระบบ"
// @Success 200 {object} models.AuthResponse
// @Failure 400 {object} map[string]string "ข้อมูลไม่ถูกต้อง"
// @Failure 401 {object} map[string]string "Email หรือ Password ไม่ถูกต้อง"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var input models.LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON request body"})
		return
	}

	if err := input.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userRepo.GetByEmail(input.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if !utils.CheckPasswordHash(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, models.AuthResponse{
		Token: token,
		User:  user,
	})
}
