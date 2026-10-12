package main

import "testing"

// 并发自增场景下，三种方案的性能对比。
// 运行方式（在项目根目录）：
//
//	go test -bench=. -benchmem -run=^$ ./concurrency
//
// 注意：这三种方案并不是「同一个东西」——atomic/Mutex 是保护共享变量，
// channel 是传递数据。这里放一起只是让你直观感受它们的底层开销差距。

// BenchmarkMutexCounter 互斥锁方案：每个 goroutine 不断 Inc。
func BenchmarkMutexCounter(b *testing.B) {
	c := &MutexCounter{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Inc()
		}
	})
}

// BenchmarkAtomicCounter 原子操作方案。
func BenchmarkAtomicCounter(b *testing.B) {
	c := &AtomicCounter{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Inc()
		}
	})
}

// BenchmarkChannelCounter channel 方案：一个专职 goroutine 接收，其余 goroutine 发送。
func BenchmarkChannelCounter(b *testing.B) {
	ch := make(chan int, 4096)
	done := make(chan struct{})
	go func() {
		for range ch { // 接收方：负责"累加"（这里只消费掉）
		}
		close(done)
	}()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ch <- 1
		}
	})

	close(ch)
	<-done
}
