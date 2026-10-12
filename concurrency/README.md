# 并发专题：goroutine 与 sync.Once

> **一句话**：Go 用 `go` 关键字让代码并发跑，用 `sync` 包解决并发带来的「数据竞争」，其中 `sync.Once` 是单例模式的官方外挂。

这个专题不是设计模式，而是**设计模式（尤其是单例）背后的地基**。搞懂它，你才知道 `singleton/` 里那句 `once.Do(...)` 为什么安全。

---

## 一、goroutine：并发执行的单位

在函数调用前加一个 `go`，这个函数就在后台**同时**跑，当前代码不等它、继续往下走：

```go
go say("后台")   // 在后台并发执行
say("主线程")     // 主 goroutine 继续执行
```

它和线程最大的区别是**轻**：

| | 线程 | goroutine |
| --- | --- | --- |
| 栈大小 | 固定 ~1MB | 初始几 KB，按需增长 |
| 创建/切换成本 | 高（系统调用） | 极低（用户态调度） |
| 常见数量 | 几百上千 | 几万、几十万 |

所以 Go 后端能「一个请求一个 goroutine」，扛高并发是刻在语言里的。

## 二、第一个坑：main 跑完，goroutine 一起死

`go` 是非阻塞的，`main` 返回，整个进程就退出，**不会等你后台的 goroutine**。

正确做法是用 **`sync.WaitGroup`** 等待：

```go
var wg sync.WaitGroup
for i := 0; i < 5; i++ {
	wg.Add(1)              // 计数 +1
	go func(id int) {
		defer wg.Done()    // 跑完计数 -1
		fmt.Println(id)
	}(i)                   // ⚠️ 把 i 当参数传，别在闭包里直接引用循环变量
}
wg.Wait()                  // 阻塞到计数归零
```

## 三、核心问题：数据竞争（Data Race）

**多个 goroutine 同时读写同一个变量**，就会出现互相覆盖。看这个计数器：

```go
var counter int
for i := 0; i < 1000; i++ {
	go func() { counter++ }()   // ⚠️ counter++ 是「读+加1+写」三步
}
```

`counter++` 不是原子操作，两个 goroutine 可能都读到 100，各自 +1 写成 101，两次操作只加了一次。结果就是 **1000 个 goroutine 各加 1 次，最终却不到 1000**。

**发现它的神器**：`go run -race` / `go test -race`。加了 `-race` 跑一遍，编译器会精确告诉你竞争发生在哪一行。**后端开发请养成用 `-race` 跑测试的习惯。**

## 四、三种解决方案 & 怎么选

### 方案 1：`sync.Mutex` 互斥锁

```go
type Counter struct {
	mu    sync.Mutex
	value int
}

func (c *Counter) Inc() {
	c.mu.Lock()   // 同一时刻只有一个 goroutine 能通过
	c.value++
	c.mu.Unlock()
}
```

**思想**：把「读改写」圈成临界区，锁保证同一时刻只有一个人进去。**读值也要加锁**（否则读的时候别人正在写，还是竞争）。

### 方案 2：`sync/atomic` 原子操作

```go
var counter int64
atomic.AddInt64(&counter, 1)   // 一条 CPU 指令完成，天然无竞争
```

**思想**：把「读+加+写」变成 CPU 级别的单条原子指令。比互斥锁**更轻量**，但只适合**简单的单值**操作（加、减、交换、CAS）。复杂逻辑还是得用锁。

### 方案 3：channel（通信代替共享内存）

```go
results := make(chan int, n)
// 每个 goroutine 把结果发到 channel，而不是去改共享变量
go func() { results <- 1 }()
// 最后由一个地方统一汇总，没有并发写
```

**思想**：Go 的名言——「不要通过共享内存来通信，而要通过通信来共享内存」。数据在 goroutine 之间**传递**，而不是被**共享修改**，从根上消灭竞争。

### 怎么选？

| 场景 | 推荐 |
| --- | --- |
| 简单计数、标志位、自增 | `atomic`（最轻） |
| 保护一段多步骤逻辑 / 一个结构体 | `Mutex`（最通用） |
| goroutine 之间传数据、协调顺序 | `channel`（最 Go） |

> 三条路都能解决竞争。新手最容易混淆的是「什么时候用锁、什么时候用 channel」——一句话：**保护共享状态用锁，传递数据用 channel**。

## 五、`sync.Once`：并发安全的「只执行一次」

单例模式最核心的需求是：**初始化只做一次，但可能被很多 goroutine 同时触发**。

如果手写 `if instance == nil { ... }`，在并发下会初始化多次（著名的「双重检查锁」问题）。`sync.Once` 把这个逻辑封装好了：

```go
var (
	db   *DB
	once sync.Once
)

func GetDB() *DB {
	once.Do(func() {
		// 无论多少个 goroutine 同时第一次调用，
		// 这里都只会执行一次，且线程安全
		db = &DB{...}
	})
	return db
}
```

要点：

- 只有**第一个到达**的 goroutine 真正执行 `Do` 里的函数，其余**阻塞等它跑完**，之后直接跳过；
- 后续所有 `Do` 调用都直接返回，**不再执行**；
- ⚠️ 如果 `Do` 里的函数 **panic**，`once` 会认为「已经执行过」，后续不再重试——初始化代码里要处理好错误，别让它 panic。

对照项目里的 [singleton/](../singleton/) 文件夹，就是这套写法的完整落地。

## 六、常见坑速查 🕳️

1. **main 退出太早** → 后台 goroutine 没跑完，用 `WaitGroup`。
2. **闭包捕获循环变量** → `go func(){ use(i) }()` 里的 `i` 可能全是最后一个值，把 `i` 作为参数传入。
3. **数据竞争** → 并发读写同一变量，用 `-race` 检测，用锁 / atomic / channel 解决。
4. **锁里做耗时操作** → 比如锁内 `time.Sleep` 或发网络请求，会拖垮所有等待者。**锁要短，进了临界区就干正事。**
5. **复制 `sync.Mutex`** → 锁一旦被值拷贝（比如 struct 值传递），锁就失效了。含锁的结构体要用指针传。
6. **`once.Do` 里 panic** → 之后不再重试，初始化代码要兜住错误。
7. **滥用单例** → 全局状态会让代码难测试、测试之间互相污染。不是「唯一」需求就别用单例，优先依赖注入。

## 怎么运行

```bash
# 在项目根目录执行
go run ./concurrency        # 看七段演示，一步步理解
go run -race ./concurrency  # 加 -race，能看到数据竞争的报错现场
go test -race ./concurrency # 跑并发测试，验证三种方案都安全
```

建议阅读顺序：先跑 `go run ./concurrency` 看输出，再对照 `main.go` 从第一节读到第七节，最后合上代码问自己：**我最近写的哪个功能里，有多个 goroutine 在抢同一个变量？**
