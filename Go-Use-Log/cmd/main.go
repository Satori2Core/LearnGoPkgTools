package main

import (
	"app/db"
	"app/middleware"
	"context"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. 初始化日志（分环境配置）
	setupLogger()

	// 2. 初始化DB（带慢查询监控）
	dsn := "host=db user=app dbname=app password=pass"
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("数据库初始化失败", "error", err)
		os.Exit(1)
	}

	// 3. 包装DB实例并注册监控
	wrappedDB := db.WithContext(gormDB, context.Background())
	wrappedDB.callbackRegister()

	// 4. 创建Gin引擎
	r := gin.New()

	// 5. 注册中间件
	r.Use(
		gin.Recovery(), // 结合上文的错误处理
		middleware.LoggingMiddleware(),
	)

	// 6. 注册路由
	r.POST("/orders", createOrderHandler)

	// 7. 启动服务
	slog.Info("服务启动成功", "port", 8080)
	r.Run(":8080")
}

func setupLogger() {
	isProduction := os.Getenv("ENV") == "production"

	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// 生产环境移除敏感路径
			if isProduction && a.Key == "path" {
				return slog.Attr{}
			}
			return a
		},
	}

	var handler slog.Handler
	if isProduction {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	slog.SetDefault(slog.New(handler))
}
