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

- **关键方法**：`Exec(query string, args ...interface{})`
- 执行写操作（INSERT/UPDATE/DELETE）

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

### 3.1 核心方法表

|类别​	|​方法签名​	|​说明​	|​使用场景​	|​返回值​|
|:-----|:-----|:-----|:-----|:-----|
|​连接管理​	|sqlx.Connect(driverName, dataSourceName string) (*DB, error)	|创建连接并立即 |ping 数据库验证	应用程序启动时建立数据库连接	|*sqlx.DB|
||sqlx.Open(driverName, dataSourceName string) (*DB, error)	|创建连接但不立即验证	|延迟验证的场景	|*sqlx.DB|
||db.Ping() error	|检查数据库连接是否有效	|健康检查、连接池维护	|error|
|​单行查询​	|db.Get(dest interface{}, query string, args ...interface{}) error	|查询单行数据到结构体	|按 ID 查询、获取单个对象	|error|
||db.GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error	|带上下文的单行查询	|需要超时控制或取消的单行查询	|error|
|​多行查询​	|db.Select(dest interface{}, query string, args ...interface{}) error	|查询多行数据到结构体切片	|列表查询、获取多个对象	|error|
||db.SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error	|带上下文的多行查询	|需要超时控制或取消的多行查询	|error|
|​原始查询​	|db.Query(query string, args ...interface{}) (*sql.Rows, error)	|执行 SQL 返回原始 Rows 对象	|复杂查询、需要手动扫描结果	|*sql.Rows, error|
||db.QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)	|带上下文的原始查询	|需要控制执行时间的复杂查询	|*sql.Rows, error|
|​写操作​	|db.Exec(query string, args ...interface{}) (sql.Result, error)	|执行写操作（INSERT/UPDATE/DELETE）	|数据修改操作	|sql.Result, error|
||db.ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)	|带上下文的写操作	需要超时控制的修改操作	|sql.Result, error|
|​命名参数操作​	|db.NamedExec(query string, arg interface{}) (sql.Result, error)	|使用命名参数的写操作	|使用结构体作为参数的数据操作	|sql.Result, error|
||db.NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)	|使用命名参数的查询	|使用结构体作为参数的查询	|*sqlx.Rows, error|
|​IN 查询​	|sqlx.In(query string, args ...interface{}) (string, []interface{}, error)	|生成 IN 查询的安全 SQL	查询多个值（如 WHERE id IN (?, ?, ?)）	(string, []interface{}, error)
|​事务管理​	|db.Beginx() (*sqlx.Tx, error)	|开启事务	|需要原子性操作的场景（转账、订单处理）	|*sqlx.Tx|
||db.BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)	|带上下文的事务开启	需|要控制事务执行时间的场景	|*sqlx.Tx|
|​预处理语句​	|db.Prepare(query string) (*sql.Stmt, error)	|创建预处理语句	|重复执行的 SQL 语句	|*sql.Stmt|
||db.Preparex(query string) (*sqlx.Stmt, error)	|创建支持 Get/Select 的预处理语句	|复用查询结构	|*sqlx.Stmt|
||db.PrepareNamed(query string) (*sqlx.NamedStmt, error)	|创建命名参数预处理语句	|复用命名参数查询	|*sqlx.NamedStmt|

---

### 3.2 事务对象方法 (sqlx.Tx)

- **前置方法**：
	- `db.Beginx() (*sqlx.Tx, error)`
	- `db.BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)`

|方法​	|​说明​|
|:----|:----|
|tx.Commit() error	|提交事务|
|tx.Rollback() error	|回滚事务|
|tx.Exec()	|与 db.Exec 相同，但在事务中执行|
|tx.ExecContext()	|与 db.ExecContext 相同，但在事务中执行|
|tx.Get()	|与 db.Get 相同，但在事务中执行|
|tx.GetContext()	|与 db.GetContext 相同，但在事务中执行|
|tx.Select()	|与 db.Select 相同，但在事务中执行|
|tx.SelectContext()	|与 db.SelectContext 相同，但在事务中执行|
|tx.NamedExec()	|与 db.NamedExec 相同，但在事务中执行|
|tx.Prepare()	|与 db.Prepare 相同，但在事务中执行|

---

### 3.3 其他方法

|​方法​	|​说明​|
|:----|:----|
|db.Rebind(query string) string	|适配不同数据库的参数占位符（? → 1,2 等）|
|db.MustExec()	|执行 SQL，出错时 panic（适合初始化脚本）|
|db.MustBegin()	|开启事务，出错时 panic|
|db.StructScan(row *sql.Rows, dest interface{}) error	手|动将行扫描到结构体|

---

### 3.4 方法的简单使用示例

#### 3.4.1 简单查询

```go
user := User{}
err := db.Get(&user, "SELECT * FROM users WHERE id = ?", 1)
```

---

#### 3.4.2 ​命名参数

```go
result, err := db.NamedExec(
    "INSERT INTO users (name, email) VALUES (:name, :email)", 
    User{Name: "Alice", Email: "alice@example.com"}
)
```

---

#### 3.4.3 事务处理

```go
tx, err := db.Beginx()
// ...执行事务操作...
if err != nil {
    tx.Rollback()
} else {
    tx.Commit()
}
```

---

#### 3.4.4 ​IN 查询

```go
query, args, _ := sqlx.In("SELECT * FROM users WHERE id IN (?)", []int{1,2,3})
users := []User{}
err := db.Select(&users, db.Rebind(query), args...)
```

---

#### 3.4.5 ​预处理

```go
stmt, err := db.Preparex("SELECT * FROM users WHERE id = ?")
// 复用预处理语句
for _, id := range ids {
    var user User
    stmt.Get(&user, id)
}
```

---