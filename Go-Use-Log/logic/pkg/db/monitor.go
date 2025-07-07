package db

import (
	"context"
	"log/slog"
	"time"

	"gorm.io/gorm"
)

// 慢查询阈值
var SlogQueryThreshold = 200 * time.Millisecond

// WrapDB 包装原生 *gorm.DB 实现慢查询监控
type WrapDB struct {
	*gorm.DB
	ctx context.Context
}

func WithContext(db *gorm.DB, ctx context.Context) *WrapDB {
	return &WrapDB{
		DB:  db,
		ctx: ctx,
	}
}

func (db *WrapDB) callbackRegister() {
	// 注册查询前后的回调
	db.Callback().Query().Before("gorm:query").Register("monitor:before_query", beforeQuery)
	db.Callback().Query().After("gorm:query").Register("monitor:after_query", afterQuery)
	// 同理注册Update/Delete/Create等...
}

func beforeQuery(db *gorm.DB) {
	db.Set("start_time", time.Now())
}

func afterQuery(db *gorm.DB) {
	start, ok := db.Get("start_time")
	if !ok {
		return
	}

	duration := time.Since(start.(time.Time))
	if duration > SlogQueryThreshold {
		// 从上下文中获取请求 Logger
		if logger, ok := db.Statement.Context.Value("logger").(*slog.Logger); ok {
			logger.LogAttrs(db.Statement.Context, slog.LevelWarn, "SLOW_QUERY",
				slog.String("table", db.Statement.Table),
				slog.String("operation", db.Statement.SQL.String()),
				slog.Int64("duration_ms", duration.Milliseconds()),
				slog.Int64("rows_affected", db.Statement.RowsAffected),
			)
		}
	}
}
