# 结构化日志库[Go] log/slog 使用指南

## 一、基础日志方法：核心 API 解析 与 使用

### 1. 直接调用方法输出（包级函数）

- **快速上手**：slog提供了顶级函数用于快速记录日志，这些函数直接操作一个默认的全局Logger。
	- **优点**​：开箱即用，适合小型项目或脚本
	- **缺点**​：无法自定义输出格式和级别过滤

```go
// 最基础用法
slog.Debug("连接数据库")                    // 默认输出到 stderr
slog.Info("服务启动", "port", 8080)         // 添加简单键值对
slog.Warn("内存使用高", "usage", 85.3)      // 自动识别基础类型
slog.Error("请求失败", err)                 // 自动处理 error 类型
```
---
**设计原理**:
- **全局默认 logger**​：slog 包维护一个全局的 logger 实例
- **自动初始化**：首次调用时会创建默认 handler（文本格式，stderr 输出）【可以查看init方法】
- **简洁性优先**​：类似标准库 log.Println 的设计哲学
---

**底层源码简化**​ (log/slog/log.go)：
```go
func Info(msg string, args ...any) {
    // 内部 Logger.log 方法支持上下文
    Default().log(context.Background(), LevelInfo, msg, args...)
}

// 上下文版本：
// InfoContext calls [Logger.InfoContext] on the default logger.
func InfoContext(ctx context.Context, msg string, args ...any) {
    // 委托给内部 log 方法
	Default().log(ctx, LevelInfo, msg, args...)
}

// 核心逻辑
var defaultLogger atomic.Pointer[Logger]

// 使用 init 初始一个默认的 Logger 实例
func init() {
    // defaultLogger.Store 源码位于：sync/atomic/type.go
    // newDefaultHandler 即，创建一个默认的 Handler
    // Logger：日志记录的门面，本身不处理格式/输出，仅转发给 Handler 作实际操作
	defaultLogger.Store(New(newDefaultHandler(loginternal.DefaultOutput)))
}

// New creates a new Logger with the given non-nil Handler.
func New(h Handler) *Logger {
	if h == nil {
		panic("nil Handler")
	}
	return &Logger{handler: h}
}

// Default returns the default [Logger].
func Default() *Logger { return defaultLogger.Load() }
```

---

### 2. 自定义 Logger 实例

#### 2.1 创建独立Logger

```go
import (
    "log/slog"
    "os"
)

func main() {
    // 1. 创建自定义Handler（JSON格式）
    jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelDebug, // 设置记录级别
    })

    // 2. 创建Logger实例
    logger := slog.New(jsonHandler)

    // 3. 使用Logger
    logger.Info("服务初始化完成", "port", 8080)
}
```

---

#### 2.2 为什么New()不直接配置？

slog采用职责分离设计：
- Logger：提供日志记录方法（Info/Error等）
- Handler：处理日志的格式、过滤和输出

这种设计支持：
- 多个Logger共享同一个Handler
- 动态替换输出目标（如测试时重定向到buffer）


---

### 3. 属性(Attrs)添加的四种方式

#### 3.1 直接参数传递（可变参数）

```go
logger.Info("用户操作", 
    "user_id", 123,
    "action", "delete",
)
// 这里的可变参数就是我们需要输出内容的键值对
// 特殊地，参数格式不匹配，不会报错，而是输出：!BADKEY="xxxx" 提示
```

---

#### 3.2 使用slog.Type()函数（性能优化）（推荐）

> 使用示例

```go
// 定义处理器基本属性
handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	Level: slog.LevelDebug,
})

// 创建日志实例
logger := slog.New(handler)

// 使用 slog.Type() 来包装输出的日志属性
logger.LogAttrs(context.Background(), slog.LevelInfo, "支付成功",
	slog.Int("user_id", 1001),
	slog.Float64("amount", 99.8),
	slog.Time("time", time.Now()),
)
```

- **为什么更高效？​​**
- 对比普通键值对：
	- **常规方式**："user_id", 1001 -> 需运行时解析为 slog.Attr
	- **Type函数：直接构造Attr结构体，减少内存分配次数**

---


#### 3.3 WithAttrs 预绑定属性（类似日志模板）

```go
// 创建带固定属性的Logger
userLogger := logger.With(
    slog.Int("user_id", 1001),
    slog.String("region", "us-west"),
)

// 后续所有日志自动携带这些属性
userLogger.Info("更改设置", "setting", "dark_mode")
```

> 输出

```json
{
	"time": "2025-07-06T22:21:13.34888977+08:00",
	"level": "INFO",
	"msg": "更改设置",
	"user_id": 1001,		// 通过 With 方法设置的固定属性
	"region": "us-west",	// 通过 With 方法设置的固定属性
	"setting": "dark_mode"
}
```

---

#### 3.4 分组属性(Group)

```go
logger.Info("数据库查询",
    slog.Group("query",
        "sql", "SELECT * FROM users",
        "params", []any{10, "admin"},
    ),
    slog.Group("metrics",
        "duration", 0.87,
        "rows", 100,
    ),
)
```

> 输出

```json
{
  "msg": "数据库查询",
  // 分组一
  "query": {
    "sql": "SELECT * FROM users",
    "params": [10,"admin"]
  },
  // 分组二
  "metrics": {
    "duration": 0.87,
    "rows": 100
  }
}
```

---

#### 3.5 直接传参与Attr的区别分析

**键值对处理机制**​ (log/slog/record.go):
> 源码
```go
// Logger 层日志的核心方法之一
func (l *Logger) log(ctx context.Context, level Level, msg string, args ...any) {
	if !l.Enabled(ctx, level) {
		return
	}
	var pc uintptr
	if !internal.IgnorePC {
		var pcs [1]uintptr
		// skip [runtime.Callers, this function, this function's caller]
		runtime.Callers(3, pcs[:])
		pc = pcs[0]
	}
	r := NewRecord(time.Now(), level, msg, pc)

	// 调用 Add 方法
	r.Add(args...)
	if ctx == nil {
		ctx = context.Background()
	}
	_ = l.Handler().Handle(ctx, r)
}

func (r *Record) Add(args ...any) {
	var a Attr
	for len(args) > 0 {
		// 在此处可以看见，slog底层存储/处理数据格式的统一方式是 Attr （进入argsToAttr方法查看）
		a, args = argsToAttr(args)
		if a.Value.isEmptyGroup() {
			continue
		}
		if r.nFront < len(r.front) {
			r.front[r.nFront] = a
			r.nFront++
		} else {
			if r.back == nil {
				r.back = make([]Attr, 0, countAttrs(args)+1)
			}
			r.back = append(r.back, a)
		}
	}
}

func argsToAttr(args []any) (Attr, []any) {
	switch x := args[0].(type) {
	case string:
		if len(args) == 1 {
			return String(badKey, x), nil
		}
		return Any(x, args[1]), args[2:]

	case Attr:
		return x, args[1:]

	default:
		return Any(badKey, x), args[1:]
	}
}

// Any returns an Attr for the supplied value.
// See [AnyValue] for how values are treated.
func Any(key string, value any) Attr {
	return Attr{key, AnyValue(value)}
}
```

---

## 二、理解 Logger 实例：核心控制单元

### 1. 创建自定义 Logger

```go
// 1. 创建自定义 handler
handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelDebug, // 日志级别
})

// 2. 创建 Logger 实例
logger := slog.New(handler)

// 3. 使用自定义 Logger
logger.Debug("调试信息") // 使用 Debug 级别
logger.Info("业务事件", "data", data)
// 输出：
// {"time":"2025-07-05T23:02:31.490867182+08:00","level":"DEBUG","msg":"调试信息"}
// {"time":"2025-07-05T23:02:31.490951368+08:00","level":"INFO","msg":"业务事件","data":"data"}
```

---

### 2. Logger 的核心能力

> 支持设置默认携带的属性

```go
// 创建带固定字段的子 logger
requestLogger := logger.With(
    slog.String("request_id", "req-123"),
    slog.String("user_agent", "Chrome/120"),
)

// 子 logger 用法（自动携带固定字段）
requestLogger.Info("请求开始", "path", "/api/users")
requestLogger.Warn("参数校验失败", "param", "user_id")

// 输出结果示例：
// {"time":"...","level":"WARN","msg":"参数校验失败","request_id":"req-123","user_agent":"Chrome/120","param":"user_id"}
```

---

## 三、API 设计原理解析 / 最佳使用指南

### 1. 方法调用链解析


当调用 logger.Info("msg", k1, v1)：

1. **​参数接收​** -> 接收消息和键值对（...any）
2. **构建 Record​** -> 包含时间、级别、消息等
3. **​转换 Attrs​** -> 将键值对转换为属性列表
4. **​传递 Handler**​ -> 调用 handler.Handle(ctx, record)
5. **格式化输出​** -> 由具体 Handler 实现

---

### 2. 三种传参方式对比

|方式​	|​使用场景  ​|性能​	|​类型安全​|
|:----|:----|:----|:----|
|键值对 (k, v)​  ​|简单日志，开发阶段	    |中等	 |❌ 运行时检查|
|slog.Attr​     |性能关键路径，生产环境   | 高	   |✅ 编译期检查|
|slog.Group​	 |结构化复杂数据	      |中等	  |❌ 运行时检查|

---

**性能关键源码**​ (log/slog/attr.go)：

```go
// log/slog/attr.go
type Attr struct {
    Key   string
    Value Value
}

// Value 的内部表示 (32字节)
type Value struct {
    // ! 使用联合体(union)优化基础类型存储
    num uint64    // 存放整数/浮点数/布尔值
    str string    // 或字符串引用
    
    // 复杂类型降级存储
    group []Attr  
    any   any     // 其他类型逃逸到堆
}
```

---

### 3. 为什么 With() 返回新 Logger？

**设计哲学​**：
- **不可变对象​**：避免并发修改问题
- ​**零配置污染**​：保持原始 logger 干净
- ​**上下文隔离**​：不同场景使用独立实例

```go
// 安全使用
func HandleRequest(ctx context.Context, r *Request) {
    // 每个请求创建独立logger
    requestLogger := baseLogger.With(
        slog.String("request_id", r.ID),
        slog.String("endpoint", r.Path),
    )
    
    // 记录请求处理日志
    requestLogger.Info("开始处理")
    
    // ...处理逻辑...
}
```

---

## 四、最佳实践指南

### 1. 性能优先实践

> 原始方式（不推荐）：直接传入 参数序列
```go
// 原始方式（有内存分配）
slog.Info("Processing", "item", item.ID, "value", item.Value)
```

---

> 优化方式一：复用 Attr 对象
- 即：定义需要指定输出的字段对象，通过提前定义 []slog.Attr 属性数组，后续复用

```go
type Item struct {
	ID int
	Value string
}

// 优化方式1：复用 Attr 对象
var itemAttrs = []slog.Attr{
    slog.Int("item", 0),
    slog.String("value", ""),
}

func logItem(item Item) {
    itemAttrs[0] = slog.Int("item", item.ID)
    itemAttrs[1] = slog.String("value", item.Value)
    slog.LogAttrs(context.Background(), slog.LevelInfo, "Processing", itemAttrs...)
}
```

---

> 优化方式二：值类型直接替换
```go
// 优化方式2：值类型直接替换
type Item struct {
    idAttr slog.Attr // 预分配属性
}

func (i *Item) SetID(id int) {
    i.idAttr = slog.Int("item", id) // 只修改属性值
}

func (i *Item) Log() {
    slog.LogAttrs(context.Background(), slog.LevelInfo, "Item", i.idAttr)
}
```

---


### 2. 错误处理

> 粗暴的方式：
```go
// 糟糕的反例
slog.Error("DB query failed: %v", err) // 错误：格式化字符串错误
```

---

> 最佳实践一：分开记录

```go
slog.Error("DB query failed", "error", err) 
```

---

> 最佳实践二：错误分组

```go
slog.Error("数据操作失败", 
    slog.Group("operation",
        "type", "update",
        "table", "users",
    ),
    slog.Group("error",
        "code", err.Code,
        "message", err.Message,
    )
)
```

---


## 五、要点总结

### 全局函数适用场景​：
- 简单脚本/工具开发
- 快速调试时的临时日志
- 不需要定制化配置的场景

---

#### ​自定义 Logger 适用场景​：
- 需要控制日志输出目标
- 需要结构化(JSON)格式化
- 不同子系统独立配置级别
- 需要固定上下文字段

---

#### 方法选择优先级

```go
// 1. 常规需求
logger.Info("message", k1, v1)

// 2. 性能优先
logger.LogAttrs(ctx, level, "message", attr1, attr2)

// 3. 复杂上下文
logger.With(fields...).Info("message")
```

---

#### 性能黄金法则​

- 在热点路径中，总是优先使用 LogAttrs + slog.Type 组合；在非热点路径中，选择最简洁的写法。

---
