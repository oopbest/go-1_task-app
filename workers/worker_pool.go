package workers

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/oopbest/task-app/metrics"
)

// Job โครงสร้างข้อมูลงานเบื้องหลังที่ต้องการให้ Worker นำไปทำ
type Job struct {
	Type    string // เช่น "TASK_CREATED", "TASK_UPDATED", "TASK_DELETED"
	UserID  int
	TaskID  int
	Title   string
	Payload any
}

// WorkerPool จัดการคิวงานและกลุ่ม Goroutines ที่ประมวลผลงานเบื้องหลัง
type WorkerPool struct {
	numWorkers int
	jobQueue   chan Job
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewWorkerPool Constructor สำหรับสร้าง WorkerPool
func NewWorkerPool(numWorkers int, queueSize int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		numWorkers: numWorkers,
		jobQueue:   make(chan Job, queueSize), // Buffered Channel
		ctx:        ctx,
		cancel:     cancel,
	}
}

// Start เริ่มต้นการทำงานของ Worker Goroutines ตามจำนวนที่กำหนด
func (wp *WorkerPool) Start() {
	for i := 1; i <= wp.numWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
	log.Printf("⚙️ Worker Pool started with %d background workers\n", wp.numWorkers)
}

// worker การทำงานของ Goroutine แต่ละตัวที่คอยดึงงานจาก Channel ไปทำ
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	for {
		select {
		case <-wp.ctx.Done():
			// ถ้าได้รับคำสั่ง Cancel ให้เคลียร์งานที่ค้างในคิวให้หมดก่อนปิด
			for job := range wp.jobQueue {
				wp.processJob(id, job)
			}
			log.Printf("🛑 Worker %d stopped gracefully\n", id)
			return

		case job, ok := <-wp.jobQueue:
			if !ok {
				// Channel ถูกปิดแล้ว
				log.Printf("🛑 Worker %d stopped (channel closed)\n", id)
				return
			}
			wp.processJob(id, job)
		}
	}
}

// processJob ฟังก์ชันจำลองการประมวลผลงานหนักเบื้องหลัง (เช่น ส่ง Email / ยิง Webhook)
func (wp *WorkerPool) processJob(workerID int, job Job) {
	log.Printf("👷 [Worker %d] START processing job: [%s] for Task #%d (User #%d: '%s')\n",
		workerID, job.Type, job.TaskID, job.UserID, job.Title)

	// จำลองเวลาในการทำงานเบื้องหลัง 1.5 วินาที
	time.Sleep(1500 * time.Millisecond)

	// 📊 บันทึกสถิติงานเข้า Prometheus
	metrics.WorkerJobsTotal.WithLabelValues(job.Type, "success").Inc()

	log.Printf("✅ [Worker %d] COMPLETED job: [%s] for Task #%d (Notification sent successfully!)\n",
		workerID, job.Type, job.TaskID)
}

// Enqueue ส่งงานเข้าคิว (Non-blocking)
func (wp *WorkerPool) Enqueue(job Job) {
	select {
	case wp.jobQueue <- job:
		// ส่งเข้า Channel สำเร็จ
	default:
		// กรณีที่คิวเต็ม 100 งาน จะแจ้งเตือนเพื่อไม่ให้บล็อก HTTP Request
		log.Printf("⚠️ Worker Pool queue is FULL! Dropping job: [%s] for Task #%d\n", job.Type, job.TaskID)
	}
}

// Stop ปิด Worker Pool อย่างปลอดภัย (Graceful Shutdown)
func (wp *WorkerPool) Stop() {
	log.Println("🛑 Stopping Worker Pool, waiting for pending jobs to finish...")
	close(wp.jobQueue) // ปิด Channel ไม่รับงานใหม่
	wp.cancel()        // แจ้งเตือน Context
	wp.wg.Wait()       // รอให้ Worker ทุกตัวทำงานที่ค้างอยู่จนเสร็จ 100%
	log.Println("✅ All background workers stopped cleanly")
}
