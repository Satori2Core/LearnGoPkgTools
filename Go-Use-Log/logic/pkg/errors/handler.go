package errors

import (
	"context"
	"log/slog"
	"runtime"

	"github.com/pkg/errors"
)

func LogError(ctx context.Context, err error) {
	if err == nil {
		return
	}

	// 获取调用栈（限制深度）
	pcs := make([]uintptr, 32)
	n := runtime.Callers(3, pcs)
	frames := runtime.CallersFrames(pcs[:n])

	// 构建错误元数据
	stack := []slog.Attr{}
	for {
		frame, more := frames.Next()
		stack = append(stack, slog.String("frame", frame.Function))
		if !more {
			break
		}
	}

	// 获取上下文中的logger
	logger := ctx.Value("logger").(*slog.Logger)

	logger.LogAttrs(ctx, slog.LevelError, "INTERNAL_ERROR",
		slog.String("message", errors.Cause(err).Error()),
		slog.Int("stack_depth", len(stack)),
		slog.Group("stack", stack...),
	)
}
