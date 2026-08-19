# 📝 Task Management REST API (Go Zero to Hero)

RESTful API สำหรับจัดการงานแบบ **Multi-User** ระดับ **Enterprise-Grade** พัฒนาด้วย **Gin Web Framework** เชื่อมต่อฐานข้อมูล **PostgreSQL** เสริมความเร็วระดับ **Microseconds ด้วย Redis Cache**, ระบบประมวลผลเบื้องหลัง **Background Worker Pool (Goroutines & Channels)**, สถาปัตยกรรม **Distributed Microservices ด้วย gRPC & Protocol Buffers**, ความปลอดภัย **JWT Authentication**, การจัดการ Schema ด้วย **Database Migrations**, เอกสาร **Swagger (OpenAPI) Interactive UI**, และระบบเฝ้าระวังครบวงจร **Observability (Structured Logging `slog` + Prometheus Metrics + Grafana Dashboards)** 🚀⚡🌐📊📈

---

## ✨ Features

- 🌐 **Microservices Architecture with gRPC & Protocol Buffers**:
  - แยก Service ออกเป็น **Task API Gateway (:8080)** และ **Notification Microservice (:50051)**
  - สื่อสารระหว่าง Services ข้ามเครือข่ายด้วย **gRPC บน HTTP/2 Binary Protocol** ที่เร็วกว่า JSON 5–10 เท่า
  - นิยาม Single Source of Truth Data Contract ด้วย **Protocol Buffers (`.proto`)**
- 📊 **Full-Stack Observability & Monitoring**:
  - **Structured JSON Logging (`log/slog`)**: บันทึก Log ทุก Request ในรูปแบบ JSON ตามมาตรฐาน Go 1.21+
  - **Prometheus Metrics Exporter**: เก็บสถิติ Request Count, Latency Histogram, และ Worker Jobs ผ่าน Endpoint `/metrics`
  - **Grafana Live Dashboards**: แสดงผลกราฟสถิติ Real-Time บนเว็บ `http://localhost:3000` (admin/admin)
- 🏎️ **Gin Web Framework & High Performance Routing**:
  - Radix Tree Router ความเร็วสูงและใช้ Memory น้อยที่สุด
  - Route Grouping (`/auth`, `/tasks`) พร้อม Gin Middleware Pipeline
  - Auto JSON Binding & Context Helpers
- 📖 **Swagger / OpenAPI Interactive UI**:
  - สร้างเอกสาร API Documentation อัตโนมัติด้วย `swaggo/swag`
  - สามารถเปิดทดสอบ API ผ่านหน้าเว็บเบราว์เซอร์ได้ที่: `http://localhost:8080/swagger/index.html`
  - รองรับการกรอก JWT Token ผ่านปุ่ม **Authorize 🔒**
- ⚡ **Background Worker Pool (Go Concurrency & Channels)**:
  - ประมวลผลงานเบื้องหลังแบบ Asynchronous (เช่น จำลองการส่ง Email / Webhook แจ้งเตือน)
  - ควบคุมจำนวน Worker คงที่ (3 Goroutines) และขนาดคิวงานใน RAM (`chan Job` ขนาด 100) ป้องกัน Server Overload
  - **Graceful Worker Shutdown**: ใช้ `sync.WaitGroup` รอให้งานที่ค้างในคิวทำให้เสร็จสมบูรณ์ก่อนปิดเซิร์ฟเวอร์
- 🏎️ **Redis In-Memory Caching (Sub-Millisecond Latency)**:
  - **Cache-Aside Pattern**: ตรวจสอบ Redis Cache ก่อนดึง PostgreSQL ลด Latency เหลือเพียง **~500 µs (ไมโครวินาที)**
  - **Cache Invalidation**: เคลียร์ Cache อัตโนมัติทันทีที่มีการ Create, Update, Delete ป้องกันข้อมูลเก่าค้าง
  - **Decorator Architecture**: แยก Caching Layer ออกจาก Business Logic ด้วย Decorator Pattern
- 🔐 **User Authentication (JWT & Bcrypt)**:
  - สมัครสมาชิก (`/auth/register`) และเข้าสู่ระบบ (`/auth/login`)
  - เข้ารหัสผ่านอย่างปลอดภัยด้วย `bcrypt`
  - ยืนยันตัวตนด้วย `JWT (JSON Web Token)` มีอายุ 24 ชั่วโมง
- 👥 **Multi-User Task Isolation**: แยกข้อมูลงานของผู้ใช้แต่ละคนอย่างเด็ดขาด (Row-level Security)
- 🗄️ **Database Migrations & Indexing (`golang-migrate`)**:
  - จัดการ Schema ด้วยไฟล์ `.up.sql` และ `.down.sql`
  - B-Tree Indexes บนคอลัมน์ `user_id` และ `completed` เพิ่มความเร็วในการค้นหา
- 🔍 **Dynamic Querying, Pagination, Search & Sorting**:
  - รองรับ `?page=1&limit=10&search=golang&completed=false&sort=created_at&order=desc`
  - ป้องกัน SQL Injection ใน `ORDER BY` ด้วย Go Map Whitelist
- 🧪 **Automated Unit Testing**: เขียน Test ด้วย `testing` และ `net/http/httptest`
- 🛑 **Graceful Shutdown**: ดักจับ OS Signals ปิดเซิร์ฟเวอร์และเคลียร์ Connections อย่างปลอดภัย
- 🐳 **Docker Compose & Multi-Stage Build**: รัน PostgreSQL + Redis + Prometheus + Grafana ครบวงจร

---

## 📁 Project Structure

```
1_task-app/
├── cmd/
│   └── notification-service/ # Notification Microservice (gRPC Server)
├── proto/
│   ├── notification.proto    # Protobuf Contract Definition
│   └── notification/         # Generated Go Protobuf & gRPC Stubs
├── docs/                     # Swagger / OpenAPI Generated Files
├── metrics/                  # Prometheus Metrics Collector & Middleware
├── migrations/             # Database Schema Migrations (.up.sql / .down.sql)
├── models/
│   ├── user.go               # Data Models สำหรับ User & Auth DTOs
│   ├── task.go               # Data Models สำหรับ Task & DTOs
│   └── pagination.go         # DTOs สำหรับ Pagination, Search & Filter
├── repository/
│   ├── user_repository.go    # PostgreSQL User Repository
│   ├── task_repository.go    # TaskRepository Interface & In-Memory Storage
│   ├── postgres_task_repository.go # PostgreSQL Dynamic Query Repository
│   ├── cached_task_repository.go   # Redis Caching Decorator Repository
│   └── migrations.go       # Migration Runner
├── workers/
│   └── worker_pool.go        # Background Worker Pool & gRPC Client
├── handlers/
│   ├── auth_handler.go       # Gin HTTP Handlers สำหรับ Register & Login
│   └── task_handler.go       # Gin HTTP Handlers สำหรับ Tasks CRUD & Worker Enqueue
├── middleware/
│   ├── gin_auth.go           # Gin JWT Authentication Middleware
│   ├── slog_logger.go        # Structured JSON Logger Middleware (log/slog)
│   ├── auth.go             # Context Auth Helper
│   └── middleware.go       # Legacy Middleware
├── utils/
│   └── auth.go             # Bcrypt Password Hashing & JWT Helpers
├── prometheus.yml            # Prometheus Scrape Configuration
├── docker-compose.yml        # Postgres + Redis + Prometheus + Grafana
├── Dockerfile                # Multi-Stage Dockerfile
├── go.mod                  # Go Module Definition
└── main.go                 # Application Entry Point & Gateway
```

---

## 🚀 Getting Started

### 📋 ข้อกำหนดเบื้องต้น (Prerequisites)
- [Go](https://go.dev/dl/) เวอร์ชัน 1.22 ขึ้นไป
- [Docker](https://www.docker.com/) & Docker Compose

### ⚙️ วิธีการติดตั้งและรันเซิร์ฟเวอร์

1. สตาร์ท Infrastructure ทั้งหมด (PostgreSQL, Redis, Prometheus, Grafana):
   ```bash
   docker compose up -d
   ```

2. สตาร์ท Notification Microservice (gRPC Server):
   ```bash
   go run cmd/notification-service/main.go
   ```

3. สตาร์ท Task API Gateway (HTTP REST API + gRPC Client):
   ```bash
   go run main.go
   ```

4. จุดเชื่อมต่อบริการต่างๆ:
   - 🚀 **REST API Server**: `http://localhost:8080`
   - 📖 **Swagger UI Docs**: `http://localhost:8080/swagger/index.html`
   - 📊 **Prometheus Metrics**: `http://localhost:8080/metrics`
   - 📈 **Grafana Dashboard**: `http://localhost:3000` *(Login: admin / admin)*
   - 📬 **Notification Microservice (gRPC)**: `localhost:50051`

---

## 🧪 การรัน Unit Tests
```bash
go test -v -cover ./...
```
