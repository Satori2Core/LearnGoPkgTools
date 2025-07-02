package main

import (
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// 全局数据库连接
var db *sqlx.DB

var (
	// 数据库连接参数
	Username = "devuser"    // "your_username"
	Password = "Dev200210_" // "your_password"
	Host     = "localhost"
	Port     = "3306"
	DbName   = "sqlx_learning"
)

func initDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", Username, Password, Host, Port)

	var err error
	db, err = sqlx.Open("mysql", dsn)
	if err != nil {
		log.Fatal("数据库连接失败：", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("数据库测试连接失败：", err)
	}

	// 配置连接池
	db.SetMaxOpenConns(20)                 // 最大打开连接数
	db.SetMaxIdleConns(5)                  // 最大空闲连接数
	db.SetConnMaxLifetime(5 * time.Minute) // 连接最长存活时间

	fmt.Println("数据库连接成功！")
}

func main() {
	initDB()
	defer db.Close()
}
