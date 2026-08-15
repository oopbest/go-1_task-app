package middleware

import (
	"log"
	"net/http"
	"time"
)

// responseWriterWrapper ช่วยดักจับ HTTP Status Code ที่ Handler ส่งออกไป เพื่อนำมาแสดงใน Log
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader ดักจับ status code ก่อนส่งต่อให้ ResponseWriter ตัวจริง
func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Logging ดักจับทุก Request เพื่อบันทึก Method, Path, Status Code และ Latency
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapper := &responseWriterWrapper{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // default status คือ 200
		}

		// ส่ง Request ต่อไปให้ Handler หรือ Middleware ตัวถัดไปทำงาน
		next.ServeHTTP(wrapper, r)

		duration := time.Since(start)
		log.Printf("[%s] %s | Status: %d | Latency: %v\n",
			r.Method,
			r.URL.Path,
			wrapper.statusCode,
			duration,
		)
	})
}

// JSONContentType กำหนด Header Content-Type เป็น application/json ให้อัตโนมัติทุก Request
func JSONContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}
