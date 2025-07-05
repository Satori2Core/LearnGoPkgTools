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

## 三、slog 核心组件 与 源码浅析


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
    // 创建一个 Record 记录
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

#### 2.4 Logger 的完整生命周期

假设调用 `logger.Info("msg", "k", v)`：

- ​参数封装​ → args...any 被转换为 []Attr（见 log/slog/record.go: Add()）
- 构建 Record​ → 合并时间/级别/消息/PC值
（NewRecord() + runtime.Callers()）
- ​Handler 路由​ → logger.handler.Handle(ctx, r)
（多态分发到 Text/JSONHandler）
- ​格式化输出​ → 根据 Handler 类型生成文本/JSON
（TextHandler.format() / JSONHandler.appendValue()）
- ​缓冲写入​ → 最终输出到目标 io.Writer
（带 bufio.Writer 缓存）

---

### 3. Handler 接口：处理器

**源码定位**​：`log/slog/handler.go`

```go
type Handler interface {
    Handle(context.Context, Record) error   // ! 关键方法：处理日志记录
    WithAttrs([]Attr) Handler               // 复制处理器并添加属性
    WithGroup(string) Handler               // 创建属性分组
    Enabled(context.Context, Level) bool    // 动态级别检查
}
```

---

#### 3.1 通用处理器：commonHandler

通用处理器定义了处理器的基础能力属性。

```go
type commonHandler struct {
    // 结构化控制字段 json bool
    // true => output JSON; false => output text
	json              bool 
    // HandlerOptions 处理器选项
	opts              HandlerOptions
	preformattedAttrs []byte
    // groupPrefix 仅用于文本处理函数
	groupPrefix string
	groups      []string // all groups started from WithGroup
	nOpenGroups int      // the number of groups opened in preformattedAttrs
	mu          *sync.Mutex
	w           io.Writer
}

// HandlerOptions are options for a [TextHandler] or [JSONHandler].
// A zero HandlerOptions consists entirely of default values.
type HandlerOptions struct {
    // 启用时（true），记录日志语句的源代码位置（如文件名和行号）。
    // 位置信息会以 SourceKey 属性（默认键为 "source"）添加到输出中。
    // ​默认值​：false（不记录源代码位置）。
	AddSource bool
    // 定义日志记录的最低级别门槛（如 LevelInfo、LevelError），​低于此级别的日志会被丢弃。
    // 如果为 nil，则默认使用 LevelInfo 作为最低级别。
    // 支持动态调整级别（如通过 LevelVar 在运行时修改）。
	Level Leveler
    // 用于自定义属性（Attribute）的预处理逻辑。
	ReplaceAttr func(groups []string, a Attr) Attr
}
```

---

#### 3.2 文本处理器：TextHandler

```go
// TextHandler 是一个 [处理器]。
// 它会将记录以键值对的形式（通过空格分隔，并在每对后跟一个换行符）写入到 [io.Writer] 中。
type TextHandler struct {
	*commonHandler
}

// NewTextHandler creates a [TextHandler] that writes to w,
// using the given options.
// If opts is nil, the default options are used.
func NewTextHandler(w io.Writer, opts *HandlerOptions) *TextHandler {
	if opts == nil {
		opts = &HandlerOptions{}
	}
	return &TextHandler{
		&commonHandler{
			json: false,
			w:    w,
			opts: *opts,
			mu:   &sync.Mutex{},
		},
	}
}

// 处理器调用
func (h *TextHandler) Handle(_ context.Context, r Record) error {
	return h.commonHandler.handle(r)
}
```

---

#### 3.3 JSON 格式处理器：JSONHandler

```go
// JSONHandler is a [Handler] that writes Records to an [io.Writer] as
// line-delimited JSON objects.
type JSONHandler struct {
	*commonHandler
}

// NewJSONHandler creates a [JSONHandler] that writes to w,
// using the given options.
// If opts is nil, the default options are used.
func NewJSONHandler(w io.Writer, opts *HandlerOptions) *JSONHandler {
	if opts == nil {
		opts = &HandlerOptions{}
	}
	return &JSONHandler{
		&commonHandler{
			json: true,
			w:    w,
			opts: *opts,
			mu:   &sync.Mutex{},
		},
	}
}

func (h *JSONHandler) Handle(_ context.Context, r Record) error {
	return h.commonHandler.handle(r)
}
```

---

#### 3.4 commonHandler.handle()

- **代码逻辑**

```plaintext
开始
  │
  ├─ 初始化缓冲区
  ├─ JSON处理: 写入 "{"
  │
  ├─ 处理内置属性
  │   ├─ 时间 → [替换函数?] → 格式化
  │   ├─ 级别 → [替换函数?] → 格式化
  │   ├─ 源码位置（可选）
  │   └─ 消息 → [替换函数?] → 格式化
  │
  ├─ 处理非内置属性（含分组嵌套）
  ├─ JSON处理: 闭合 "}"
  ├─ 写入换行符
  │
  └─ 加锁 → 写入输出流 → 解锁
```

- **方法源码**
> 仅展示大致内容
```go
// handle is the internal implementation of Handler.Handle
// used by TextHandler and JSONHandler.
func (h *commonHandler) handle(r Record) error {
    // 初始化缓冲区
	state := h.newHandleState(buffer.New(), true, "")
	defer state.free()

    // JSON处理: 写入 "{"    
	if h.json {
		state.buf.WriteByte('{')
	}

    // 处理内置属性
	// time
	if !r.Time.IsZero() {
		// ...
	}
	// level
	key := LevelKey
	val := r.Level
	if rep == nil {
		// ...
	} else {
		// ...
	}
	// source
	if h.opts.AddSource {
		// ...
	}
	key = MessageKey
	msg := r.Message
	if rep == nil {
		// ...
	} else {
		// ...
	}

    // JSON处理: 闭合 "}"
    state.appendNonBuiltIns(r)
	state.buf.WriteByte('\n')

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.w.Write(*state.buf)
	return err
}
```

---


### 4. Record：日志数据的原子单元

**源码位置**：`log/slog/record.go`

```go
const nAttrsInline = 5

type Record struct {
    Time    time.Time 
    Message string    
    Level   Level     
    PC      uintptr   

    front [nAttrsInline]Attr  
    nFront int                 
    back   []Attr              
}
```

---

#### 4.1 核心字段解析

**基础信息字段**

```go
Time    time.Time // 日志事件发生的时间
Message string    // 日志消息内容
Level   Level     // 日志级别（如 Debug/Info/Warn/Error）
PC      uintptr   // 程序计数器（用于溯源调用位置）

// PC 字段通过 runtime.Callers 获取，​仅用于​ runtime.CallersFrames 解析调用栈
// 特别注意：不可传给 runtime.FuncForPC，可能导致错误
```

---

**属性存储优化设计**

> 通过复合结构优化小量属性的内存分配

```go
front [nAttrsInline]Attr   // 内联固定大小数组（nAttrsInline=5）
nFront int                 // 实际存储在 front 中的属性数量
back   []Attr              // 超量属性的动态切片，存储实际日志数据
```

---

**Attr字段**

```go
// An Attr is a key-value pair.
type Attr struct {
	Key   string
	Value Value
}
```

---

#### 4.2 设计思想

- 内存优化策略​
    - ​小属性内联存储​：当属性 ≤ nAttrsInline 时，直接使用栈上数组，避免堆分配
    - ​动态扩展机制​：超量属性自动转存 back 切片
    - ​空元素检测​：未使用的 front 元素保持零值，便于错误检测

- 安全使用约束
    - 浅拷贝风险​：副本共享 back 切片的底层数组
    - ​安全操作要求​：
        - 需要修改时 → 用 Record.Clone() 创建深拷贝副本
        - 创建实例 → 必须通过 NewRecord() 工厂函数

- 使用场景约束
    - 安全传递
        - 只读场景可直接传递 Record
    - 需要修改时，通过`Clone`方法，建立独立副本
    - 属性访问
        - 应使用 Record.Attrs(f func(Attr)) 方法遍历属性，而非直接访问 front/back
```go
// 安全传递
// handle方法
func (h *commonHandler) handle(r Record) error {}

// 需要修改时
// 错误：可能污染其他引用
record := origRecord 
record.Message = "modified"

// 正确：创建独立副本
record := origRecord.Clone()
record.Message = "safe modification"
```

---

#### 4.3 Record 创建机制

```go
// Logger.log() 内部 (log/slog/logger.go)
func (l *Logger) log(ctx context.Context, level Level, msg string, args ...any) {
    if !l.Enabled(ctx, level) { // 级别过滤
        return
    }
    var pc uintptr
	if !internal.IgnorePC {
		var pcs [1]uintptr
		// skip [runtime.Callers, this function, this function's caller]
		runtime.Callers(3, pcs[:])
		pc = pcs[0]
	} // 获取调用栈信息
    r := NewRecord(time.Now(), level, msg, pcs[0])
    r.Add(args...) // 添加额外属性
    _ = l.Handler().Handle(ctx, r) // 委托给Handler
}
```

---

#### 4.4 Attr：高效的数据载体

- **关键设计**​：优化键值对的内存分配

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

- 错误使用对比

```go
// ❌ 低效：引发内存逃逸，触发堆分配
slog.Info("login", "userID", getUserID()) 

// ✅ 高效：union存储，无内存逃逸
slog.Info("login", slog.Int("userID", getUserID()))
```