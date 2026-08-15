# 📝 Task Management REST API (Go Zero to Hero — Project 1)

RESTful API สำหรับจัดการรายการงาน (Task Management) ที่พัฒนาด้วยภาษา **Go (Golang)** โดยใช้ Standard Library (`net/http`) ล้วน ๆ เพื่อเรียนรู้พื้นฐานและโครงสร้างสถาปัตยกรรมโค้ดที่สะอาด (Clean Layered Architecture)

---

## ✨ Features

- ⚡ **CRUD Operations**: สร้าง (Create), เรียกดู (Read), แก้ไข (Update), และลบ (Delete) รายการ Task
- 🔒 **Thread-Safe In-Memory Storage**: จัดการข้อมูลในหน่วยความจำอย่างปลอดภัยในสภาวะ Concurrency ด้วย `sync.RWMutex`
- 🛣️ **Modern Go 1.22+ Routing**: ใช้ `http.NewServeMux` พร้อม Path Parameters (`/tasks/{id}`)
- 🛡️ **Custom Middleware**:
  - `Logging`: บันทึก Request Method, Path, Status Code และ Latency
  - `JSONContentType`: ตั้งค่า `Content-Type: application/json` ให้อัตโนมัติทุก Request
- 🎯 **Clean Architecture**: แบ่งแยกหน้าที่ของโค้ดชัดเจน (`models`, `repository`, `handlers`, `middleware`)

---

## 📁 Project Structure

```
1_task-app/
├── models/
│   └── task.go             # Data Models, DTOs และ Validation
├── repository/
│   └── task_repository.go  # Interface และ In-Memory Storage (Thread-safe)
├── handlers/
│   └── task_handler.go     # HTTP Handlers และการแปลง JSON
├── middleware/
│   └── middleware.go       # Logging และ Response Header Middleware
├── go.mod                  # Go Module Definition
├── main.go                 # Entry Point และการเชื่อมต่อ Routing
└── README.md               # คู่มือการใช้งานโปรเจกต์
```

---

## 🚀 Getting Started

### 📋 ข้อกำหนดเบื้องต้น (Prerequisites)
- [Go](https://go.dev/dl/) เวอร์ชัน 1.22 ขึ้นไป

### ⚙️ การติดตั้งและรันเซิร์ฟเวอร์
1. Clone หรือเปิดโฟลเดอร์โปรเจกต์:
   ```bash
   cd 1_task-app
   ```

2. รันแอปพลิเคชัน:
   ```bash
   go run main.go
   ```

3. เซิร์ฟเวอร์จะเริ่มต้นทำงานที่: `http://localhost:8080`

---

## 📡 API Endpoints

| Method | Endpoint | คำอธิบาย | Status Code |
| :--- | :--- | :--- | :--- |
| `GET` | `/health` | ตรวจสอบสถานะการทำงานของ Server | `200 OK` |
| `GET` | `/tasks` | ดึงรายการ Task ทั้งหมด | `200 OK` |
| `GET` | `/tasks/{id}` | ดึงรายละเอียด Task ตาม ID | `200 OK` / `404 Not Found` |
| `POST` | `/tasks` | สร้าง Task ใหม่ | `201 Created` / `400 Bad Request` |
| `PUT` | `/tasks/{id}` | แก้ไขข้อมูล Task ตาม ID | `200 OK` / `400 Bad Request` / `404 Not Found` |
| `DELETE` | `/tasks/{id}` | ลบ Task ตาม ID | `200 OK` / `404 Not Found` |

---

## 🧪 ตัวอย่างการทดสอบ API

### 1. ดึงรายการ Task ทั้งหมด (GET `/tasks`)
```bash
curl -X GET http://localhost:8080/tasks
```
**PowerShell:**
```powershell
Invoke-RestMethod -Uri 'http://localhost:8080/tasks'
```

### 2. สร้าง Task ใหม่ (POST `/tasks`)
```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "เรียนรู้ Go REST API", "description": "ศึกษาเรื่อง Go Structs และ Handlers"}'
```
**PowerShell:**
```powershell
$body = @{
    title = "เรียนรู้ Go REST API"
    description = "ศึกษาเรื่อง Go Structs และ Handlers"
} | ConvertTo-Json

Invoke-RestMethod -Uri 'http://localhost:8080/tasks' -Method Post -Body $body -ContentType 'application/json'
```

### 3. อัปเดตสถานะ Task (PUT `/tasks/1`)
```bash
curl -X PUT http://localhost:8080/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"completed": true}'
```
**PowerShell:**
```powershell
$update = @{ completed = $true } | ConvertTo-Json
Invoke-RestMethod -Uri 'http://localhost:8080/tasks/1' -Method Put -Body $update -ContentType 'application/json'
```

### 4. ลบ Task (DELETE `/tasks/1`)
```bash
curl -X DELETE http://localhost:8080/tasks/1
```
**PowerShell:**
```powershell
Invoke-RestMethod -Uri 'http://localhost:8080/tasks/1' -Method Delete
```

---

## 🧠 Key Concepts Learned

- **Go Structs & Tags**: การใช้ Struct Tags (`json:"id"`) เพื่อแปลงข้อมูลระหว่าง Go และ JSON
- **Pointers**: การใช้ Pointer (`*string`, `*bool`) ใน DTO สำหรับทำ Partial Update
- **Implicit Interface**: การออกแบบ Interface สำหรับ Dependency Injection
- **Concurrency & Goroutines**: การใช้ `sync.RWMutex` และ `defer` ป้องกัน Data Race
- **Go Standard Library**: การใช้ `net/http` และ `json` โดยไม่ต้องพึ่งพา 3rd-party Frameworks
