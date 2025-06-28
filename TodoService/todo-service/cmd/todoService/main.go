package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todo-service/database"
	"todo-service/handlers"
	"todo-service/middleware"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// 初始化DB
	database.Init()
	defer database.Close()

	router := gin.Default()

	// 添加Prometheus监控中间件
	router.Use(middleware.MetricsMiddleware())

	// 注册路由
	api := router.Group("/api")
	{
		api.POST("/tasks", handlers.CreateTask)
		api.PATCH("/tasks/:id/start", handlers.StartTask)
		api.PATCH("/tasks/:id/complete", handlers.CompleteTask)
		api.DELETE("/tasks/:id", handlers.DeleteTask)
		api.GET("/tasks", handlers.ListTasks)
	}

	// 添加Prometheus指标端点
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// 启动连接池监控
	middleware.StartDBPoolMonitor(database.DB)

	// 配置服务器
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// 启动服务器
	go func() {
		log.Println("服务启动在 http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务启动失败: %v\n", err)
		}
	}()

	// 设置优雅关机
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务器...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("服务关闭失败: %v\n", err)
	}

	log.Println("服务器已关闭")
}
