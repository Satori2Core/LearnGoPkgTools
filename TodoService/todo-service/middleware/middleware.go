package middleware

import (
	"database/sql"
	"time"
	"todo-service/metrics"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP请求耗时",
			Buckets: []float64{0.001, 0.01, 0.1, 0.5, 1, 5, 10},
		},
		[]string{"method", "path", "status"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestDuration)
}

// 添加监控中间件
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		duration := time.Since(start).Seconds()
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.FullPath()

		httpRequestDuration.WithLabelValues(method, path, string(status)).Observe(duration)
	}
}

// 定期收集连接池状态
func StartDBPoolMonitor(db *sql.DB) {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			metrics.RecordDBPoolStats(db)
		}
	}()
}
