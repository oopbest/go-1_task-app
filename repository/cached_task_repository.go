package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/oopbest/task-app/models"
	"github.com/redis/go-redis/v9"
)

// CachedTaskRepository Decorator สำหรับทำ Caching ครอบ TaskRepository
type CachedTaskRepository struct {
	next TaskRepository
	rdb  *redis.Client
	ttl  time.Duration
}

// NewCachedTaskRepository Constructor สำหรับสร้าง Cached Repository
func NewCachedTaskRepository(next TaskRepository, rdb *redis.Client, ttl time.Duration) *CachedTaskRepository {
	return &CachedTaskRepository{
		next: next,
		rdb:  rdb,
		ttl:  ttl,
	}
}

// สร้าง Cache Key สำหรับรายการ Tasks (เช่น tasks:u:1:p:1:l:10:s:test:c:all:sb:created_at:o:desc)
func (r *CachedTaskRepository) listCacheKey(userID int, filter models.TaskFilter) string {
	completedStr := "all"
	if filter.Completed != nil {
		completedStr = fmt.Sprintf("%v", *filter.Completed)
	}
	return fmt.Sprintf("tasks:u:%d:p:%d:l:%d:s:%s:c:%s:sb:%s:o:%s",
		userID, filter.Page, filter.Limit, filter.Search, completedStr, filter.SortBy, filter.Order)
}

// สร้าง Cache Key สำหรับ Task เดี่ยว
func (r *CachedTaskRepository) singleCacheKey(id, userID int) string {
	return fmt.Sprintf("task:u:%d:id:%d", userID, id)
}

// invalidateUserCache ลบ Cache ทั้งหมดของ User คนนั้นเมื่อมีการเปลี่ยนแปลงข้อมูล
func (r *CachedTaskRepository) invalidateUserCache(ctx context.Context, userID int) {
	pattern := fmt.Sprintf("tasks:u:%d:*", userID)
	iter := r.rdb.Scan(ctx, 0, pattern, 0).Iterator()

	keys := make([]string, 0)
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		log.Printf("⚠️ Redis Scan error: %v", err)
	}

	if len(keys) > 0 {
		_ = r.rdb.Del(ctx, keys...).Err()
	}
}

// GetAll (Cache-Aside Pattern)
func (r *CachedTaskRepository) GetAll(userID int, filter models.TaskFilter) (models.PaginatedTasks, error) {
	ctx := context.Background()
	key := r.listCacheKey(userID, filter)

	// 1. ลองอ่านจาก Redis Cache ก่อน
	cachedData, err := r.rdb.Get(ctx, key).Result()
	if err == nil {
		var result models.PaginatedTasks
		if err := json.Unmarshal([]byte(cachedData), &result); err == nil {
			return result, nil // ⚡ Cache Hit!
		}
	}

	// 2. ถ้าไม่มีใน Cache (Cache Miss) ให้ไปดึงจาก PostgreSQL
	result, err := r.next.GetAll(userID, filter)
	if err != nil {
		return models.PaginatedTasks{}, err
	}

	// 3. นำผลลัพธ์ไปบันทึกลง Redis พร้อมตั้ง TTL
	if dataBytes, err := json.Marshal(result); err == nil {
		_ = r.rdb.Set(ctx, key, dataBytes, r.ttl).Err()
	}

	return result, nil
}

// GetByID (Cache-Aside Pattern)
func (r *CachedTaskRepository) GetByID(id int, userID int) (models.Task, error) {
	ctx := context.Background()
	key := r.singleCacheKey(id, userID)

	cachedData, err := r.rdb.Get(ctx, key).Result()
	if err == nil {
		var task models.Task
		if err := json.Unmarshal([]byte(cachedData), &task); err == nil {
			return task, nil // ⚡ Cache Hit!
		}
	}

	task, err := r.next.GetByID(id, userID)
	if err != nil {
		return models.Task{}, err
	}

	if dataBytes, err := json.Marshal(task); err == nil {
		_ = r.rdb.Set(ctx, key, dataBytes, r.ttl).Err()
	}

	return task, nil
}

// Create บันทึกลง DB แล้วลบ Cache ของ User ทิ้งทันที
func (r *CachedTaskRepository) Create(input models.CreateTaskInput, userID int) models.Task {
	task := r.next.Create(input, userID)
	r.invalidateUserCache(context.Background(), userID)
	return task
}

// Update บันทึกลง DB แล้วลบ Cache ของ User ทิ้งทันที
func (r *CachedTaskRepository) Update(id int, input models.UpdateTaskInput, userID int) (models.Task, error) {
	task, err := r.next.Update(id, input, userID)
	if err != nil {
		return models.Task{}, err
	}

	ctx := context.Background()
	_ = r.rdb.Del(ctx, r.singleCacheKey(id, userID)).Err()
	r.invalidateUserCache(ctx, userID)

	return task, nil
}

// Delete ลบจาก DB แล้วลบ Cache ของ User ทิ้งทันที
func (r *CachedTaskRepository) Delete(id int, userID int) error {
	if err := r.next.Delete(id, userID); err != nil {
		return err
	}

	ctx := context.Background()
	_ = r.rdb.Del(ctx, r.singleCacheKey(id, userID)).Err()
	r.invalidateUserCache(ctx, userID)

	return nil
}
