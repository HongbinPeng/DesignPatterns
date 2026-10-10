package main

import (
	"sync"
	"testing"
)

// TestSingletonSameInstance 验证多次调用拿到的是同一个实例。
func TestSingletonSameInstance(t *testing.T) {
	a := GetConfig()
	b := GetConfig()

	if a != b {
		t.Fatalf("期望两次调用拿到同一个实例，但实际不是同一个")
	}
}

// TestSingletonConcurrent 验证并发场景下也只有一个实例，
// 且 sync.Once 保证初始化只执行一次（不会数据竞争）。
func TestSingletonConcurrent(t *testing.T) {
	const goroutines = 100

	var wg sync.WaitGroup
	results := make([]*Config, goroutines)

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			results[idx] = GetConfig()
		}(i)
	}
	wg.Wait()

	// 所有 goroutine 拿到的必须是同一个实例
	first := results[0]
	for i := 1; i < goroutines; i++ {
		if results[i] != first {
			t.Fatalf("并发下出现了不同实例，单例被破坏")
		}
	}
}
