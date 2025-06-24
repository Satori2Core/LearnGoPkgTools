# Go Use MySQL —— 使用 Go 语言操作 MySQL

---

## 驱动安装与数据库基本连接

### 1. 安装 MySQL 驱动（Go包）

在项目目录中执行：

```bash
go get github.com/go-sql-driver/mysql
```

---

### 2. Go 程序中的基础连接设置

前置需求：
- 本地有数据库（安装：[【环境搭建】项目开发数据库选择指南：从类型特性到实战决策 —— Mysql&Redis](https://satori2core.github.io/notes/noteroot/project/%E6%95%B0%E6%8D%AE%E5%BA%93%E9%80%89%E6%8B%A9%E7%AD%96%E7%95%A5%E4%B8%8E%E5%AE%89%E8%A3%85.html)）
- 建立一个`testdb`库

DSN(Data Source Name)配置格式：
- 格式：`[username]:[password]@tcp([host]:[port])/[database]`

```go
// 注意：import _ "github.com/go-sql-driver/mysql" // 匿名导入驱动

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
```

---

**【程序连接验证】**

```bash
# 示例：开始 -----------------------------------------------------------------------------
➜  ConnectionToMySQL git:(learn/GoUseMySQL) ✗ ll
总计 4.0K
-rw-rw-r-- 1 devuser devuser 637  6月 24 23:05 main.go
➜  ConnectionToMySQL git:(learn/GoUseMySQL) ✗ go build -o conntest    # 编译
➜  ConnectionToMySQL git:(learn/GoUseMySQL) ✗ ./conntest              # 执行测试
成功连接到MySQL服务器!
➜  ConnectionToMySQL git:(learn/GoUseMySQL) ✗ 
# 示例：结束 -----------------------------------------------------------------------------
```

---