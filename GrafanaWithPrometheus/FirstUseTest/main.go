package main

import (
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// 定义三个指标
var (
	randomNumber = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "random_number",
		Help: "随机生成（0～100）数值",
	})

	apiCalls = promauto.NewCounter(prometheus.CounterOpts{
		Name: "api_calls_total",
		Help: "API调用总数",
	})

	requestDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "request_duration_seconds",
		Help:    "API请求处理时间",
		Buckets: []float64{0.1, 0.3, 0.5, 1.0},
	})
)

func main() {
	// 后台生成随机数
	go func() {
		for {
			value := rand.Float64() * 100
			randomNumber.Set(value)
			time.Sleep(10 * time.Second) // 每两秒更新一次
		}
	}()

	// 业务API
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			duration := time.Since(start).Seconds()
			requestDuration.Observe(duration)
		}()

		// 模拟API处理请求
		time.Sleep(time.Duration(rand.IntN(500)) * time.Millisecond)
		apiCalls.Inc()
		w.Write([]byte("API请求成功！"))
	})

	// 暴露监控端点
	http.Handle("/metrics", promhttp.Handler())

	// 启动服务
	println("服务运行在 :2112 端口")
	http.ListenAndServe(":2112", nil)
}
