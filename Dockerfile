# ==========================================
# Stage 1: Build Phase
# ==========================================
FROM golang:1.24-alpine AS builder

# ตั้งค่า Working Directory ภายใน Container
WORKDIR /app

# คัดลอก go.mod เพื่อดาวน์โหลด dependencies ล่วงหน้า (ถ้ามี)
COPY go.mod ./
RUN go mod download

# คัดลอก Source Code ทั้งหมด
COPY . .

# คอมไพล์ Go Code เป็น Static Binary สำหรับ Linux
# CGO_ENABLED=0 : ปิดการใช้งาน C bindings เพื่อให้รันบน Alpine ได้แบบ 100% standalone
# -ldflags="-s -w" : ตัด debug information ออกเพื่อลดขนาดไฟล์ binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o task-app .

# ==========================================
# Stage 2: Minimal Runtime Phase
# ==========================================
FROM alpine:latest

# ติดตั้ง tzdata สำหรับจัดการ Timezone และ ca-certificates สำหรับ HTTPS
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# คัดลอกเฉพาะ Binary ที่คอมไพล์แล้วมาจาก Stage 1
COPY --from=builder /app/task-app .

# กำหนด Port
ENV PORT=8080
EXPOSE 8080

# คำสั่งรันเซิร์ฟเวอร์
CMD ["./task-app"]
