package middleware

import (
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	// ตั้งค่า Default Logger ของ Go ให้เป็น JSON Handler
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)
}

// StructuredLogger Middleware บันทึก Log ทุก Request ในรูปแบบ JSON
func StructuredLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		userID, _ := c.Get("user_id")

		// คัดแยก Log Level ตาม Status Code
		attrs := []any{
			slog.String("method", method),
			slog.String("path", path),
			slog.String("query", query),
			slog.Int("status", status),
			slog.String("ip", clientIP),
			slog.Duration("latency", latency),
			slog.Any("user_id", userID),
		}

		if len(c.Errors) > 0 {
			attrs = append(attrs, slog.String("error", c.Errors.String()))
		}

		if status >= 500 {
			slog.Error("HTTP Server Error", attrs...)
		} else if status >= 400 {
			slog.Warn("HTTP Client Error", attrs...)
		} else {
			slog.Info("HTTP Request Completed", attrs...)
		}
	}
}
