package main

import (
	"log/slog"
)

func main() {
	slog.Info("查询完成", "duration_ms", slog.Int("value", 85))

}
