package main

import (
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

const (
	DBName = "sqlx_learning"
)

type User struct {
	ID        int64     `db:"id" json:"id"`
	Username  string    `db:"username" json:"username"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func main() {

	createUsersTable()

	dropDatabase()
	if db != nil {
		db.Close()
	}
}

func init() {
	creataDatabase()
}

func createUsersTable() {
	fmt.Println("test")

	if db == nil {
		fmt.Println("nil")
		return
	}

	_, err := db.Exec("use sqlx_learning")
	if err != nil {
		log.Fatalf("切换数据库失败: %v", err)
	}

	sql := `
        CREATE TABLE IF NOT EXISTS users (
            id INT AUTO_INCREMENT PRIMARY KEY,
            username VARCHAR(50) NOT NULL UNIQUE,
            email VARCHAR(100) NOT NULL UNIQUE,
            password VARCHAR(255) NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
    `

	_, err = db.Exec(sql)
	if err != nil {
		log.Fatalf("创建users表失败: %v", err)
	}
	fmt.Println("users 表创建成功！")
}

func creataDatabase() {
	dsn := "devuser:Dev200210_@tcp(localhost:3306)/"
	var err error
	db, err = sqlx.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// 创建数据库
	sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s DEFAULT CHARSET utf8mb4", DBName)
	_, err = db.Exec(sql)
	if err != nil {
		log.Fatalf("创建数据库失败: %v", err)
	}

	fmt.Printf("%s 数据库创建成功！\n", DBName)

	// 配置连接池
	db.SetMaxOpenConns(20)                 // 最大打开连接数
	db.SetMaxIdleConns(5)                  // 最大空闲连接数
	db.SetConnMaxLifetime(5 * time.Minute) // 连接最长存活时间

	fmt.Printf("%s 数据库连接成功！\n", DBName)
}

func dropDatabase() {
	if db == nil {
		return
	}

	// 删除数据库
	sql := fmt.Sprintf("DROP DATABASE IF EXISTS %s ", DBName)
	_, err := db.Exec(sql)
	if err != nil {
		log.Fatalf("删除数据库失败: %v", err)
	}

	fmt.Printf("%s 数据库删除成功！\n", DBName)
}
