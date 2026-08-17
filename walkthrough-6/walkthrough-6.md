# Walkthrough 6: Caching & High Performance with Redis (Cache-Aside & Decorator Pattern)

ในส่วนนี้เป็นการยกระดับประสิทธิภาพของ **Task Management REST API** ให้มีความเร็วระดับ **Microseconds** โดยนำ **Redis In-Memory Database** มาทำหน้าที่เป็น Caching Layer ด่านหน้าก่อนเข้าถึง PostgreSQL พร้อมใช้ **Decorator Pattern** เพื่อรักษาหลักการ Clean Architecture 🏎️⚡

---

## 🏗️ โครงสร้างสถาปัตยกรรม Caching Layer

```
1_task-app/
├── docker-compose.yml              # [UPDATE] เพิ่ม Redis Service (redis:7-alpine พอร์ต 6379)
├── repository/
│   └── cached_task_repository.go   # [NEW] Decorator Caching Layer (Cache-Aside & Eviction)
├── main.go                         # [UPDATE] เชื่อมต่อ Redis และห่อหุ้ม Postgres ด้วย Cached Repo
└── go.mod                          # [UPDATE] เพิ่ม github.com/redis/go-redis/v9
```

---

## 🔍 6.1 ลำดับการทำงานของ Cache-Aside & Invalidation Pattern

```mermaid
sequenceDiagram
    autonumber
    actor Client
    participant Handler as TaskHandler
    participant Cache as CachedTaskRepository (Redis)
    participant DB as PostgresTaskRepository (PostgreSQL)

    Note over Client,DB: 1. อ่านข้อมูล (Read Operation)
    Client->>Handler: GET /tasks
    Handler->>Cache: GetAll(userID, filter)
    
    alt ⚡ Cache Hit (มีข้อมูลใน Redis)
        Cache-->>Handler: คืนค่า JSON จาก RAM ทันที (502 µs)
        Handler-->>Client: 200 OK
    else 🐢 Cache Miss (ไม่มีข้อมูลใน Redis)
        Cache->>DB: Query จาก PostgreSQL (4.8 ms)
        DB-->>Cache: รายการ Tasks
        Cache->>Cache: บันทึกลง Redis พร้อม TTL = 5 นาที
        Cache-->>Handler: รายการ Tasks
        Handler-->>Client: 200 OK
    end

    Note over Client,DB: 2. เขียน/แก้ไขข้อมูล (Write Operation)
    Client->>Handler: POST /tasks (สร้างงานใหม่)
    Handler->>Cache: Create(input, userID)
    Cache->>DB: บันทึกลง PostgreSQL
    DB-->>Cache: งานใหม่
    Cache->>Cache: 🧹 ลบ Cache เก่าของ User ทิ้งทันที (Cache Eviction)
    Cache-->>Handler-->>Client: 201 Created
```

---

## ⚡ 6.2 จุดเด่นของสถาปัตยกรรม (Design Patterns)

### 1. Decorator Pattern (Open-Closed Principle)
`CachedTaskRepository` ถูกสร้างขึ้นโดย Implement `repository.TaskRepository` ทำให้สามารถนำไป **"สวมทับ (Wrap)"** PostgreSQL Repository ได้ทันที โดย **ไม่ต้องแก้ไข Business Logic ใน Handler เลยแม้แต่บรรทัดเดียว!**

```go
postgresTaskRepo := repository.NewPostgresTaskRepository(db)
cachedTaskRepo := repository.NewCachedTaskRepository(postgresTaskRepo, rdb, 5*time.Minute)
taskHandler := handlers.NewTaskHandler(cachedTaskRepo)
```

### 2. Cache Invalidation Strategy (Zero Stale Data)
เมื่อเกิดการ `Create`, `Update`, หรือ `Delete` ระบบจะทำการลบ Cache Key ของ User คนนั้นทิ้งทันที (`tasks:u:<userID>:*`) เพื่อรับประกันว่าในการเรียก `GET` ครั้งถัดไป ผู้ใช้จะได้ข้อมูลที่สดใหม่อยู่เสมอ

---

## 🧪 6.3 ผลการทดสอบประสิทธิภาพจริง (Benchmark Results)

จากผลการรันจริงในระบบ:

| ลำดับ Request | Endpoint | สถานะ Cache | Latency ที่วัดได้จริง | อัตราความเร็ว |
| :--- | :--- | :--- | :--- | :--- |
| **Request 1** | `GET /tasks` | 🐢 **Cache Miss** | **`4.8288 ms`** | Query ผ่าน PostgreSQL |
| **Request 2** | `GET /tasks` | ⚡ **Cache Hit** | **`1.2329 ms`** | อ่านจาก Redis RAM |
| **Request 3** | `POST /tasks` | 🧹 **Invalidation** | **`5.1154 ms`** | บันทึก DB + ลบ Cache เก่า |
| **Request 4** | `GET /tasks` | 🐢 **Cache Miss** | **`2.1619 ms`** | Query DB ข้อมูลใหม่ + Set Cache |
| **Request 5** | `GET /tasks` | ⚡ **Cache Hit** | **`1.1093 ms`** | อ่านจาก Redis RAM |
| **Request 6** | `GET /tasks` | ⚡ **Cache Hit** | **`674.2 µs`** | **0.67 มิลลิวินาที** 🚀 |
| **Request 7** | `GET /tasks` | ⚡ **Cache Hit** | **`502.7 µs`** | **0.50 มิลลิวินาที** 🏎️💨 |

> 🚀 **สรุปผลลัพธ์**: ความเร็วในการดึงข้อมูลเพิ่มขึ้นเกือบ **10 เท่า** (จาก 4.8ms ลดลงเหลือเพียง **502 ไมโครวินาที**) โดยใช้ทรัพยากร CPU และ Database I/O ลดลงอย่างมหาศาล!

---

## 🧠 สรุป Concept สำคัญของ Section นี้

| หัวข้อ | Concept สำคัญในภาษา Go |
| :--- | :--- |
| **Cache-Aside Pattern** | การตรวจสอบ Cache ก่อนดึง Database และการ Set ข้อมูลพร้อม TTL |
| **Cache Invalidation** | การทำ Cache Eviction เมื่อเกิด Write Operations ป้องกันข้อมูลค้าง |
| **Decorator Design Pattern** | การเขียน Struct ครอบ Interface เดิมเพื่อเพิ่มฟังก์ชัน Caching โดยไม่กระทบ Handler |
| **Redis Go Driver (`go-redis/v9`)** | การจัดการ Connection Pool, `Get`, `Set`, `Del`, และ `Scan` Iterator |
