package models

import (
	"errors"
	"strings"
	"time"
)

// Struct: คือประเภทข้อมูลแบบโครงสร้าง (คล้าย Class ในภาษาอื่น แต่ Go ไม่มี Class / Inheritance)
type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	UserID      int       `json:"user_id"` // เพิ่มฟิลด์นี้
	CreatedAt   time.Time `json:"created_at"`
}

// CreateTaskInput คือข้อมูลที่ Client ต้องส่งมาตอนสร้าง Task (POST /tasks
type CreateTaskInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (input *CreateTaskInput) Validate() error {
	if strings.TrimSpace(input.Title) == "" {
		return errors.New("title is required")
	}
	return nil
}

// UpdateTaskInput คือข้อมูลที่ Client ต้องส่งมาตอนอัปเดต Task (PUT /tasks/:id)
// Pointer (*string, *bool): ใช้ใน DTO Update เพื่อแยกแยะว่าผู้ใช้ "ส่งค่าว่างมา" หรือ "ไม่ได้ส่งฟิลด์นี้มาเลย (nil)"
type UpdateTaskInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Completed   *bool   `json:"completed"`
}

// Validate ตรวจสอบว่ามีการส่งฟิลด์ใดฟิลด์หนึ่งมาอัปเดตหรือไม่
func (input *UpdateTaskInput) Validate() error {
	if input.Title != nil && strings.TrimSpace(*input.Title) == "" {
		return errors.New("title is required")
	}
	if input.Description != nil && strings.TrimSpace(*input.Description) == "" {
		return errors.New("description is required")
	}
	return nil
}
