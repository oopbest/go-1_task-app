# 📝 Task Management Enterprise SaaS Platform (Go Zero to Hero)

ระบบจัดการงานระดับ **Enterprise Full-Stack SaaS Platform** ขับเคลื่อนด้วย **Next.js 15 (App Router)** ฝั่งหน้าบ้าน เชื่อมต่อกับ **Go (Gin Web Framework)** ฝั่งหลังบ้าน เสริมความเร็วด้วย **PostgreSQL + Redis In-Memory Cache (~500µs)**, ประมวลผลเบื้องหลังด้วย **Go Background Worker Pools**, สถาปัตยกรรม **Distributed Microservices ด้วย gRPC (HTTP/2 Binary Protocol)**, ความปลอดภัย **JWT Authentication**, การบริหารจัดการ Schema ด้วย **Database Migrations**, เอกสาร **Swagger (OpenAPI) Interactive UI**, และระบบเฝ้าระวังครบวงจร **Observability (Structured Logging `slog` + Prometheus Metrics + Grafana Dashboards)** 🚀🎨💻⚡🌐📊📈

---

## ✨ Full-Stack Features

- 💻 **Modern Next.js 15 SaaS Frontend (`/frontend`)**:
  - พัฒนาด้วย **Next.js 15 (Turbopack)**, **React 19**, **TypeScript**, และ **Tailwind CSS**
  - **Optimistic UI Updates (TanStack Query v5)**: ติ๊กเปลี่ยนสถานะงานแล้ว UI ตอบสนองทันทีใน **0 ms**
  - **Glassmorphic Dark UI**: ดีไซน์ระดับ World-Class พร้อม Glowing Background Orbs และ Sonner Toast Notifications
  - **Live Debounced Search & Filter**: ค้นหางานแบบ Real-Time กรองสถานะงาน และสลับการเรียงลำดับ เก่า ➡️ ใหม่
  - **Session Guard**: จัดการสิทธิ์การเข้าถึงหน้าเว็บด้วย Auth Context & Cookies
- 🌐 **Microservices Architecture with gRPC & Protocol Buffers**:
  - แยก Service ออกเป็น **Task API Gateway (:8080)** และ **Notification Microservice (:50051)**
  - สื่อสารระหว่าง Services ข้ามเครือข่ายด้วย **gRPC บน HTTP/2 Binary Protocol** ที่เร็วกว่า JSON 5–10 เท่า
- 📊 **Full-Stack Observability & Monitoring**:
  - **Structured JSON Logging (`log/slog`)**: บันทึก Log ทุก Request ในรูปแบบ JSON ตามมาตรฐาน Go 1.21+
  - **Prometheus Metrics Exporter**: เก็บสถิติ Request Count, Latency Histogram, และ Worker Jobs ผ่าน Endpoint `/metrics`
  - **Grafana Live Dashboards**: แสดงผลกราฟสถิติ Real-Time บนเว็บ `http://localhost:3000` (admin/admin)
- 🏎️ **Gin Web Framework & High Performance Routing**:
  - Radix Tree Router ความเร็วสูงและใช้ Memory น้อยที่สุด พร้อม Route Groups (`/auth`, `/tasks`) และ CORS Middleware
- 📖 **Swagger / OpenAPI Interactive UI**:
  - สร้างเอกสาร API Documentation อัตโนมัติด้วย `swaggo/swag` เปิดทดสอบที่ `http://localhost:8080/swagger/index.html`
- ⚡ **Background Worker Pool (Go Concurrency & Channels)**:
  - ประมวลผลงานเบื้องหลังแบบ Asynchronous (3 Workers, Queue Size 100) พร้อม **Zero-Data-Loss Graceful Shutdown**
- 🏎️ **Redis In-Memory Caching (Sub-Millisecond Latency)**:
  - **Cache-Aside Pattern**: ตรวจสอบ Redis Cache ก่อนดึง PostgreSQL ลด Latency เหลือเพียง **~500 µs (ไมโครวินาที)**
- 🔐 **User Authentication (JWT & Bcrypt)**:
  - สมัครสมาชิก (`/auth/register`) และเข้าสู่ระบบ (`/auth/login`) ด้วย `bcrypt` และ `JWT`
- 👥 **Multi-User Task Isolation**: แยกข้อมูลงานของผู้ใช้แต่ละคนอย่างเด็ดขาด (Row-level Security)
- 🗄️ **Database Migrations & Indexing (`golang-migrate`)**:
  - จัดการ Schema ด้วย `.up.sql` / `.down.sql` พร้อม B-Tree Indexes
- 🐳 **Docker Compose & Multi-Stage Build**: รัน PostgreSQL + Redis + Prometheus + Grafana ครบวงจร

---

## 📁 Project Structure

```
1_task-app/
├── frontend/                 # [NEW] Next.js 15 Full-Stack SaaS Web App (:3001)
│   ├── src/
│   │   ├── app/              # Next.js App Router (login, register, dashboard)
│   │   ├── context/          # Auth Context & Session Guard
│   │   ├── providers/        # React Query & Sonner Toaster Provider
│   │   ├── lib/              # API Client & Utilities
│   │   └── types/            # TypeScript Interfaces
│   └── package.json
├── cmd/
│   └── notification-service/ # Notification Microservice (gRPC Server on :50051)
├── proto/
│   ├── notification.proto    # Protobuf Contract Definition
│   └── notification/         # Generated Go Protobuf & gRPC Stubs
├── docs/                     # Swagger / OpenAPI Generated Files
├── metrics/                  # Prometheus Metrics Collector & Middleware
├── migrations/               # Database Schema Migrations (.up.sql / .down.sql)
├── models/                   # Data Models & DTOs
├── repository/               # PostgreSQL, In-Memory & Redis Cached Repositories
├── workers/                  # Background Worker Pool & gRPC Client
├── handlers/                 # Gin HTTP Handlers
├── middleware/               # CORS, Auth, Logger Middlewares
├── utils/                    # Password Hashing & JWT Helpers
├── docker-compose.yml        # Postgres + Redis + Prometheus + Grafana
├── Dockerfile                # Multi-Stage Dockerfile
├── go.mod                    # Go Module Definition
└── main.go                   # Go Backend API Gateway (:8080)
```

---

## 🚀 Getting Started

### 📋 ข้อกำหนดเบื้องต้น (Prerequisites)
- [Go](https://go.dev/dl/) เวอร์ชัน 1.22 ขึ้นไป
- [Node.js](https://nodejs.org/) เวอร์ชัน 18+ (สำหรับ Frontend)
- [Docker](https://www.docker.com/) & Docker Compose

### ⚙️ วิธีการติดตั้งและรัน Full-Stack System

1. สตาร์ท Infrastructure ทั้งหมด (PostgreSQL, Redis, Prometheus, Grafana):
   ```bash
   docker compose up -d
   ```

2. สตาร์ท Notification Microservice (gRPC Server):
   ```bash
   go run cmd/notification-service/main.go
   ```

3. สตาร์ท Go Backend API Gateway:
   ```bash
   go run main.go
   ```

4. สตาร์ท Next.js 15 Frontend:
   ```bash
   cd frontend
   npm run dev
   ```

5. จุดเชื่อมต่อบริการทั้งหมด:
   - 💻 **Frontend Web App**: `http://localhost:3001`
   - 🚀 **REST API Server**: `http://localhost:8080`
   - 📖 **Swagger UI Docs**: `http://localhost:8080/swagger/index.html`
   - 📊 **Prometheus Metrics**: `http://localhost:8080/metrics`
   - 📈 **Grafana Dashboard**: `http://localhost:3000` *(Login: admin / admin)*
   - 📬 **Notification Microservice (gRPC)**: `localhost:50051`

---

## 🧪 การรัน Unit Tests
```bash
# Go Backend Tests
go test -v -cover ./...

# Frontend Lint & Type Checks
cd frontend && npm run lint && npx tsc --noEmit
```
