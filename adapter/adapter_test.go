package main

import "testing"

// TestAliyunAdapter 验证阿里云适配器实现了统一接口，且能正常发送。
func TestAliyunAdapter(t *testing.T) {
	var sender SmsSender = &AliyunAdapter{sdk: &AliyunSmsSDK{}}

	if err := sender.Send("13800138000", "test"); err != nil {
		t.Fatalf("阿里云适配器发送失败：%v", err)
	}
}

// TestTencentAdapter 验证腾讯云适配器同样实现了统一接口。
func TestTencentAdapter(t *testing.T) {
	var sender SmsSender = &TencentAdapter{sdk: &TencentSmsSDK{}}

	if err := sender.Send("13900139000", "test"); err != nil {
		t.Fatalf("腾讯云适配器发送失败：%v", err)
	}
}

// TestBothImplementSameInterface 验证两个适配器都是 SmsSender 类型，
// 这正说明业务代码可以无差别地使用它们、自由切换。
func TestBothImplementSameInterface(t *testing.T) {
	var _ SmsSender = &AliyunAdapter{}
	var _ SmsSender = &TencentAdapter{}
}
