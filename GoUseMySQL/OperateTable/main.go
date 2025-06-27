package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dsn := "devuser:Dev200210_@tcp(localhost:3306)/"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("连接失败：", err)
	}
	defer db.Close()

	// 选择使用 testdb 数据库
	_, err = db.Exec("use testdb")
	if err != nil {
		log.Fatal("切换数据库失败：", err)
	}

	// 表操作
	createUserTable(db)

	// 列操作演示
	addColumn(db)
	modifyColumn(db)
	renameColumn(db)
	dropColumn(db)

	// 索引操作演示
	addIndex(db)
	addUniqueIndex(db)
	addCompositeIndex(db)
	dropIndex(db)

	// 6. 清理（可选）
	dropUserTable(db)
}

// 创建基本表结构
func createUserTable(db *sql.DB) {
	sql := `
		CREATE TABLE IF NOT EXISTS users (
			user_id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
			username VARCHAR(50) NOT NULL COMMENT '用户名',
			email VARCHAR(100) NOT NULL COMMENT '邮箱地址',
			create_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (user_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
	`

	_, err := db.Exec(sql)
	if err != nil {
		log.Fatal("创建表失败：", err)
	}
	fmt.Println("user表创建成功！")
}

// 删除表
func dropUserTable(db *sql.DB) {
	_, err := db.Exec("DROP TABLE IF EXISTS users")
	if err != nil {
		log.Fatal("删除表失败：", err)
	}
	fmt.Println("user表删除成功！")
}

// 添加新列
func addColumn(db *sql.DB) {
	_, err := db.Exec("ALTER TABLE users ADD COLUMN last_login DATETIME DEFAULT NULL AFTER email")
	if err != nil {
		log.Fatal("添加列失败: ", err)
	}
	fmt.Println("成功添加last_login列")
}

// 修改列类型
func modifyColumn(db *sql.DB) {
	_, err := db.Exec("ALTER TABLE users MODIFY COLUMN username VARCHAR(70) NOT NULL")
	if err != nil {
		log.Fatal("修改列失败: ", err)
	}
	fmt.Println("成功修改username列类型")
}

// 重命名列
func renameColumn(db *sql.DB) {
	_, err := db.Exec("ALTER TABLE users RENAME COLUMN last_login TO last_login_rename_test")
	if err != nil {
		log.Fatal("重命名列失败: ", err)
	}
	fmt.Println("成功重命名email列")
}

// 删除列
func dropColumn(db *sql.DB) {
	_, err := db.Exec("ALTER TABLE users DROP COLUMN last_login_rename_test")
	if err != nil {
		log.Fatal("删除列失败: ", err)
	}
	fmt.Println("成功删除last_login列")
}

// 添加普通索引
func addIndex(db *sql.DB) {
	_, err := db.Exec("CREATE INDEX idx_username ON users (username)")
	if err != nil {
		log.Fatal("添加索引失败: ", err)
	}
	fmt.Println("成功添加username索引")
}

// 添加唯一索引
func addUniqueIndex(db *sql.DB) {
	_, err := db.Exec("CREATE UNIQUE INDEX uidx_email ON users (email)")
	if err != nil {
		log.Fatal("添加唯一索引失败: ", err)
	}
	fmt.Println("成功添加email唯一索引")
}

// 添加复合索引
func addCompositeIndex(db *sql.DB) {
	_, err := db.Exec("CREATE INDEX idx_user_status ON users (username, email)")
	if err != nil {
		log.Fatal("添加复合索引失败: ", err)
	}
	fmt.Println("成功添加复合索引")
}

// 删除索引
func dropIndex(db *sql.DB) {
	_, err := db.Exec("DROP INDEX idx_username ON users")
	if err != nil {
		log.Fatal("删除索引失败: ", err)
	}
	fmt.Println("成功删除username索引")
}
