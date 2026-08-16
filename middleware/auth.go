package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/oopbest/task-app/utils"
)

type contextKey string

const UserIDContextKey contextKey = "user_id"

// GetUserIDFromContext ฟังก์ชัน Helper สำหรับดึง UserID ออกมาจาก Request Context ใน Handler
func GetUserIDFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(UserIDContextKey).(int)
	return userID, ok
}

// AuthMiddleware ตรวจสอบ JWT Bearer Token ก่อนอนุญาตให้เข้าถึง API
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error":"Authorization header is required"}`, http.StatusUnauthorized)
			return
		}

		// รูปแบบต้องเป็น "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, `{"error":"Invalid authorization format. Format must be Bearer <token>"}`, http.StatusUnauthorized)
			return
		}

		tokenStr := parts[1]

		// ตรวจสอบความถูกต้องของ JWT Token
		claims, err := utils.ValidateJWT(tokenStr)
		if err != nil {
			http.Error(w, `{"error":"Invalid or expired token"}`, http.StatusUnauthorized)
			return
		}

		// ฝาก UserID ไว้ใน Request Context เพื่อให้ Handler หยิบไปใช้
		ctx := context.WithValue(r.Context(), UserIDContextKey, claims.UserID)

		// ส่งต่อไปยัง Handler ตัวถัดไปพร้อม Context ใหม่
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
