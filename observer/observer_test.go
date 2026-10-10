package main

import "testing"

// countingObserver 是一个测试用的观察者，用来记录收到通知的次数。
type countingObserver struct {
	count int
}

func (c *countingObserver) Update(event string) {
	c.count++
}

// TestNotifyAll 验证注册了几个观察者，通知时就会调用几次。
func TestNotifyAll(t *testing.T) {
	order := &Order{}
	obs := &countingObserver{}

	order.Register(obs)
	order.Register(obs) // 注册两次（实际业务一般不会，这里只为验证）

	order.Notify("测试事件")

	if obs.count != 2 {
		t.Fatalf("期望被通知 2 次，实际 %d 次", obs.count)
	}
}

// TestDeregister 验证退订后不再收到通知。
func TestDeregister(t *testing.T) {
	order := &Order{}
	obs := &countingObserver{}

	order.Register(obs)
	order.Notify("第一次")
	order.Deregister(obs)
	order.Notify("第二次")

	if obs.count != 1 {
		t.Fatalf("退订后不应再收到通知，期望 1 次，实际 %d 次", obs.count)
	}
}
