package main

import "fmt"

// Observer 观察者接口：收到主题通知后要做什么。
type Observer interface {
	Update(event string)
}

// Subject 主题接口：负责管理观察者，以及广播通知。
type Subject interface {
	Register(o Observer)   // 订阅
	Deregister(o Observer) // 退订
	Notify(event string)   // 广播
}

// ---- 具体主题：订单 ----

// Order 订单：当订单状态变化（比如支付成功）时，通知所有观察者。
type Order struct {
	observers []Observer
}

func (o *Order) Register(obs Observer) {
	o.observers = append(o.observers, obs)
}

func (o *Order) Deregister(obs Observer) {
	for i, ob := range o.observers {
		if ob == obs {
			// 删除切片中的这个元素
			o.observers = append(o.observers[:i], o.observers[i+1:]...)
			return
		}
	}
}

func (o *Order) Notify(event string) {
	for _, ob := range o.observers {
		ob.Update(event)
	}
}

// ---- 各种观察者：各管各的响应 ----

// SmsNotifier 发短信
type SmsNotifier struct{}

func (SmsNotifier) Update(event string) {
	fmt.Printf("  [短信服务] 收到事件「%s」，发送短信通知用户\n", event)
}

// EmailNotifier 发邮件
type EmailNotifier struct{}

func (EmailNotifier) Update(event string) {
	fmt.Printf("  [邮件服务] 收到事件「%s」，发送邮件通知用户\n", event)
}

// StockService 扣库存
type StockService struct{}

func (StockService) Update(event string) {
	fmt.Printf("  [库存服务] 收到事件「%s」，扣减库存\n", event)
}

func main() {
	order := &Order{}

	// 注册观察者（订阅）
	order.Register(SmsNotifier{})
	order.Register(EmailNotifier{})
	order.Register(StockService{})

	// 订单支付成功，广播通知所有观察者
	fmt.Println("订单支付成功，开始通知：")
	order.Notify("订单支付成功")
}
