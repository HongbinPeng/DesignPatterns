package main

import (
	"sync"
	"testing"
)

// runConcurrently 用 n 个 goroutine 并发执行 fn，并等待全部完成。
// 适合给「不需要下标」的并发测试复用。
func runConcurrently(n int, fn func()) {
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fn()
		}()
	}
	wg.Wait()
}

// TestMutexCounterConcurrent 验证互斥锁方案在并发下计数正确。
func TestMutexCounterConcurrent(t *testing.T) {
	c := &MutexCounter{}
	runConcurrently(1000, c.Inc)
	if got := c.Value(); got != 1000 {
		t.Fatalf("期望计数为 1000，实际为 %d", got)
	}
}

// TestAtomicCounterConcurrent 验证原子操作方案在并发下计数正确。
func TestAtomicCounterConcurrent(t *testing.T) {
	c := &AtomicCounter{}
	runConcurrently(1000, c.Inc)
	if got := c.Value(); got != 1000 {
		t.Fatalf("期望计数为 1000，实际为 %d", got)
	}
}

// TestCountViaChannelConcurrent 验证 channel 方案在并发下计数正确。
func TestCountViaChannelConcurrent(t *testing.T) {
	if got := countViaChannel(1000); got != 1000 {
		t.Fatalf("期望计数为 1000，实际为 %d", got)
	}
}

// TestGetDBSingletonConcurrent 验证 sync.Once 单例在并发下只有一个实例。
func TestGetDBSingletonConcurrent(t *testing.T) {
	const n = 100
	results := make([]*DB, n)

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			results[idx] = GetDB()
		}(i)
	}
	wg.Wait()

	first := results[0]
	for i := 1; i < n; i++ {
		if results[i] != first {
			t.Fatalf("并发下出现了不同实例，单例被破坏")
		}
	}
}
