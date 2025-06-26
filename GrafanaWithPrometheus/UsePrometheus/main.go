package main

import (
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// 定义指标
var (
	// 示例一：随机温度计（仪表盘）
	tempGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "room_temperature_celsius",
			Help: "当前室温摄氏度",
		},
	)

	// 示例二：API调用计数器
	apiCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "api_requests_total",
			Help: "API总调用次数",
		},
	)
)

func main() {
	// 示例一：随机温度计（仪表盘）
	ExampleOne()

	// 示例二：API调用计数器
	ExampleTwo()

	// 露监控端点
	http.Handle("/metrics", promhttp.Handler())

	println("服务启动：:8080")
	http.ListenAndServe(":8080", nil)
}

// 示例一：随机温度计（仪表盘）
func ExampleOne() {
	// 注册到默认注册表
	prometheus.MustRegister(tempGauge)

	// 后台定时更新指标值
	// 模拟温度度变化
	go func() {
		for {
			// 生成20-30之间的随机温度
			temp := 20 + rand.Float64()*10
			tempGauge.Set(temp)
			time.Sleep(5 * time.Second)
		}
	}()
}

// 示例二：API调用计数器
func ExampleTwo() {
	// 注册到默认注册表
	prometheus.MustRegister(apiCounter)

	// 业务接口
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		// 每次调用API时计数器+1
		apiCounter.Inc()
		w.Write([]byte("API请求成功"))
	})
}
