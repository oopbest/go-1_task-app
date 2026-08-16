package models

import "strings"

// TaskFilter พารามิเตอร์สำหรับการค้นหา กรอง และแบ่งหน้า
type TaskFilter struct {
	Page      int    // หน้าปัจจุบัน (Default: 1)
	Limit     int    // จำนวนรายการต่อหน้า (Default: 10, Max: 100)
	Search    string // คำค้นหาใน title หรือ description
	Completed *bool  // กรองตามสถานะ true / false (ถ้าเป็น nil คือดึงทั้งหมด)
	SortBy    string // เรียงตามฟิลด์ id, title, created_at (Default: created_at)
	Order     string // asc หรือ desc (Default: desc)
}

// Sanitize ตรวจสอบและตั้งค่า Default ให้กับ Filter เพื่อความปลอดภัย
func (f *TaskFilter) Sanitize() {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Limit <= 0 {
		f.Limit = 10
	}
	if f.Limit > 100 {
		f.Limit = 100 // จำกัดไม่ให้ดึงเกิน 100 รายการต่อครั้งเพื่อป้องกัน memory เต็ม
	}

	f.Search = strings.TrimSpace(f.Search)

	// อนุญาตเฉพาะฟิลด์ที่ปลอดภัย ป้องกัน SQL Injection ใน ORDER BY
	validSorts := map[string]bool{"id": true, "title": true, "created_at": true}
	if !validSorts[strings.ToLower(f.SortBy)] {
		f.SortBy = "created_at"
	}

	if strings.ToLower(f.Order) != "asc" {
		f.Order = "desc"
	} else {
		f.Order = "asc"
	}
}

// PaginatedTasks โครงสร้าง JSON Response สำหรับส่งรายการข้อมูลพร้อม Pagination Metadata
type PaginatedTasks struct {
	Data       []Task `json:"data"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	TotalItems int    `json:"total_items"`
	TotalPages int    `json:"total_pages"`
}
