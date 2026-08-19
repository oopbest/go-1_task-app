package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	pb "github.com/oopbest/task-app/proto/notification"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// notificationServer โครงสร้าง Server ที่ implement NotificationServiceServer
type notificationServer struct {
	pb.UnimplementedNotificationServiceServer
	totalSent int64 // Atomic Counter เก็บจำนวนการแจ้งเตือนทั้งหมด
}

// SendNotification รับคำสั่งส่งการแจ้งเตือนจาก Service อื่นผ่าน gRPC
func (s *notificationServer) SendNotification(ctx context.Context, req *pb.NotificationRequest) (*pb.NotificationResponse, error) {
	atomic.AddInt64(&s.totalSent, 1)

	log.Printf("📬 [gRPC Server] RECEIVED: Event='%s' | TaskID=%d | UserID=%d | Title='%s'\n",
		req.GetEventType(), req.GetTaskId(), req.GetUserId(), req.GetTitle())

	// จำลองเวลาส่ง Notification (เช่น ยิง SMTP Email หรือ Firebase Notification)
	time.Sleep(500 * time.Millisecond)

	return &pb.NotificationResponse{
		Success: true,
		Message: fmt.Sprintf("Notification for Task #%d successfully sent to User #%d", req.GetTaskId(), req.GetUserId()),
		SentAt:  time.Now().Format(time.RFC3339),
	}, nil
}

// GetNotificationStats ดูสถิติการแจ้งเตือนทั้งหมด
func (s *notificationServer) GetNotificationStats(ctx context.Context, req *pb.StatsRequest) (*pb.StatsResponse, error) {
	count := atomic.LoadInt64(&s.totalSent)
	return &pb.StatsResponse{
		TotalSent: count,
	}, nil
}

func main() {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50051"
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	// สร้าง gRPC Server
	grpcServer := grpc.NewServer()
	serverInstance := &notificationServer{}
	pb.RegisterNotificationServiceServer(grpcServer, serverInstance)

	// เปิด Reflection สำหรับการ Debug ด้วยเครื่องมือ gRPC GUI (เช่น Postman / Evans)
	reflection.Register(grpcServer)

	go func() {
		fmt.Println("==================================================")
		fmt.Printf("📬 Notification Microservice (gRPC) listening on port :%s\n", port)
		fmt.Println("==================================================")
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("gRPC Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down Notification Microservice gracefully...")
	grpcServer.GracefulStop()
	log.Println("✅ Notification Microservice stopped cleanly")
}
