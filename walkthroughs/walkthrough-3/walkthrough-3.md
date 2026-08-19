# Walkthrough 3: PostgreSQL Database Integration (Clean Architecture & Persistence)

ในส่วนนี้เป็นการยกระดับ **Task Management REST API** จากการเก็บข้อมูลชั่วคราวใน Memory สู่การเชื่อมต่อฐานข้อมูล **PostgreSQL** จริง โดยอาศัยพลังของ **Interface & Clean Architecture** ที่ทำให้เราไม่ต้องแก้ไขโค้ดในส่วน Handler เลยแม้แต่บรรทัดเดียว! 🚀🐘

---

## 🏗️ โครงสร้างไฟล์ที่เพิ่มขึ้นในโปรเจกต์

```
1_task-app/
├── docker-compose.yml                 # [NEW] จัดการ PostgreSQL Container
├── repository/
│   ├── task_repository.go             # Interface กลาง (Contract)
│   ├── memory_task_repository.go      # In-Memory Storage (เดิม)
│   └── postgres_task_repository.go    # [NEW] PostgreSQL Storage พร้อม Raw SQL
└── main.go                            # [UPDATE] เชื่อมต่อ Connection Pool & สลับ Repo
```

---

## 🐘 3.1 การตั้งค่า PostgreSQL ด้วย Docker Compose

สร้างไฟล์ `docker-compose.yml` เพื่อรัน PostgreSQL 16 บน Container:

- **Image**: `postgres:16-alpine`
- **Database**: `taskdb`
- **User / Password**: `postgres` / `mysecretpassword`
- **Port**: `5432:5432`
- **Volume**: `postgres_data` เพื่อเก็บข้อมูลให้อยู่ถาวรแม้ Container จะถูกหยุดทำงาน

```bash
docker compose up -d
```

---

## 💾 3.2 การสร้าง `PostgresTaskRepository`

พัฒนา `PostgresTaskRepository` ที่ implement ตาม `TaskRepository` interface:

### 1. Auto Table Migration
สร้างตาราง `tasks` อัตโนมัติเมื่อ Repository ถูกเริ่มต้น:
```sql
CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### 2. Parameterized Queries (ป้องกัน SQL Injection 100%)
- **Create**: `INSERT INTO tasks (title, description) VALUES ($1, $2) RETURNING ...`
- **GetByID**: `SELECT ... FROM tasks WHERE id = $1`
- **Update**: `UPDATE tasks SET title=$1, description=$2, completed=$3 WHERE id=$4 RETURNING ...`
- **Delete**: `DELETE FROM tasks WHERE id = $1`

### 3. Go Linter Best Practices (`sqlrowserr`)
ตรวจสอบ `rows.Err()` เสมอหลังจบลูป `for rows.Next()` เพื่อดักจับ Error ที่อาจเกิดขึ้นระหว่างการสตรีมข้อมูลจากฐานข้อมูล:
```go
for rows.Next() {
    // scan ข้อมูลแต่ละแถว
}
if err := rows.Err(); err != nil {
    return []models.Task{}
}
```

---

## 🔌 3.3 การจัดการ Database Connection Pool ใน `main.go`

ใน `main.go` เราตั้งค่า Connection Pool ของ Go `database/sql` ให้มีประสิทธิภาพระดับ Production:

```go
db, err := sql.Open("postgres", dbURL)

// ตั้งค่า Connection Pool
db.SetMaxOpenConns(25)                  // จำนวน Connection สูงสุดที่เปิดพร้อมกัน
db.SetMaxIdleConns(25)                  // จำนวน Connection ที่พักรอไว้ใน Pool
db.SetConnMaxLifetime(5 * time.Minute)  // อายุสูงสุดของแต่ละ Connection
```

---

## 🧪 3.4 ผลการทดสอบ (Verification Results)

### 1. Log การเชื่อมต่อ Server & Database:
```text
2026/08/15 10:42:54 🐘 Connected to PostgreSQL database successfully!
==================================================
🚀 Task API (PostgreSQL) is running on http://localhost:8080
==================================================
```

### 2. ผลการทดสอบผ่าน Postman (`GET /tasks`):
```json
[
  {
    "id": 2,
    "title": "just do it",
    "description": "desc do it",
    "completed": false,
    "created_at": "2026-08-15T03:53:31.995628Z"
  }
]
```

### 3. Data Persistence Test:
เมื่อทดลองกด `Ctrl+C` ปิด Server แล้วรัน `go run main.go` ใหม่อีกครั้ง ข้อมูลทั้งหมดใน PostgreSQL **ยังคงอยู่ครบถ้วน 100% ไม่สูญหาย** 🎉

---

## 🧠 สรุป Concept สำคัญของ Section นี้

| หัวข้อ | Concept สำคัญในภาษา Go |
| :--- | :--- |
| **Interface Polymorphism** | การสลับ Data Source จาก Memory เป็น Postgres โดยไม่ต้องแก้ Handler |
| **Blank Identifier Import** | การใช้ `_ "github.com/lib/pq"` เพื่อลงทะเบียน Driver กับ `database/sql` |
| **Connection Pooling** | การทำงานของ `*sql.DB` ที่จัดการ Pool อัตโนมัติ ปลอดภัยต่อ Concurrency |
| **QueryRow vs Query vs Exec** | การเลือกใช้ฟังก์ชัน SQL ให้เหมาะกับชนิดคำสั่งและจำนวนแถวที่ต้องการ |
| **Scanning & Type Mapping** | การใช้ `rows.Scan()` จับคู่ Type จาก Postgres (เช่น `TIMESTAMPTZ` -> `time.Time`) |
