package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 日志级别阈值设置
const (
	SlowRequestThreshold = 500 * time.Millisecond
	ErrorThreshold       = 5 // 内部错误超过五个，则输出完整日志
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 创建请求专用日志 Logger
		start := time.Now()
		path := c.Request.URL.Path

		// 2. 创建子 logger 添加请求元数据
		reqLogger := slog.Default().With(
			slog.Group("request",
				slog.String("id", c.GetHeader("X-Request-ID")),
				slog.String("method", c.Request.Method),
				slog.String("path", path),
				slog.String("ip", c.ClientIP()),
				slog.String("user_agent", c.Request.UserAgent()),
			),
		)

		// 3. 将 logger 存入 context
		ctx := context.WithValue(c.Request.Context(), "logger", reqLogger)
		c.Request = c.Request.WithContext(ctx)

		// 4. 请求错误收集器（轻量级记录）
		errCounter := 0

		// 5. 业务处理
		c.Next()

		// 6. 请求处理完成记录日志
		latency := time.Since(start)
		status := c.Writer.Status()

		logAttrs := []slog.Attr{
			slog.Int("status", status),
			slog.Int64("latency_ms", latency.Milliseconds()),
			slog.Int("response_size", c.Writer.Size()),
		}

		// 7. 动态日志级别决策
		switch {
		case latency > SlowRequestThreshold:
			reqLogger.LogAttrs(ctx, slog.LevelWarn, "SLOW_REQUEST", logAttrs...)
		case latency >= http.StatusInternalServerError:
			reqLogger.LogAttrs(ctx, slog.LevelError, "SERVER_ERROR", logAttrs...)
		case latency > ErrorThreshold:
			reqLogger.LogAttrs(ctx, slog.LevelError, "BUSINESS_ERRORS",
				append(logAttrs, slog.Int("error_count", errCounter))...)
		default:
			reqLogger.LogAttrs(ctx, slog.LevelInfo, "REQUEST_COMPLELTED", logAttrs...)
		}
	}
}
