package main

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	DbName = "sqlx_learning"
)

func main() {
	// 创建/删除 —— 数据库
	createAndDropDatabaseTest()
}

// 创建/删除 —— 数据库
func createAndDropDatabaseTest() {
	// 首先连接不指定数据库的实例
	dsn := "devuser:Dev200210_@tcp(localhost:3306)/"
	DB, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer DB.Close()

	// 创建数据库
	sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s DEFAULT CHARSET utf8mb4", DbName)
	_, err = DB.Exec(sql)
	if err != nil {
		log.Fatalf("创建数据库失败: %v", err)
	}

	fmt.Printf("%s 数据库创建成功！\n", DbName)

	// 删除数据库
	sql = fmt.Sprintf("DROP DATABASE IF EXISTS %s ", DbName)
	_, err = DB.Exec(sql)
	if err != nil {
		log.Fatalf("删除数据库失败: %v", err)
	}

	fmt.Printf("%s 数据库删除成功！\n", DbName)
}
