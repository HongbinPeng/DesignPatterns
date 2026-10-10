package main

import "testing"

// TestCacheDecoratorHit 验证缓存装饰器：第二次读取命中缓存，且结果正确。
func TestCacheDecoratorHit(t *testing.T) {
	mem := NewMemorySource()
	cached := &CachingSource{inner: mem, cache: map[string]string{}}

	// 第一次：未命中缓存，走底层
	v1, err := cached.Read("name")
	if err != nil || v1 != "小明" {
		t.Fatalf("第一次读取失败：v=%q err=%v", v1, err)
	}

	// 第二次：应该命中缓存
	v2, err := cached.Read("name")
	if err != nil || v2 != "小明" {
		t.Fatalf("第二次读取失败：v=%q err=%v", v2, err)
	}

	// 缓存里应该有这条记录了
	if _, ok := cached.cache["name"]; !ok {
		t.Fatal("期望 name 已被缓存，但缓存里没有")
	}
}

// TestLoggingDecorator 验证日志装饰器会转发请求并正确返回结果。
func TestLoggingDecorator(t *testing.T) {
	mem := NewMemorySource()
	logged := &LoggingSource{inner: mem}

	v, err := logged.Read("level")
	if err != nil || v != "vip" {
		t.Fatalf("读取失败：v=%q err=%v", v, err)
	}
}
