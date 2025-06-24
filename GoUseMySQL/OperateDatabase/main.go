package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dsn := "devuser:Dev200210_@tcp(127.0.0.1:3306)/"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("连接失败：", err)
	}
	defer db.Close()

	// 创建数据库
	createDatabase(db)

	// 删除数据库
	deleteDatabase(db)
}

// 创建数据库
func createDatabase(db *sql.DB) {
	_, err := db.Exec("CREATE DATABASE IF NOT EXISTS go_db DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci")
	if err != nil {
		log.Fatal("，数据库创建失败：", err)
	}
	fmt.Println("数据库创建成功！")

	// 使用新建的数据库
	_, err = db.Exec("use go_db")
	if err != nil {
		log.Fatal("切换数据库失败: ", err)
	}
}

// 删除数据库
func deleteDatabase(db *sql.DB) {
	_, err := db.Exec("DROP DATABASE IF EXISTS go_db")
	if err != nil {
		log.Fatal("删除数据库失败: ", err)
	}
	fmt.Println("数据库删除成功!")
}
