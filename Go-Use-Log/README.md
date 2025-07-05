# Go 结构化日志库 log/slog

> 【一、什么是结构化日志】、【二、slog 与 三方库】为 deepseek 给出的解读。

> 可能由于 deepseek 参照的版本不是最新的，部分内容存在差异

> 从【三、slog 核心组件】开始，结合 ai 给出的大纲自行内容总结。

## 一、什么是结构化日志

**传统日志**（如：`log.Println`）：

```text
2025/07/05 12:00:00 用户 Alice 登录，IP: 192.168.1.1
```

**结构化日志**（slog）：

- **优势**​：机器可解析、键值对清晰、便于日志分析系统处理。

```json
{
  "time": "2023-10-01T12:00:00Z",
  "level": "INFO",
  "msg": "用户登录",
  "user": "Alice",
  "ip": "192.168.1.1"
}
```

---

## 二、slog 与 三方库

### 1. slog 日志库特性

|特性​  |​说明​|
|:----|:----|
|​官方标准​	|Go 团队维护，随语言版本升级（1.21+），兼容性有保证|
|​结构化日志​	|默认输出键值对（如 "user=Alice count=3"），取代传统纯文本，方便日志分析系统处理|
|​高性能设计​	|底层优化减少内存分配，性能逼近 Zap/Zerolog|
|​零依赖​	|无需 go get，开箱即用|
|​扩展性强​	|支持自定义日志处理器（Handler），轻松对接文件、JSON、云服务等|
|​兼容旧项目​	|可通过适配器接入 Logrus/Zap 等第三方库|

---

### 2. slog vs 第三方库（核心区别）

|场景​	|​slog​	|​Zap/Logrus​|
|:----|:----|:----|
|​新项目​	|✅ 首选！官方标准+未来生态	|⚠️ 需额外依赖|
|​旧项目迁移​	|✅ 用 slog.SetDefault() |无痛接入	🔄 需改调用写法|
|​性能极致需求​	|🟢 接近 Zap（差 10%~15%）	|⚡ Zap 仍略微领先|
|​自定义输出格式​	|🟢 通过 Handler 灵活实现	|🟢 各库有自己的插件体系|

---

### 3. 注意事项

- **版本要求​**：仅 `Go ​1.21+`​​ 支持（用 `go version` 确认）
- ​**性能关键点**​：传递键值对时避免内存逃逸，推荐使用 `slog.Int()` 等类型函数

```go
// 更高效（减少内存分配）
slog.Info("查询完成", "duration_ms", slog.Int("value", 85))
```

---

## 三、slog 核心组件


|概念	|作用|
|:-----|:-----|
|​Logger​	|主日志对象，调用 .Info()/.Error() 等方法记录日志|
|​Handler​	|日志处理器（核心！），决定日志如何输出、格式和过滤条件|
|​Record​	|单条日志内容（包含时间、级别、消息、键值对）|
|​Attr​	|键值对数据单元（如 slog.String("ip", "1.1.1.1"))|

---

### 1. 四大核心组件关联图

```plaintext
[Logger] -> 记录日志 (Info/Error/Debug/...)
  │
  ▼ 生成 [Record] (包含时间/级别/消息/Attrs)
  │
  ▼
[Handler] (处理器接口)  <─┬─ TextHandler (文本)
  ▲                     ├─ JSONHandler (JSON)
  │                     └─ 自定义 Handler
  │
  └─ 输出目标 (os.File, buffer, 网络等)
```

---

### 2. Logger：日志记录的门面

**源码定位**：`log/slog/logger.go`

```go
type Logger struct {
    handler Handler     // for structured logging
}
```
---

#### 2.1 简单使用示例

**【示例一】**：仅输出指定信息

```go
// 最基本的使用示例：仅输出指定信息
func main() {
    // Info(msg string, args ...any)
    slog.Info("hello world")
}

// 2025/07/04 22:58:13 INFO hello world    // 日志输出
```

最简单的日志输出信息：
- 时间：2025/07/04 22:48:10
- 日志等级：INFO
- 指定的msg信息：hello world

---

**【示例二】**：输出结构化日志信息

- 指定的入参：`args ...any`
- 输出的格式：key=value

```go
func main() {
    // Info(msg string, args ...any)
	slog.Info("查询完成", "duration_ms", slog.Int("value", 85))
}

// 2025/07/04 21:31:25 INFO 查询完成 duration_ms="value=85"
```

---

#### 2.2 日志方法源码浅识（以Info为例）

**源码**：`log/slog/logger.go`

```go
// 无上下文版本
// Info calls [Logger.Info] on the default logger.
func Info(msg string, args ...any) {
    // 委托给内部 log 方法
    // 内部方法实际是需要上下文的，但是【无上下文版本】方法会在方法里初始一个空白上下文
	Default().log(context.Background(), LevelInfo, msg, args...)
}

// 带上下文版本
// InfoContext calls [Logger.InfoContext] on the default logger.
func InfoContext(ctx context.Context, msg string, args ...any) {
    // 委托给内部 log 方法
	Default().log(ctx, LevelInfo, msg, args...)
}
```

在业务使用中，我们可以直接调用全局方法来进行日志的指定输出，如上源码：
- 对于全局日志方法是有一个默认实例的，通过默认实例来调用`Logger`内部方法
- 默认实例方法如下：
```go
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

**`Logger.log`方法**

```go
// log is the low-level logging method for methods that take ...any.
// It must always be called directly by an exported logging method
// or function, because it uses a fixed call depth to obtain the pc.
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
	r.Add(args...)
	if ctx == nil {
		ctx = context.Background()
	}
	_ = l.Handler().Handle(ctx, r)
}
```

---

#### 2.3 设计亮点
- 轻量级封装​：Logger 本身不处理格式/输出，仅转发给 Handler
- ​不可变设计​：WithAttrs() 返回新 Logger 实例，避免并发冲突

```go
// 不可变设计
func (l *Logger) With(args ...any) *Logger {
	if len(args) == 0 {
		return l
	}
    // 对原始日志实例进行深拷贝
	c := l.clone()
    // 调用 Handler 的 WithAttrs 方法一个新的 Handler
	c.handler = l.handler.WithAttrs(argsToAttrSlice(args))
	return c
}

func (l *Logger) clone() *Logger {
    // 深拷贝 Logger 返回新的 实例
	c := *l
	return &c
}
```

---

### 3. Handler 接口：处理器

**源码定位**​：`log/slog/handler.go`

---