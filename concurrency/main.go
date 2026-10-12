package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// 这个文件围绕一个核心问题展开：
//
//   「多个 goroutine 同时读写同一个变量，会发生什么？怎么解决？」
//
// 从 goroutine 基础，到数据竞争，再到三种解决方案，最后到 sync.Once 单例。
// 建议配合 README.md 一起看，代码里的每种写法在 README 里都有对应讲解。

// say 打印一段文字，用来演示 goroutine 的并发执行。
func say(name string) {
	for i := 0; i < 3; i++ {
		fmt.Printf("   [%s] 第 %d 次\n", name, i)
	}
}

// ---- 一、goroutine 基础 ----

// demoGoroutineBasics 演示 go 关键字：函数在后台并发执行。
func demoGoroutineBasics() {
	fmt.Println("=== 一、goroutine：用 go 关键字并发执行 ===")
	go say("后台 goroutine")
	say("主 goroutine")
	// ⚠️ 如果不等待，main 跑完程序就退出，后台 goroutine 可能根本没跑完。
	// 这里用 Sleep 简单等一下（仅为演示）；正确姿势见下一节的 WaitGroup。
	time.Sleep(500 * time.Millisecond)
	fmt.Println()
}

// ---- 二、sync.WaitGroup ----

// demoWaitGroup 演示用 WaitGroup 等待所有 goroutine 跑完。
func demoWaitGroup() {
	fmt.Println("=== 二、sync.WaitGroup：等所有 goroutine 跑完 ===")
	var wg sync.WaitGroup
	for i := 1; i <= 5; i++ {
		wg.Add(1) // 每启动一个 goroutine，计数 +1
		go func(id int) {
			defer wg.Done() // 跑完计数 -1
			fmt.Printf("   任务 %d 完成\n", id)
		}(i) // 把 i 作为参数传入，避免闭包捕获同一个循环变量
	}
	wg.Wait() // 阻塞，直到计数归零
	fmt.Println("   全部任务完成 ✅")
	fmt.Println()
}

// ---- 三、数据竞争（错误示范）----

// UnsafeCounter 不安全的计数器：多个 goroutine 同时 Inc 会互相覆盖。
// 仅用于演示数据竞争，生产环境千万别这么写。
type UnsafeCounter struct{ value int }

// Inc 让计数器 +1（存在数据竞争）。
func (c *UnsafeCounter) Inc() { c.value++ }

// Value 返回当前值。
func (c *UnsafeCounter) Value() int { return c.value }

// demoDataRace 演示：1000 个 goroutine 同时 +1，结果却不到 1000。
func demoDataRace() {
	fmt.Println("=== 三、数据竞争：多个 goroutine 同时改一个变量 ===")
	var c UnsafeCounter
	var wg sync.WaitGroup
	const n = 1000
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Inc() // ⚠️ 1000 个 goroutine 同时执行 c.value++
		}()
	}
	wg.Wait()
	fmt.Printf("    期望 %d，实际 %d（每次运行都可能不同）\n", n, c.Value())
	fmt.Println("    这就是数据竞争：用 `go run -race ./concurrency` 能抓出现场")
	fmt.Println()
}

// ---- 四、解决方案一：sync.Mutex ----

// MutexCounter 用互斥锁保护的计数器，并发安全。
type MutexCounter struct {
	mu    sync.Mutex
	value int
}

// Inc 进入临界区再加 1，同一时刻只有一个 goroutine 能通过 Lock。
func (c *MutexCounter) Inc() {
	c.mu.Lock() // 加锁：临界区开始
	c.value++
	c.mu.Unlock() // 解锁：临界区结束
}

// Value 读值也要加锁，避免读的时候正好有人写。
func (c *MutexCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

// demoMutex 演示互斥锁方案。
func demoMutex() {
	fmt.Println("=== 四、解决方案一：sync.Mutex 互斥锁 ===")
	var c MutexCounter
	var wg sync.WaitGroup
	const n = 1000
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Inc()
		}()
	}
	wg.Wait()
	fmt.Printf("    加锁后：期望 %d，实际 %d ✅\n", n, c.Value())
	fmt.Println()
}

// ---- 五、解决方案二：sync/atomic ----

// AtomicCounter 用原子操作实现的计数器，比互斥锁更轻量，适合简单计数。
type AtomicCounter struct{ value int64 }

// Inc 原子加 1，一条 CPU 指令完成，天然无竞争。
func (c *AtomicCounter) Inc() { atomic.AddInt64(&c.value, 1) }

// Value 原子读取。
func (c *AtomicCounter) Value() int64 { return atomic.LoadInt64(&c.value) }

// demoAtomic 演示原子操作方案。
func demoAtomic() {
	fmt.Println("=== 五、解决方案二：sync/atomic 原子操作 ===")
	var c AtomicCounter
	var wg sync.WaitGroup
	const n = 1000
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Inc()
		}()
	}
	wg.Wait()
	fmt.Printf("    原子操作后：期望 %d，实际 %d ✅\n", n, c.Value())
	fmt.Println()
}

// ---- 六、解决方案三：channel ----

// countViaChannel 用「通信代替共享内存」：
// 每个 goroutine 把自己的贡献通过 channel 上交，由调用方汇总。
// 全程没有共享的可变变量，天然无竞争。
func countViaChannel(n int) int {
	results := make(chan int, n)
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			results <- n // 贡献 1，通过 channel 传出去
		}(i)
	}
	wg.Wait()
	close(results)

	sum := 0
	for v := range results {
		fmt.Printf("我是第%d个线程输入的信息", v) // 只有这一处累加，没有并发写
	}
	return sum
}

// demoChannel 演示 channel 方案。
func demoChannel() {
	fmt.Println("=== 六、解决方案三：channel（通信代替共享内存）===")
	const n = 1000
	fmt.Printf("    汇总结果：期望 %d，实际 %d ✅\n", n, countViaChannel(n))
	fmt.Println()
}

// ---- 七、sync.Once 与单例 ----

// DB 模拟一个「昂贵」的数据库连接，整个程序只需要一份。
type DB struct{ id int }

var (
	db   *DB
	once sync.Once
)

// GetDB 返回全局唯一的 DB 实例，并发安全。
// 不管多少个 goroutine 同时第一次调用，初始化代码都只执行一次。
// 这就是单例模式在 Go 里的标准写法（对照 singleton/ 文件夹）。
func GetDB() *DB {
	once.Do(func() {
		db = &DB{id: 1}
		fmt.Println("   [初始化] 建立数据库连接（只应该打印这一次）")
	})
	return db
}

// demoOnce 演示 sync.Once 保证并发下单例只初始化一次。
func demoOnce() {
	fmt.Println("=== 七、sync.Once：并发安全的单例 ===")
	var wg sync.WaitGroup
	const n = 10
	ids := make([]*DB, n)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			ids[idx] = GetDB()
		}(i)
	}
	wg.Wait()

	same := true
	for i := 1; i < n; i++ {
		if ids[i] != ids[0] {
			same = false
		}
	}
	fmt.Printf("    10 个 goroutine 拿到的都是同一个实例吗？%v\n", same)
	fmt.Println("    对照：项目里的 singleton/ 文件夹就是完整的单例模式实现")
	fmt.Println()
}
func main() {
	// demoGoroutineBasics()
	// demoWaitGroup()
	// demoDataRace()
	// demoMutex()
	// demoAtomic()
	// demoChannel()
	demoOnce()
}
