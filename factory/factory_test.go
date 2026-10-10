package main

import (
	"fmt"
	"testing"
)

func TestNewPayment(t *testing.T) {
	cases := []struct {
		payType  string
		wantType string // 期望返回的具体类型（%T 的格式）
	}{
		{"alipay", "*main.Alipay"},
		{"wechat", "*main.WechatPay"},
		{"union", "*main.UnionPay"},
	}

	for _, c := range cases {
		p, err := NewPayment(c.payType)
		if err != nil {
			t.Fatalf("创建 %s 失败：%v", c.payType, err)
		}
		// %T 打印变量的具体类型，用来验证工厂返回了正确的实现
		if got := fmt.Sprintf("%T", p); got != c.wantType {
			t.Fatalf("期望类型 %s，实际 %s", c.wantType, got)
		}
	}
}

// TestNewPaymentUnknown 验证未知类型会返回错误。
func TestNewPaymentUnknown(t *testing.T) {
	if _, err := NewPayment("bitcoin"); err == nil {
		t.Fatal("期望未知类型返回错误，但返回了 nil")
	}
}
