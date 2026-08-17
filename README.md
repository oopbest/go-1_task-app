# 📝 Task Management REST API (Go Zero to Hero)

RESTful API สำหรับจัดการงานแบบ **Multi-User** ระดับ **Enterprise-Grade** พัฒนาด้วย **Gin Web Framework** เชื่อมต่อฐานข้อมูล **PostgreSQL** เสริมความเร็วระดับ **Microseconds ด้วย Redis Cache**, ระบบประมวลผลเบื้องหลัง **Background Worker Pool (Goroutines & Channels)**, ความปลอดภัย **JWT Authentication**, การจัดการ Schema ด้วย **Database Migrations**, เอกสาร **Swagger (OpenAPI) Interactive UI**, และออกแบบตามหลักการ **Clean Layered Architecture (Decorator Pattern)** 🚀⚡📖

---

## ✨ Features

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
- 🐳 **Docker Compose & Multi-Stage Build**: รัน PostgreSQL + Redis และบิลด์ Static Binary บน Alpine Linux ขนาดเล็กลงเหลือเพียง **24.5 MB**

---

## 📁 Project Structure

```
1_task-app/
├── docs/                   # Swagger / OpenAPI Generated Files
├── migrations/             # Database Schema Migrations (.up.sql / .down.sql)
├── models/
│   ├── user.go             # Data Models สำหรับ User & Auth DTOs
│   ├── task.go             # Data Models สำหรับ Task & DTOs
│   └── pagination.go       # DTOs สำหรับ Pagination, Search & Filter
├── repository/
│   ├── user_repository.go  # PostgreSQL User Repository
│   ├── task_repository.go  # TaskRepository Interface & In-Memory Storage
│   ├── postgres_task_repository.go # PostgreSQL Dynamic Query Repository
│   ├── cached_task_repository.go   # Redis Caching Decorator Repository
│   └── migrations.go       # Migration Runner
├── workers/
│   └── worker_pool.go      # Background Worker Pool (Goroutines & Channels)
├── handlers/
│   ├── auth_handler.go     # Gin HTTP Handlers สำหรับ Register & Login
│   └── task_handler.go     # Gin HTTP Handlers สำหรับ Tasks CRUD & Worker Enqueue
├── middleware/
│   ├── gin_auth.go         # Gin JWT Authentication Middleware
│   ├── auth.go             # Context Auth Helper
│   └── middleware.go       # Logging & JSON Content-Type Middleware
├── utils/
│   └── auth.go             # Bcrypt Password Hashing & JWT Helpers
├── docker-compose.yml      # PostgreSQL & Redis Services
├── Dockerfile              # Multi-Stage Dockerfile
├── .dockerignore           # Docker ignore rules
├── go.mod                  # Go Module Definition
└── main.go                 # Application Entry Point, Gin Router & Swagger
```

---

## 🚀 Getting Started

### 📋 ข้อกำหนดเบื้องต้น (Prerequisites)
- [Go](https://go.dev/dl/) เวอร์ชัน 1.22 ขึ้นไป
- [Docker](https://www.docker.com/) & Docker Compose

### ⚙️ วิธีการติดตั้งและรันเซิร์ฟเวอร์

1. สตาร์ท PostgreSQL Database & Redis Cache:
   ```bash
   docker compose up -d
   ```

2. รันแอปพลิเคชัน Go:
   ```bash
   go run main.go
   ```

3. เซิร์ฟเวอร์จะเริ่มต้นทำงานที่: `http://localhost:8080`
4. เปิดหน้าเว็บ Swagger UI เพื่อทดสอบ API: **`http://localhost:8080/swagger/index.html`**

---

## 📡 API Endpoints

### 🔓 Public Endpoints (ไม่ต้องใช้ Token)
| Method | Endpoint | คำอธิบาย |
| :--- | :--- | :--- |
| `GET` | `/health` | ตรวจสอบสถานะ Server, Database, Cache, และ Workers |
| `GET` | `/swagger/*any` | เอกสาร Interactive Swagger API Documentation |
| `POST` | `/auth/register` | สมัครสมาชิกใหม่ (รับ `email`, `password`) |
| `POST` | `/auth/login` | เข้าสู่ระบบเพื่อรับ JWT Token |

### 🔒 Protected Endpoints (ต้องแนบ `Authorization: Bearer <token>`)
| Method | Endpoint | คำอธิบาย |
| :--- | :--- | :--- |
| `GET` | `/tasks` | ดึงรายการ Task (รองรับ `page`, `limit`, `search`, `completed`, `sort`, `order`) |
| `GET` | `/tasks/{id}` | ดึงรายละเอียด Task ตาม ID (ดึงจาก Redis Cache) |
| `POST` | `/tasks` | สร้าง Task ใหม่ + ล้าง Cache + โยนเข้า Worker Pool |
| `PUT` | `/tasks/{id}` | แก้ไข Task + ล้าง Cache + โยนเข้า Worker Pool |
| `DELETE` | `/tasks/{id}` | ลบ Task + ล้าง Cache + โยนเข้า Worker Pool |

---

## 🧪 การรัน Unit Tests
```bash
go test -v -cover ./...
```
