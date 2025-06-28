package metrics

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// 数据库操作计数器
	DBOpCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_operations_total",
			Help: "总数据库操作次数",
		},
		[]string{"operation", "status"},
		// operation: create_task, start_task, etc.; status: success, error, invalid
	)

	// 数据库操作耗时
	DBOpDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_operation_duration_seconds",
			Help:    "数据库操作耗时",
			Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 2, 5, 10},
		},
		[]string{"operation"},
	)

	// 连接池监控
	DBPoolOpen = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "db_pool_open_connections",
		Help: "当前打开的数据库连接数",
	})

	DBPoolInUse = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "db_pool_in_use_connections",
		Help: "当前正在使用的数据库连接数",
	})

	DBPoolIdle = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "db_pool_idle_connections",
		Help: "当前空闲的数据库连接数",
	})
)

// 记录连接池状态
func RecordDBPoolStats(db *sql.DB) {
	stats := db.Stats()
	DBPoolOpen.Set(float64(stats.OpenConnections))
	DBPoolInUse.Set(float64(stats.InUse))
	DBPoolIdle.Set(float64(stats.Idle))
}
