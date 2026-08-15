# Walkthrough 2: Production-Ready Upgrades for Task Management API

ในส่วนนี้เป็นการต่อยอด **Task Management REST API** ให้พร้อมสำหรับการนำไปใช้งานจริงบนระดับ **Production** โดยเพิ่มระบบ Automated Testing, การปิดเซิร์ฟเวอร์อย่างปลอดภัย (Graceful Shutdown), และการทำ Containerization ด้วย Docker Multi-Stage Build 🚀

---

## 🎯 สิ่งที่ได้ทำและพัฒนาเพิ่มเติม

```
1_task-app/
├── .dockerignore              # ป้องกันไฟล์ไม่จำเป็นเข้า Docker Context
├── Dockerfile                 # Multi-Stage Build สำหรับสร้าง Container ขนาดจิ๋ว
├── handlers/
│   ├── task_handler.go
│   └── task_handler_test.go   # [NEW] Unit Test สำหรับ HTTP Handlers
├── repository/
│   ├── task_repository.go
│   └── task_repository_test.go# [NEW] Unit Test สำหรับ CRUD Operations
└── main.go                    # [UPDATE] เพิ่ม Graceful Shutdown & Port จาก ENV
```

---

## 🧪 1.1 Unit Testing & Code Coverage

สร้าง Unit Test ด้วย Standard Library (`testing` และ `net/http/httptest`) โดยไม่ต้องพึ่งพา 3rd-party Framework:

- **Repository Test (`task_repository_test.go`)**: ทดสอบ CRUD Functionality ทั้ง Create, GetByID, Update, Delete
- **Handler Test (`task_handler_test.go`)**: ทดสอบ HTTP Endpoints โดยจำลอง Request และ ResponseRecorder

### 📊 ผลการรันคำสั่ง `go test -v -cover ./...`:
```text
=== RUN   TestTaskHandler_GetAllTasks
--- PASS: TestTaskHandler_GetAllTasks (0.00s)
=== RUN   TestTaskHandler_CreateTask
--- PASS: TestTaskHandler_CreateTask (0.00s)
PASS
coverage: 18.0% of statements
ok      github.com/oopbest/task-app/handlers    coverage: 18.0% of statements

=== RUN   TestMemoryTaskRepository_CRUD
--- PASS: TestMemoryTaskRepository_CRUD (0.00s)
PASS
coverage: 78.0% of statements
ok      github.com/oopbest/task-app/repository  coverage: 78.0% of statements
```

---

## 🛡️ 1.2 Graceful Shutdown & Environment Configuration

อัปเกรด `main.go` ให้รองรับการทำงานในสภาพแวดล้อมจริง:

1. **Environment Variables**: อ่านค่า `PORT` ผ่าน `os.Getenv("PORT")` โดยมี fallback เป็น `8080`
2. **Goroutine Execution**: รัน `server.ListenAndServe()` ใน Background Goroutine เพื่อไม่ให้บล็อก Main Thread
3. **Signal Trapping**: ใช้ Go Channel `chan os.Signal` ร่วมกับ `signal.Notify` ดักจับสัญญาณ `Ctrl+C` (`SIGINT`) และ `SIGTERM` จาก Docker/Kubernetes
4. **Clean Exit**: ใช้ `context.WithTimeout(..., 5*time.Second)` และ `server.Shutdown(ctx)` เพื่อรอให้ Request ที่กำลังประมวลผลอยู่ทำงานเสร็จก่อนปิด

### 📝 ผลลัพธ์เมื่อกดปิดเซิร์ฟเวอร์ (Ctrl+C):
```text
==================================================
🚀 Server is running on http://localhost:8080
==================================================
2026/08/15 09:20:18 🛑 Shutdown signal received, shutting down gracefully...
2026/08/15 09:20:18 ✅ Server exited cleanly
```

---

## 🐳 1.3 Docker Multi-Stage Build

สร้าง Docker Image โดยใช้เทคนิค **Multi-Stage Build**:

- **Stage 1 (Builder)**: ใช้ `golang:1.24-alpine` คอมไพล์ Source Code เป็น Static Binary Linux ด้วยคำสั่ง:
  ```bash
  CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o task-app .
  ```
- **Stage 2 (Runtime)**: สลับไปใช้ Base Image `alpine:latest` ขนาดเล็ก แล้วดึงเฉพาะ Binary `task-app` มาใช้งาน
- **ผลลัพธ์**: ขนาด Image รวมทั้งหมดเล็กเพียง **24.5 MB** (Content Size เพียง 6.76 MB)

```text
IMAGE         ID             DISK USAGE   CONTENT SIZE
task-app:v1   8d8f137ba218   24.5MB       6.76MB
```

---

## 🌐 1.4 การทดสอบ Container จริง (Live Verification)

รัน Container ผ่านคำสั่ง:
```bash
docker run --rm -p 8080:8080 --name my-task-app task-app:v1
```

- **Health Check (`GET /health`)**: ทดสอบผ่าน Browser ส่งคืน JSON Status `ok` พร้อมเวลา
- **Get Tasks (`GET /tasks`)**: ทดสอบผ่าน Postman ได้รับรายการ Task ถูกต้องจากภายใน Container

---

## 🧠 สรุป Concept สำคัญของ Section นี้

| หัวข้อ | Concept และ Best Practice ใน Go |
| :--- | :--- |
| **Go Testing** | การใช้ `testing.T`, `httptest.NewRecorder()`, และคำสั่ง `go test -cover` |
| **Go Channels & Signals** | การใช้ Channel ดักจับ OS Signals เพื่อควบคุม Lifecycle ของโปรแกรม |
| **Context Cancellation** | การใช้ `context.WithTimeout` ควบคุม Deadline ของการปิด Server |
| **Static Binary Compilation** | การใช้ `CGO_ENABLED=0` และ `-ldflags="-s -w"` ตัด Debug info ลดขนาดไฟล์ |
| **Multi-Stage Docker** | การแยก Build Environment ออกจาก Runtime Environment เพื่อความปลอดภัยและประหยัดพื้นที่ |
