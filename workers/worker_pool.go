package workers

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/oopbest/task-app/metrics"
	pb "github.com/oopbest/task-app/proto/notification"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	grpcClient pb.NotificationServiceClient
	grpcConn   *grpc.ClientConn
}

// NewWorkerPool Constructor สำหรับสร้าง WorkerPool พร้อมเชื่อมต่อ gRPC Microservice
func NewWorkerPool(numWorkers int, queueSize int, grpcServerAddr string) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	// เชื่อมต่อ gRPC Notification Microservice
	conn, err := grpc.NewClient(grpcServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	var client pb.NotificationServiceClient
	if err != nil {
		log.Printf("⚠️ Warning: Could not connect to gRPC Notification Service at %s: %v\n", grpcServerAddr, err)
	} else {
		client = pb.NewNotificationServiceClient(conn)
		log.Printf("🌐 Connected to gRPC Notification Service at %s\n", grpcServerAddr)
	}

	return &WorkerPool{
		numWorkers: numWorkers,
		jobQueue:   make(chan Job, queueSize),
		ctx:        ctx,
		cancel:     cancel,
		grpcClient: client,
		grpcConn:   conn,
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
			for job := range wp.jobQueue {
				wp.processJob(id, job)
			}
			log.Printf("🛑 Worker %d stopped gracefully\n", id)
			return

		case job, ok := <-wp.jobQueue:
			if !ok {
				log.Printf("🛑 Worker %d stopped (channel closed)\n", id)
				return
			}
			wp.processJob(id, job)
		}
	}
}

// processJob ฟังก์ชันยิง gRPC Call ข้ามไปยัง Notification Microservice
func (wp *WorkerPool) processJob(workerID int, job Job) {
	log.Printf("👷 [Worker %d] START processing job: [%s] for Task #%d (User #%d: '%s')\n",
		workerID, job.Type, job.TaskID, job.UserID, job.Title)

	// ⚡ ยิง gRPC Call ข้าม Service ไปยัง Notification Microservice
	if wp.grpcClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		resp, err := wp.grpcClient.SendNotification(ctx, &pb.NotificationRequest{
			EventType: job.Type,
			TaskId:    int64(job.TaskID),
			UserId:    int64(job.UserID),
			Title:     job.Title,
		})

		if err != nil {
			log.Printf("❌ [Worker %d] gRPC Error calling Notification Service: %v\n", workerID, err)
			metrics.WorkerJobsTotal.WithLabelValues(job.Type, "error").Inc()
			return
		}

		log.Printf("✅ [Worker %d] gRPC SUCCESS: %s (at %s)\n", workerID, resp.GetMessage(), resp.GetSentAt())
	} else {
		// Fallback จำลองถ้าไม่มี gRPC
		time.Sleep(1000 * time.Millisecond)
		log.Printf("✅ [Worker %d] COMPLETED job locally: [%s] for Task #%d\n", workerID, job.Type, job.TaskID)
	}

	metrics.WorkerJobsTotal.WithLabelValues(job.Type, "success").Inc()
}

// Enqueue ส่งงานเข้าคิว (Non-blocking)
func (wp *WorkerPool) Enqueue(job Job) {
	select {
	case wp.jobQueue <- job:
	default:
		log.Printf("⚠️ Worker Pool queue is FULL! Dropping job: [%s] for Task #%d\n", job.Type, job.TaskID)
	}
}

// Stop ปิด Worker Pool และปิด gRPC Connection
func (wp *WorkerPool) Stop() {
	log.Println("🛑 Stopping Worker Pool, waiting for pending jobs to finish...")
	close(wp.jobQueue)
	wp.cancel()
	wp.wg.Wait()

	if wp.grpcConn != nil {
		_ = wp.grpcConn.Close()
	}
	log.Println("✅ All background workers and gRPC connections stopped cleanly")
}
