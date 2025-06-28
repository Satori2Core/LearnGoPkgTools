package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Init() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "devuser:Dev200210_@tcp(127.0.0.1:3306)/"
	}

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("数据库连接失败：", err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatal("数据库无法连通：", err)
	}

	// 设置连接池
	DB.SetMaxOpenConns(50)
	DB.SetMaxIdleConns(20)
	DB.SetConnMaxLifetime(10 * time.Minute)

	log.Println("成功连接MySQL数据库")

	createTodoDatabase()

	// 使用新建的数据库
	_, err = DB.Exec("use testdb")
	if err != nil {
		log.Fatal("切换数据库失败: ", err)
	}
	fmt.Println("切换数据库成功！")

	createTodoTable()

}

func Close() {
	if DB != nil {
		DB.Close()
		log.Println("关闭数据库连接")
	}
}

// createTodoDatabase 创建数据库
func createTodoDatabase() {
	// 创建数据库
	_, err := DB.Exec("CREATE DATABASE IF NOT EXISTS testdb DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;")
	if err != nil {
		log.Fatal("数据库创建失败：", err)
	}
	fmt.Println("数据库创建成功！")
}

// createTodoTable 创建数据表
func createTodoTable() {
	_, err := DB.Exec(
		`
		CREATE TABLE IF NOT EXISTS todo_tasks (
			id INT AUTO_INCREMENT PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			creator VARCHAR(100) NOT NULL,
			assignee VARCHAR(100) NOT NULL,
			status ENUM('pending', 'in_progress', 'completed', 'deleted') NOT NULL DEFAULT 'pending',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			started_at DATETIME,
			completed_at DATETIME,
			deleted_at DATETIME,
			INDEX idx_creator (creator),
			INDEX idx_assignee (assignee),
			INDEX idx_status (status),
			INDEX idx_completed_at (completed_at)
		);
		`,
	)
	if err != nil {
		log.Fatal("创建表失败：", err)
	}
	fmt.Println("todo_tasks表创建成功！")
}
