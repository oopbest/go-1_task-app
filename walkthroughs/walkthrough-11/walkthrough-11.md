# Walkthrough 11: Enterprise Full-Stack SaaS (Next.js 15, React Query v5 & Go Backend Integration)

ในส่วนนี้เป็นการพัฒนา **Frontend Web Application ระดับ Enterprise SaaS** ด้วย **Next.js 15 (App Router)**, **TypeScript**, **Tailwind CSS**, และ **TanStack Query v5** เพื่อเชื่อมต่อกับ **Go Backend API Gateway**, **Redis Cache (500µs)**, **PostgreSQL**, และ **gRPC Microservice** อย่างไร้รอยต่อ 🚀🎨💻✨

---

## 🏗️ โครงสร้างไฟล์ Frontend ที่พัฒนาในส่วนนี้

```
1_task-app/
├── middleware/
│   └── cors.go                        # [NEW] Gin CORS Middleware อนุญาต Request จาก Next.js
├── frontend/                          # [NEW] Next.js 15 App Router Project
│   ├── src/
│   │   ├── app/
│   │   │   ├── layout.tsx             # Root Layout สวม Providers (Inter Font + Dark Theme)
│   │   │   ├── page.tsx               # Root Home Page (Auto Redirect Guard)
│   │   │   ├── login/page.tsx         # Modern Glassmorphic Login Page
│   │   │   ├── register/page.tsx      # User Registration Page
│   │   │   └── dashboard/page.tsx     # SaaS Task Management Dashboard
│   │   ├── context/
│   │   │   └── auth-context.tsx       # Auth Context & JWT Session Guard
│   │   ├── providers/
│   │   │   └── query-provider.tsx     # TanStack React Query & Sonner Toaster Provider
│   │   ├── lib/
│   │   │   ├── api.ts                 # Type-Safe API Client (Bearer JWT Fetcher)
│   │   │   └── utils.ts               # Tailwind Class Merger (clsx + tailwind-merge)
│   │   └── types/
│   │       └── index.ts               # TypeScript Interfaces matching Go Models
│   └── package.json
└── README.md                          # [UPDATE] เอกสารประกอบโปรเจกต์ฉบับสมบูรณ์
```

---

## 🔍 11.1 สถาปัตยกรรม Full-Stack System Pipeline

```mermaid
flowchart TD
    subgraph FrontendTier["💻 Frontend Tier (Next.js 15 on :3001)"]
        Browser["🌐 User Browser (Next.js App)"]
        AuthCtx["🔐 AuthContext (JWT Cookies / Storage)"]
        Query["⚡ TanStack Query v5<br>(Optimistic Updates & Client Cache)"]
        UI["🎨 Tailwind CSS + Glassmorphism UI"]
    end

    subgraph BackendTier["🚀 Backend Tier (Go Gin Gateway on :8080)"]
        CORS["🛡️ CORSMiddleware()"]
        GinRouter["Radix Tree Router"]
        WorkerPool["⚙️ Worker Pool (Goroutines)"]
    end

    subgraph StorageTier["💾 Data & Microservices Tier"]
        Redis["⚡ Redis Cache (~500µs)"]
        Postgres["🐘 PostgreSQL Database (:5432)"]
        gRPC["📬 Notification Microservice (:50051)"]
    end

    Browser --> UI
    UI --> Query
    Query -->|HTTP / JSON with Bearer Token| CORS
    CORS --> GinRouter
    GinRouter <--> Redis
    GinRouter <--> Postgres
    GinRouter --> WorkerPool
    WorkerPool -->|gRPC Protobuf (HTTP/2)| gRPC
```

---

## ⚡ 11.2 จุดเด่นของระบบ Full-Stack SaaS

### 1. Optimistic UI Updates (TanStack Query v5)
- เมื่อผู้ใช้กดคลิกสลับสถานะงาน (เสร็จ / รอดำเนินการ) หน้าจอจะอัปเดตสถานะทันทีใน **0 มิลลิวินาที** โดยไม่ต้องรอการตอบกลับจาก Server
- หากเกิด Network Error ระบบจะ Rollback กลับสถานะเดิมอัตโนมัติพร้อมแจ้งเตือนผ่าน Toast

### 2. Live Debounced Search & Dynamic Filters
- มีระบบ Debounce 300ms เมื่อพิมพ์ค้นหางาน เพื่อลดภาระการยิง Request ถี่เกินไป
- เชื่อมต่อกับ Backend SQL `ILIKE` และระบบ Sorting สลับ ลำดับใหม่ ➡️ เก่า ได้ทันที

### 3. Modern Glassmorphic Dark UI
- ดีไซน์สวยหรูระดับ State-of-the-Art ด้วย Glowing Background Orbs, Translucent Cards (`backdrop-blur-xl`), และ Micro-Animations
- ระบบแจ้งเตือนด้วย **Sonner Toast Notification** ที่มุมขวาบน

---

## 🧪 11.3 ผลการทดสอบจริง (Live Dashboard Screenshot)

จากผลการทดสอบผ่านหน้าเว็บ `http://localhost:3001/dashboard`:
- **ผู้ใช้งานปัจจุบัน**: `test@example.com` (ID #4)
- **Status Cards**: แสดงผลการเชื่อมต่อกับ `Gin REST :8080`, `Redis ~500µs`, `gRPC :50051`, และ `Total Tasks: 2`
- **Task Cards**: แสดงงาน `#Task 18: Learn Go` และ `#Task 19: Learn Java` พร้อมฟังก์ชันสลับสถานะ, แก้ไข, และลบงานแบบ Real-Time! 🚀

---

## 🧠 สรุป Concept สำคัญของ Section นี้

| หัวข้อ | Concept สำคัญ (Full-Stack Engineering) |
| :--- | :--- |
| **CORS Policy** | การเปิด Header `Access-Control-Allow-Origin` เพื่อให้ Browser คุยกับ API ต่าง Port |
| **Next.js 15 App Router** | การจัดโครงสร้างแบบ Nested Layouts, Client Components (`'use client'`), และ Navigation Guard |
| **TanStack Query v5** | การใช้งาน `useQuery`, `useMutation`, `queryClient.setQueryData`, และ `cancelQueries` สำหรับ Optimistic UI |
| **JWT Session Guard** | การจัดเก็บ Token ใน Cookie + LocalStorage และการเขียน Interceptor แนบ Header อัตโนมัติ |
