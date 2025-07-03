# Go & SQLx库 —— 数据库表操作使用指南

## 1. 安装与初始化

### 1.1 安装依赖

```bash
# sqlx 驱动
go get github.com/jmoiron/sqlx

# go mysql 驱动
go get github.com/go-sql-driver/mysql
```

---

### 1.2 初始化数据库连接

- **关键方法**：`sqlx.Open("mysql", dsn)`

```go
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
```

> 执行测试

```bash
# 示例：开始 -----------------------------------------------------------------------------
➜  ConnectionToMysql git:(learn/Go-Use-sqlx) ✗ go run main.go 
数据库连接成功！
➜  ConnectionToMysql git:(learn/Go-Use-sqlx) ✗
# 示例：结束 -----------------------------------------------------------------------------
```

---

## 2. 数据库/表操作

### 2.1 创建/删除 —— 数据库

```go
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
```

> 可能出现的问题：如下错误表明，当前用户没有数据库操作权限

```bash
# 示例：开始 -----------------------------------------------------------------------------
➜  OperateDatabase git:(learn/Go-Use-sqlx) ✗ go run main.go
2025/07/02 22:10:17 创建数据库失败: Error 1044 (42000): Access denied for user 'devuser'@'localhost' to database 'sqlx_learning'
exit status 1
➜  OperateDatabase git:(learn/Go-Use-sqlx) ✗ 
# 示例：结束 -----------------------------------------------------------------------------
```

> 解决办法：使用root用户登陆，对devuser用户授予权限

```sql
-- 此处直接授权对指定数据库的全部操作权限
GRANT CREATE, DROP, ALTER, INDEX, INSERT, SELECT, UPDATE, DELETE, 
    EXECUTE, CREATE VIEW, SHOW VIEW, CREATE TEMPORARY TABLES, 
    LOCK TABLES, REFERENCES, TRIGGER
ON 数据库名.* TO 'devuser'@'localhost';
```

> 执行测试

```bash
# 示例：开始 -----------------------------------------------------------------------------
➜  OperateDatabase git:(learn/Go-Use-sqlx) ✗ go run main.go
sqlx_learning 数据库创建成功！
sqlx_learning 数据库删除成功！
➜  OperateDatabase git:(learn/Go-Use-sqlx) ✗ 
# 示例：结束 -----------------------------------------------------------------------------
```

---

### 2.2 创建数据表

- **关键方法**：`Exec(sql)`

#### 2.2.1 模型定义

```go
type User struct {
    ID        int64     `db:"id" json:"id"`
    Username  string    `db:"username" json:"username"`
    Email     string    `db:"email" json:"email"`
    Password  string    `db:"password" json:"-"`
    CreatedAt time.Time `db:"created_at" json:"created_at"`
    UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type Task struct {
    ID          int64     `db:"id" json:"id"`
    UserID      int64     `db:"user_id" json:"user_id"`
    Title       string    `db:"title" json:"title"`
    Description string    `db:"description" json:"description"`
    Status      string    `db:"status" json:"status"` // 'pending', 'in_progress', 'completed'
    DueDate     time.Time `db:"due_date" json:"due_date"`
    CreatedAt   time.Time `db:"created_at" json:"created_at"`
    UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
```

---

#### 2.2.2 创建数据表

```go
func createUsersTable() {
    query := `
        CREATE TABLE IF NOT EXISTS users (
            id INT AUTO_INCREMENT PRIMARY KEY,
            username VARCHAR(50) NOT NULL UNIQUE,
            email VARCHAR(100) NOT NULL UNIQUE,
            password VARCHAR(255) NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
    `
    
    _, err := db.Exec(query)
    if err != nil {
        log.Fatalf("创建 users 表失败: %v", err)
    }
    
    fmt.Println("users 表创建成功！")
}

func createTasksTable() {
    query := `
        CREATE TABLE IF NOT EXISTS tasks (
            id INT AUTO_INCREMENT PRIMARY KEY,
            user_id INT NOT NULL,
            title VARCHAR(255) NOT NULL,
            description TEXT,
            status ENUM('pending', 'in_progress', 'completed') DEFAULT 'pending',
            due_date DATETIME,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
            FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
    `
    
    _, err := db.Exec(query)
    if err != nil {
        log.Fatalf("创建 tasks 表失败: %v", err)
    }
    
    fmt.Println("tasks 表创建成功！")
}
```

---


## 3. 基本 CRUD 操作

> 基于前文创建的数据表与模型