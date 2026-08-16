package models

import (
	"errors"
	"strings"
	"time"
)

// User โครงสร้างข้อมูลของผู้ใช้ในระบบ
type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // เครื่องหมาย "-" เพื่อไม่ให้ส่ง Password ออกไปใน JSON response
	CreatedAt time.Time `json:"created_at"`
}

// RegisterInput ข้อมูลที่รับเข้ามาตอนสมัครสมาชิก
type RegisterInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate ตรวจสอบความถูกต้องของการสมัคร
func (input *RegisterInput) Validate() error {
	if strings.TrimSpace(input.Email) == "" || !strings.Contains(input.Email, "@") {
		return errors.New("a valid email is required")
	}
	if len(input.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}
	return nil
}

// LoginInput ข้อมูลที่รับเข้ามาตอนเข้าสู่ระบบ
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate ตรวจสอบความถูกต้องของการ Login
func (input *LoginInput) Validate() error {
	if strings.TrimSpace(input.Email) == "" {
		return errors.New("email is required")
	}
	if strings.TrimSpace(input.Password) == "" {
		return errors.New("password is required")
	}
	return nil
}

// AuthResponse ผลลัพธ์ที่ส่งกลับเมื่อ Login/Register สำเร็จ
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
