package service

import (
	"context"
	"log/slog"

	"github.com/Satori2Core/LearnGoPkgTools/Go-Use-Log/logic/pkg/db"
)

type OrderService struct {
	db *db.WrapDB
}

type Item struct {
	ID        string
	ProductID int
	Quantity  int
}

func (s *OrderService) CreateOrder(ctx context.Context, userID int64, items []Item) error {
	// 从上下文中获取请求 logger
	logger := ctx.Value("logger").(*slog.Logger)

	logger.Info("开始创建订单", slog.Int64("user_id", userID), slog.Int("item_count", len(items)))

	tx := s.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			logger.Error("订单创建失败：事务回滚", slog.Any("reason", r))
		}
	}()

	// 业务处理（略）

	// 关键业务单独记录
	logger.Debug("库存扣除操作", slog.Int("items", logItems(items)))

	if err := tx.Commit().Error; err != nil {
		logger.Error("订单提交失败", "error", err)
	}

	logger.Info("订单创建成功", slog.String("order_id", createdOrder.ID))
	return nil
}

// 安全记录敏感数据示例
func logItems(items []Item) slog.Value {
	safeItems := make([]slog.Value, 0, len(items))
	for _, item := range items {
		safeItems = append(safeItems, slog.GroupValue(
			slog.Int("product_id", item.ProductID),
			slog.Int("quantity", item.Quantity),
			// 不记录敏感价格信息
		))
	}
	return slog.GroupValue(safeItems...)
}
