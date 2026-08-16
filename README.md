# 📝 Task Management REST API (Go Zero to Hero)

RESTful API สำหรับจัดการงานแบบ **Multi-User** ที่พัฒนาด้วยภาษา **Go (Golang)** เชื่อมต่อฐานข้อมูล **PostgreSQL** พร้อมระบบรักษาความปลอดภัยด้วย **JWT (JSON Web Token)** และออกแบบตามหลักการ **Clean Layered Architecture** 🚀

---

## ✨ Features

- 🔐 **User Authentication (JWT & Bcrypt)**:
  - สมัครสมาชิก (`/auth/register`) และเข้าสู่ระบบ (`/auth/login`)
  - เข้ารหัสผ่านอย่างปลอดภัยด้วย `bcrypt`
  - ยืนยันตัวตนด้วย `JWT (JSON Web Token)` มีอายุ 24 ชั่วโมง
- 👥 **Multi-User Task Management**: แยกข้อมูลงานของผู้ใช้แต่ละคนอย่างเด็ดขาด (งานใครงานมัน)
- ⚡ **Full CRUD Operations**: สร้าง (Create), ดูรายการ (Read), แก้ไข (Update), และลบ (Delete) งานของตัวเอง
- 🐘 **PostgreSQL Database & Connection Pool**: จัดการ Connection Pool ด้วย `database/sql` และรันบน Docker Compose
- 🛡️ **Middleware Pipeline**:
  - `AuthMiddleware`: ตรวจสอบ Bearer Token และฝาก `user_id` เข้า Request Context
  - `Logging`: บันทึก Method, Path, Status Code, และ Latency
  - `JSONContentType`: ตั้งค่า Header `Content-Type: application/json` อัตโนมัติ
- 🧪 **Automated Unit Testing**: เขียน Test ด้วย `testing` และ `net/http/httptest`
- 🛑 **Graceful Shutdown**: ดักจับ OS Signals ปิดเซิร์ฟเวอร์และเคลียร์ Connections อย่างปลอดภัย
- 🐳 **Docker Multi-Stage Build**: บิลด์เป็น Static Binary บน Alpine Linux ขนาดเล็กลงเหลือเพียง **24.5 MB**

---

## 📁 Project Structure

```
1_task-app/
├── models/
│   ├── user.go             # Data Models สำหรับ User & Auth DTOs
│   └── task.go             # Data Models สำหรับ Task & DTOs
├── repository/
│   ├── user_repository.go  # Interface & PostgreSQL User Storage
│   ├── task_repository.go  # Interface & In-Memory Task Storage
│   └── postgres_task_repository.go # PostgreSQL Task Storage
├── handlers/
│   ├── auth_handler.go     # HTTP Handlers สำหรับ Register & Login
│   └── task_handler.go     # HTTP Handlers สำหรับ Tasks CRUD
├── middleware/
│   ├── auth.go             # JWT Authentication Middleware & Context Helper
│   └── middleware.go       # Logging & JSON Content-Type Middleware
├── utils/
│   └── auth.go             # Bcrypt Password Hashing & JWT Helpers
├── docker-compose.yml      # PostgreSQL Database Service
├── Dockerfile              # Multi-Stage Dockerfile
├── .dockerignore           # Docker ignore rules
├── go.mod                  # Go Module Definition
└── main.go                 # Application Entry Point & Routes
```

---

## 🚀 Getting Started

### 📋 ข้อกำหนดเบื้องต้น (Prerequisites)
- [Go](https://go.dev/dl/) เวอร์ชัน 1.22 ขึ้นไป
- [Docker](https://www.docker.com/) & Docker Compose

### ⚙️ วิธีการติดตั้งและรันเซิร์ฟเวอร์

1. สตาร์ท PostgreSQL Database:
   ```bash
   docker compose up -d
   ```

2. รันแอปพลิเคชัน Go:
   ```bash
   go run main.go
   ```

3. เซิร์ฟเวอร์จะเริ่มต้นทำงานที่: `http://localhost:8080`

---

## 📡 API Endpoints

### 🔓 Public Endpoints (ไม่ต้องใช้ Token)
| Method | Endpoint | คำอธิบาย |
| :--- | :--- | :--- |
| `GET` | `/health` | ตรวจสอบสถานะการทำงานของ Server |
| `POST` | `/auth/register` | สมัครสมาชิกใหม่ (รับ `email`, `password`) |
| `POST` | `/auth/login` | เข้าสู่ระบบเพื่อรับ JWT Token |

### 🔒 Protected Endpoints (ต้องแนบ `Authorization: Bearer <token>`)
| Method | Endpoint | คำอธิบาย |
| :--- | :--- | :--- |
| `GET` | `/tasks` | ดึงรายการ Task ทั้งหมดเฉพาะของตัวเอง |
| `GET` | `/tasks/{id}` | ดึงรายละเอียด Task ตาม ID |
| `POST` | `/tasks` | สร้าง Task ใหม่ผูกกับ User ตัวเอง |
| `PUT` | `/tasks/{id}` | แก้ไขข้อมูล Task ของตัวเอง |
| `DELETE` | `/tasks/{id}` | ลบ Task ของตัวเอง |

---

## 🧪 การรัน Unit Tests
```bash
go test -v -cover ./...
```

---

## 🐳 การรันผ่าน Docker
```bash
# Build Image
docker build -t task-app:v1 .

# Run Container
docker run --rm -p 8080:8080 task-app:v1
```
