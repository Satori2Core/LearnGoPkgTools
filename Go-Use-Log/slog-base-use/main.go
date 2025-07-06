package main

import (
	"log/slog"
	"os"
)

func main() {
	withAttrs()
}

func withAttrs() {
	// 1. 创建自定义 handler
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug, // 日志级别
	})

	// 2. 创建 Logger 实例
	logger := slog.New(handler)

	// 创建带固定属性的Logger
	userLogger := logger.With(
		slog.Int("user_id", 1001),
		slog.String("region", "us-west"),
	)

	// 后续所有日志自动携带这些属性
	userLogger.Info("更改设置", "setting", "dark_mode")
}

func customLogger() {
	// 1. 创建自定义 handler
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug, // 日志级别
	})

	// 2. 创建 Logger 实例
	logger := slog.New(handler)

	// 3. 使用自定义 Logger
	logger.Debug("调试信息") // 使用 Debug 级别
	logger.Info("业务事件", "data", "data")

}

// test01
func test01() {
	slog.Info("用户登录", "user_id", 123, "ip", "192.168.1.1")
	slog.Info("订单创建",
		slog.Int("order_id", 1001),
		slog.Float64("amount", 99.99),
	)
	slog.Info("API调用",
		slog.Group("request",
			"method", "POST",
			"path", "/users",
		),
		slog.Group("response",
			"status", 201,
			"latency_ms", 45.2,
		),
	)
}
