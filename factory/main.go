package main

import "fmt"

// Payment 支付接口：所有支付方式都实现它。
// 调用方只需要认识这个接口，不需要知道背后是支付宝还是微信。
type Payment interface {
	Pay(amount float64) string
}

// ---- 各种支付实现 ----

// Alipay 支付宝
type Alipay struct{}

func (Alipay) Pay(amount float64) string {
	return fmt.Sprintf("【支付宝】支付 %.2f 元", amount)
}

// WechatPay 微信支付
type WechatPay struct{}

func (WechatPay) Pay(amount float64) string {
	return fmt.Sprintf("【微信】支付 %.2f 元", amount)
}

// UnionPay 银联
type UnionPay struct{}

func (UnionPay) Pay(amount float64) string {
	return fmt.Sprintf("【银联】支付 %.2f 元", amount)
}

// NewPayment 工厂函数：根据传入的类型，返回对应的支付实现。
//
// 这是工厂模式的核心——「创建哪个对象」的决定，都集中在这里。
// 以后要新增一种支付方式（比如京东支付），只需：
//   1. 加一个实现了 Payment 接口的结构体
//   2. 在这个 switch 里多一个 case
// 其他调用方的代码一行都不用改。
func NewPayment(payType string) (Payment, error) {
	switch payType {
	case "alipay":
		return &Alipay{}, nil
	case "wechat":
		return &WechatPay{}, nil
	case "union":
		return &UnionPay{}, nil
	default:
		return nil, fmt.Errorf("不支持的支付方式：%s", payType)
	}
}

func main() {
	// 调用方完全不关心具体实现，只跟工厂和接口打交道。
	for _, payType := range []string{"alipay", "wechat", "union"} {
		p, err := NewPayment(payType)
		if err != nil {
			fmt.Println("创建失败：", err)
			continue
		}
		fmt.Println(p.Pay(100.50))
	}

	// 传入一个未知类型，看看工厂怎么报错。
	if _, err := NewPayment("bitcoin"); err != nil {
		fmt.Println("创建失败：", err)
	}
}
