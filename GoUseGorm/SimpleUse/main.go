package main

import (
	"log/slog"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "devuser:Dev200210_@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("数据库连接失败：" + err.Error())
	}

	// 获取底层数据库连接池
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	slog.Info("数据库连接成功")

	sqlDB.SetMaxOpenConns(100) // 最大连接数
	sqlDB.SetMaxIdleConns(10)  // 最大空闲连接
}

type User struct {
	ID        int64     `gorm:"primaryKey"`
	Name      string    `gorm:"type:varchar(100);not null"`
	Email     string    `gorm:"type:varchar(100);uniqueIndex"`
	CreatedAt time.Time // 自动记录创建时间（GORM约定）
	UpdatedAt time.Time // 自动记录更新时间（GORM约定）
}

// 封装公用字段
type BaseModel struct {
	ID        int64 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Product struct {
	BaseModel // 嵌入
	Name      string
	Price     float64
}
