package metrics

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// 1. กำหนด Prometheus Metrics
var (
	// HttpRequestsTotal นับจำนวน Request ทั้งหมด แยกตาม method, path, และ status code
	HttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests processed by the Task API",
		},
		[]string{"method", "path", "status"},
	)

	// HttpRequestDuration วัดการกระจายตัวของระยะเวลาประมวลผล (Latency) ในหน่วยวินาที
	HttpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of response latency for HTTP requests",
			Buckets: []float64{0.0005, 0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5}, // 500µs ถึง 2.5s
		},
		[]string{"method", "path"},
	)

	// WorkerJobsTotal นับจำนวนงานที่ Background Worker ทำเสร็จ แยกตาม type และ status
	WorkerJobsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "task_worker_jobs_total",
			Help: "Total number of background worker jobs executed",
		},
		[]string{"type", "status"},
	)
)

// PrometheusMiddleware ดักจับทุก Request และบันทึก Metrics ส่งให้ Prometheus อัตโนมัติ
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		method := c.Request.Method

		// บันทึกตัวเลขสถิติ
		HttpRequestsTotal.WithLabelValues(method, path, status).Inc()
		HttpRequestDuration.WithLabelValues(method, path).Observe(duration)
	}
}
