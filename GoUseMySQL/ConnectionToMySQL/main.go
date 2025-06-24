package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// DSN 格式：`[username]:[password]@tcp([host]:[port])/[database]`
	dsn := "devuser:Dev200210_@tcp(127.0.0.1:3306)/testdb"
	// dsn := "devuser:Dev200210_@tcp(localhost:3306)/"	// 不指定具体数据库

	// 建立数据库连接（此时已未制定具体数据库）
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("连接失败：", err)
	}
	defer db.Close()

	// 验证连接
	if err := db.Ping(); err != nil {
		log.Fatal("连接验证失败：", err)
	}
	fmt.Println("成功连接到MySQL服务器!")
}
